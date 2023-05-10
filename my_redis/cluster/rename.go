package cluster

import (
	"my_redis/interface/resp"
	"my_redis/resp/reply"
)

//rename k1 k2
func rename(cluster * ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply{
	if len(cmdArgs) != 3{
		return reply.MakeStandardErrReply("ERR wrong number args")
	}

	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])

	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)

	if srcPeer != destPeer{
		return reply.MakeStandardErrReply("ERR rename must within one peer")
	}

	return cluster.relay(srcPeer, conn, cmdArgs)
}

