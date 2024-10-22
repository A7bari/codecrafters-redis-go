package config

import (
	"flag"
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
	Replicas         []*Node
	Master           *Node
}

var (
	configs  *Config
	once     sync.Once
	fieldMap = map[string]*string{
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

		configs = &Config{
			Role:       "master",
			Dir:        *dir,
			Dbfilename: *dbfilename,
			Port:       *port,
		}

		if *replicaof != "" {
			configs.Role = "slave"
			configs.MasterHost = strings.Split(*replicaof, " ")[0]
			configs.MasterPort = strings.Split(*replicaof, " ")[1]
			master, _ := NewNode(
				configs.MasterHost + ":" + configs.MasterPort,
			)

			configs.Master = master
		} else {
			configs.MasterReplid = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
			configs.MasterReplOffset = "0"
		}
	})

	return configs
}

func GetConfigHandler(params []resp.RESP) []byte {
	if len(params) > 1 && params[0].Bulk == "GET" {
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

func AddReplicat(port string) {
	replica, _ := NewNode("0.0.0.0:" + port)
	configs.Replicas = append(configs.Replicas, replica)
}
