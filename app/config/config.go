package config

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type Config struct {
	Role             string
	Dir              string
	Dbfilename       string
	Port             string
	MasterHost       string
	MasterPort       string
	MasterReplid     string
	MasterReplOffset string
	Offset           int
	Replicas         []*Node
	Master           Node
}

var (
	configs  *Config = &Config{}
	once     sync.Once
	mu       sync.RWMutex = sync.RWMutex{}
	fieldMap              = map[string]*string{
		"dir":                &configs.Dir,
		"dbfilename":         &configs.Dbfilename,
		"port":               &configs.Port,
		"master_host":        &configs.MasterHost,
		"master_port":        &configs.MasterPort,
		"master_replid":      &configs.MasterReplid,
		"master_repl_offset": &configs.MasterReplOffset,
	}
)

func Get() *Config {
	once.Do(func() {
		dir := flag.String("dir", "", "Directory to serve static files from")
		dbfilename := flag.String("dbfilename", "dump.rdb", "Filename to save the DB to")
		port := flag.String("port", "6379", "Port to listen on")
		replicaof := flag.String("replicaof", "", "Replicate to another Redis server")
		flag.Parse()

		configs.Dir = *dir
		configs.Dbfilename = *dbfilename
		configs.Port = *port
		configs.Role = "master"

		if *replicaof != "" {
			configs.Role = "slave"
			configs.MasterHost = strings.Split(*replicaof, " ")[0]
			configs.MasterPort = strings.Split(*replicaof, " ")[1]
			masterConn, err := net.Dial("tcp", configs.MasterHost+":"+configs.MasterPort)
			if err != nil {
				fmt.Println("Error connecting to master: ", err.Error())
			} else {
				configs.Master = *NewNode(masterConn)
			}
		} else {
			configs.MasterReplid = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
		}

		configs.Replicas = make([]*Node, 0)
	})

	return configs
}

func GetConfigHandler(params []resp.RESP) []byte {
	if len(params) > 1 && params[0].Bulk == "GET" {
		mu.RLock()
		defer mu.RUnlock()
		value, ok := fieldMap[params[1].Bulk]
		if !ok {
			return resp.Nil().Marshal()
		}

		return resp.Array(
			resp.Bulk(params[1].Bulk),
			resp.Bulk(*value),
		).Marshal()
	}

	return resp.Error("ERR wrong number of arguments for 'config' command").Marshal()
}

func AddReplicat(conn net.Conn) {
	mu.Lock()
	configs.Replicas = append(configs.Replicas, NewNode(conn))
	mu.Unlock()
}

func IncOffset(num int) {
	mu.Lock()
	configs.Offset += num
	mu.Unlock()
}
