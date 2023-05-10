package cluster

import (
	"my_redis/interface/resp"
)

func makeRouter() map[string]CmdFunc{
	router := make(map[string]CmdFunc)

	router["exists"] = defaultFunc
	router["type"] = defaultFunc
	router["get"] = defaultFunc
	router["getset"] = defaultFunc
	router["set"] = defaultFunc
	router["setnx"] = defaultFunc

	router["ping"] = ping
	router["rename"] = rename
	router["renamenx"] = rename
	router["flushdb"] = flushdb
	router["del"] = del
	router["select"] = execSelect

	return router
}

//直接自身或者转发执行的常规命令
func defaultFunc(cluster * ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply{
	key := string(cmdArgs[1])
	peer := cluster.peerPicker.PickNode(key)

	return cluster.relay(peer, conn, cmdArgs)
}
