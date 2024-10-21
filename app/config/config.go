package config

import "github.com/codecrafters-io/redis-starter-go/app/resp"

var configs = make(map[string]string)

func Set(key, value string) {
	configs[key] = value
}

func Get(key string) string {
	value := configs[key]
	return value
}

func GetConfigHandler(params []resp.RESP) []byte {
	if len(params) > 1 && params[0].Bulk == "GET" {
		value, ok := configs[params[1].Bulk]
		if !ok {
			return resp.Nil().Marshal()
		}

		return resp.Array(
			resp.Bulk(params[1].Bulk),
			resp.Bulk(value),
		).Marshal()
	}

	return resp.Error("ERR wrong number of arguments for 'config' command").Marshal()
}
