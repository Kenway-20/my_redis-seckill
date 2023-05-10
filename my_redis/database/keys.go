package database

import (
	"my_redis/interface/resp"
	"my_redis/lib/utils"
	"my_redis/lib/wildcard"
	"my_redis/resp/reply"
)

//实现以下命令
//DEL
//EXISTS
//KEYS
//FLUSHDB
//TYPE
//RENAME
//RENAMENX

func execDel(db *DB, args [][]byte) resp.Reply{
	keys := make([]string, len(args))
	for i, v := range args{
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)

	if deleted > 0{
		data := utils.ToCmdLine2("del", args...)
		db.addAof(data)
	}
	return reply.MakeIntReply(int64(deleted))
}

func execExists (db *DB, args [][]byte) resp.Reply{
	result := int64(0)
	for _, arg := range args{
		key := string(arg)
		_, exist := db.GetEntity(key)
		if exist == true{
			result++
		}
	}

	return reply.MakeIntReply(result)
}

func execFlushDB (db *DB, args [][]byte) resp.Reply {
	db.Flush()

	data := utils.ToCmdLine2("flushdb", args...)
	db.addAof(data)
	return reply.MakeOkReply()
}

//TYPE KEY
func execType (db *DB, args [][]byte) resp.Reply{
	key := args[0]
	entity, exists := db.GetEntity(string(key))
	if exists == false{
		return reply.MakeStatusReply("none")
	}

	switch entity.Data.(type){
	case []byte: //因为字符串的底层是用[]byte存的
		return reply.MakeStatusReply("string")
	}

	return &reply.UnknownErrReply{}
}

func execRename (db *DB, args [][]byte) resp.Reply{
	src := string(args[0])
	dest := string(args[1])

	entity, exists := db.GetEntity(src)
	if exists == false{
		return reply.MakeStatusReply("none")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)

	data := utils.ToCmdLine2("rename", args...)
	db.addAof(data)

	return reply.MakeOkReply()
}

func execRenamenx (db *DB, args [][]byte) resp.Reply{
	src := string(args[0])
	dest := string(args[1])

	entity, exists := db.GetEntity(src)
	if exists == false{
		return reply.MakeStatusReply("none")
	}
	_, exists = db.GetEntity(dest)
	if exists == true{
		return reply.MakeIntReply(0)
	}

	db.PutEntity(dest, entity)
	db.Remove(src)

	data := utils.ToCmdLine2("renamenx", args...)
	db.addAof(data)

	return reply.MakeIntReply(1)
}

func execKeys (db *DB, args [][]byte) resp.Reply{
	pattern := wildcard.CompilePattern(string(args[0])) //pattern用于实现通配符逻辑的

	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		match := pattern.IsMatch(key) //这行表示如果当前key符合通配符逻辑，那么match就为true，否则match为false
		if match == true{
			result = append(result, []byte(key))
		}
		return true
	})

	return reply.MakeMultiBulkReply(result)
}


func init(){
	RegisterCommand("DEL", execDel, -2) //del k1, k2...
	RegisterCommand("EXISTS", execExists, -2) //exist k1, k2...
	RegisterCommand("FLUSHDB", execFlushDB, 1) //flush
	RegisterCommand("TYPE", execType, 2) //type key
	RegisterCommand("RENAME", execRename, 3) //rename k1, k2
	RegisterCommand("RENAMENX", execRenamenx, 3) //renamenx k1, k2
	RegisterCommand("KEYS", execKeys, 2) //keys *
}