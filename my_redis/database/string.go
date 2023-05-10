package database

import (
	databaseface "my_redis/interface/database"
	"my_redis/interface/resp"
	"my_redis/lib/utils"
	"my_redis/resp/reply"
)

//实现以下命令
//GET
//SET
//SETNX
//GETSET
//STRLEN

//GET
func execGet(db *DB, args [][]byte) resp.Reply{
	key := string(args[0])
	entity, exist := db.GetEntity(key)

	if exist == false{
		return reply.MakeNullBulkReply()
	}

	bytes := entity.Data.(string)

	return reply.MakeBulkReply([]byte(bytes))
}

//SET
func execSet(db *DB, args [][]byte) resp.Reply{
	key := string(args[0])
	val := string(args[1])

	entity := &databaseface.DataEntity{Data: val}
	db.PutEntity(key, entity)


	data := utils.ToCmdLine2("set", args...)
	db.addAof(data)

	return reply.MakeOkReply()
}

//SETNX
func execSetnx(db *DB, args [][]byte) resp.Reply{
	key := string(args[0])
	val := string(args[1])

	entity := &databaseface.DataEntity{Data: val}
	result := db.PutIfAbsent(key, entity)

	data := utils.ToCmdLine2("setnx", args...)
	db.addAof(data)
	return reply.MakeIntReply(int64(result))
}

//GETSET
func execGetSet(db *DB, args [][]byte) resp.Reply{
	key := string(args[0])
	val := string(args[1])

	entityNew := &databaseface.DataEntity{Data: val}
	db.PutEntity(key, entityNew)

	entity, exist := db.GetEntity(key)
	if exist == false{
		return reply.MakeNullBulkReply()
	}
	bytes := entity.Data.([]byte)

	data := utils.ToCmdLine2("getset", args...)
	db.addAof(data)
	return reply.MakeBulkReply(bytes)
}

//STRLEN
func execStrLen(db *DB, args [][]byte) resp.Reply{
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if exist == false{
		return reply.MakeNullBulkReply()
	}

	bytes := len(entity.Data.([]byte))

	db.addAof(args)
	return reply.MakeIntReply(int64(bytes))
}


func init(){
	RegisterCommand("GET", execGet, 2) //get key
	RegisterCommand("SET", execSet, 3) //set key val
	RegisterCommand("SETNX", execSetnx, 2) //SET key val
	RegisterCommand("GETSET", execGetSet, 3) //getset key val
	RegisterCommand("STRLEN", execStrLen, 2) //strlen key
}