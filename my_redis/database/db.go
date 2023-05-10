package database

import (
	"my_redis/datastruct/dict"
	"my_redis/interface/database"
	"my_redis/interface/resp"
	"my_redis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data dict.Dict
	addAof func(line CmdLine)
}

//执行用户的指令，例如set key-val、ping
type ExecFunc func(db *DB, args [][]byte) resp.Reply // 所有的redis执行函数都要写成这个形式，传db和指令数组，返回一个reply

type CmdLine [][]byte

func MakeDB() *DB{
	return &DB{
		data:dict.MakeSyncDict(),
		addAof: func(line CmdLine) {}, //用于保证恢复数据阶段的时候db里的addAof也有值
	}
}

func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply{
	//cmdName是PING、SET
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName] //这个map不用考虑并发安全问题，因为只在初始化的时候去写，平时使用都是读
	if ok == false{
		return reply.MakeStandardErrReply("ERR Unknown command " + cmdName)
	}

	if validateArity(cmd.arity, cmdLine) == false{
		return reply.MakeArgNumErrReply(cmdName)
	}

	fun := cmd.exector
	//cmdLine是 SET k v，cmdLine[1:]是去掉了SET，只保留k v
	return fun(db, cmdLine[1:])

}

func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

/* ---- data Access ----- */

// GetEntity returns DataEntity bind to given key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {

	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)

	return entity, true
}

// PutEntity a DataEntity into DB
func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists edit an existing DataEntity
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

// PutIfAbsent insert an DataEntity only if the key not exists
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

// Remove the given key from db
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush clean database
func (db *DB) Flush() {
	db.data.Clear()
}

