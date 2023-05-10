package cluster

import (
	"my_redis/interface/resp"
	"my_redis/resp/reply"
)

func del(cluster * ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply{
	replies := cluster.broadcast(conn, cmdArgs)
	var errReply reply.ErrorReply

	var deleted int64 = 0
	for _, r := range replies{
		if reply.IsErrReply(r){
			errReply = r.(reply.ErrorReply)
			return errReply
		}

		intReply, ok := r.(*reply.IntReply)
		if ok == false{
			return reply.MakeStandardErrReply("error")
		}

		deleted += intReply.Code
	}

	if errReply != nil{
		return reply.MakeStandardErrReply("error: " + errReply.Error())
	}

	return reply.MakeIntReply(deleted)
}