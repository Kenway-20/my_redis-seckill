package cluster

import (
	"context"
	"errors"
	"my_redis/interface/resp"
	"my_redis/lib/logger"
	"my_redis/lib/utils"
	"my_redis/resp/client"
	"my_redis/resp/reply"
	"strconv"
)

//拿到目标地址peer的client对象
func (cluster *ClusterDatabase) getPeerClient(peer string) (*client.Client, error){
	pool, ok := cluster.peerConnection[peer]
	if ok == false{
		return nil, errors.New("connection not found")
	}
	object, err := pool.BorrowObject(context.Background())
	if err != nil{
		logger.Error(err)
		return nil, err
	}
	client, ok := object.(*client.Client)
	if ok == false{
		return nil, errors.New("wrong type")
	}

	return client, nil
}

//把getPeerClient()里占据的连接释放回连接池
func (cluster *ClusterDatabase) returnPeerClient(peer string, client *client.Client)  error{
	pool, ok := cluster.peerConnection[peer]
	if ok == false{
		return errors.New("connection not found")
	}

	return pool.ReturnObject(context.Background(), client)
}

//转发命令给其他节点
func (cluster *ClusterDatabase) relay(peer string, conn resp.Connection, args [][]byte) resp.Reply{
	if peer == cluster.self{
		return cluster.db.Exec(conn, args)
	}
	peerClient, err := cluster.getPeerClient(peer)
	if err != nil{
		return reply.MakeStandardErrReply(err.Error())
	}
	defer cluster.returnPeerClient(peer, peerClient)

	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(conn.GetDBIndex())))

	return peerClient.Send(args)
}

//广播，即某些自己执行后其他节点也需要执行的命令，例如
func (cluster *ClusterDatabase) broadcast(conn resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, peer := range cluster.nodes{
		result := cluster.relay(peer, conn, args)
		results[peer] = result
	}

	return results
}
