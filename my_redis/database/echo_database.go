package database

import (
	"my_redis/interface/resp"
	"my_redis/resp/reply"
)

type EchoDatabase struct {
	
}

func NewEchoDatabase() *EchoDatabase{
	return &EchoDatabase{}
}

//将收到的已经编码了的args，再原封不动发回去，去测试解码功能
func (e *EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e *EchoDatabase) Close() {

}

func (e *EchoDatabase) AfterClientClose(c resp.Connection) {

}
