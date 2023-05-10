package cluster

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"hash/crc32"
	"my_redis/config"
	database2 "my_redis/database"
	"my_redis/interface/database"
	"my_redis/interface/resp"
	"my_redis/lib/consistenthash"
	"my_redis/lib/logger"
	"my_redis/resp/reply"
	"strings"
)

type ClusterDatabase struct {
	self string
	nodes []string
	peerPicker *consistenthash.NodeMap
	peerConnection map[string] *pool.ObjectPool
	db database.Database
}

func NewClusterDatabase() *ClusterDatabase{
	cluster := &ClusterDatabase{
		self: config.Properties.Self,
		db: database2.NewStandaloneDatabase(),
		peerPicker: consistenthash.NewNodeMap(crc32.ChecksumIEEE),
		peerConnection: make(map[string]*pool.ObjectPool),
	}

	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers{
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)

	cluster.nodes = nodes
	cluster.peerPicker.AddNode(nodes...)

	for _, peer := range config.Properties.Peers{
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(context.Background(), &connectionFactory{peer: peer})
	}

	return cluster
}

type CmdFunc func(cluster * ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply
var router = makeRouter()

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil{
			logger.Error(err)
			result = &reply.UnknownErrReply{}
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if ok == false{
		return reply.MakeStandardErrReply("not supported cmd")
	}

	return cmdFunc(cluster, client, args)
}

func (cluster *ClusterDatabase) Close() {
	cluster.db.Close()
}

func (cluster *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	cluster.db.AfterClientClose(conn)
}

