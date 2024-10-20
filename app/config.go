package main

import "github.com/codecrafters-io/redis-starter-go/app/resp"

var configs = map[string]string{}

func SetConfig(key, value string) {
	configs[key] = value
}

func GetConfig(params []resp.RESP) resp.RESP {
	if len(params) > 1 && params[0].Bulk == "GET" {
		value, ok := configs[params[1].Bulk]
		if !ok {
			return resp.RESP{
				Type: "nil",
			}
		}

		return resp.RESP{
			Type: "array",
			Array: []resp.RESP{
				{
					Type: "bulk",
					Bulk: params[1].Bulk,
				},
				{
					Type: "bulk",
					Bulk: value,
				},
			},
		}
	}

	return resp.RESP{
		Type: "error",
		Bulk: "CONFIG GET: Invalid command",
	}
}

func GetInfo(params []resp.RESP) resp.RESP {
	return resp.RESP{
		Type: "string",
		Bulk: "node:master",
	}
}
