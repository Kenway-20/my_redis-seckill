package database

import (
	"my_redis/interface/resp"
	"my_redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply{
	return reply.MakePongReply()
}

//利用init方法自动注册Ping方法
func init(){
	RegisterCommand("ping", Ping, 0)
}