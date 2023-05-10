package handler

import (
	"context"
	"io"
	"my_redis/cluster"
	"my_redis/config"
	"my_redis/database"
	databaseface "my_redis/interface/database"
	"my_redis/lib/logger"
	"my_redis/lib/sync/atomic"
	"my_redis/resp/connection"
	"my_redis/resp/parser"
	"my_redis/resp/reply"
	"net"
	"strings"
	"sync"
)

var(
	unknownErrReplyBytes = []byte("-ERR Unknown \r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db databaseface.Database
	closing atomic.Boolean
}


func MakeHandler() *RespHandler{
	var db databaseface.Database

	//判断是用集群模式还是单机模式
	if config.Properties.Self != "" && len(config.Properties.Peers) > 0{
		db = cluster.NewClusterDatabase()
	} else {
		db = database.NewStandaloneDatabase()
	}

	return &RespHandler{db: db}
}

//关闭单个客户端
func (r *RespHandler) CloseClient(client *connection.Connection) error {
	client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client) //因为这个client关闭了，所以要从表示活动连接的map里去除掉
	return nil
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() ==  true{
		conn.Close()
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, 1)

	ch := parser.ParseStream(conn)
	for payload := range ch{
		//出现错误
		if payload.Err != nil{
			//出现连接错误
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF || strings.Contains(payload.Err.Error(), "use of closed network connection"){ //表示客户端正在或已经进行TCP四次挥手要断开连接
				r.CloseClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			//出现协议错误
			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Writer(errReply.ToBytes())
			if err != nil{
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		//exec
		if payload.Data== nil{
			continue
		}
		reply, ok := payload.Data.(*reply.MultiBulkReply)
		if ok == false{
			logger.Error("require multi bulk reply")
			continue
		}

		result := r.db.Exec(client, reply.Args)

		if result != nil{
			client.Writer(result.ToBytes())
		}else{
			client.Writer(unknownErrReplyBytes)
		}

	}
}

//关闭整个服务
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	//关掉map里存的每一个连接
	r.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*connection.Connection)
		client.Close()
		return true
	})
	r.db.Close()
	return nil
}
