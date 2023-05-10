package cluster

import (
	"my_redis/interface/resp"
	"my_redis/resp/reply"
)

func flushdb(cluster * ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply{
	replies := cluster.broadcast(conn, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replies{
		if reply.IsErrReply(r){
			errReply = r.(reply.ErrorReply)
			return errReply
		}
	}

	if errReply != nil{
		return reply.MakeStandardErrReply("error: " + errReply.Error())
	}

	return reply.MakeOkReply()
}
