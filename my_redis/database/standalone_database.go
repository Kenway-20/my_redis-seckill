package database

import (
	"my_redis/aof"
	"my_redis/config"
	"my_redis/interface/resp"
	"my_redis/lib/logger"
	"my_redis/resp/reply"
	"strconv"
	"strings"
)

type StandaloneDatabase struct {
	dbSet []*DB
	aofHandle *aof.AofHandler
}

func NewStandaloneDatabase() *StandaloneDatabase {
	database := &StandaloneDatabase{dbSet: make([]*DB, config.Properties.Databases)}
	for i, _ := range database.dbSet{
		db := MakeDB()
		db.index = i
		database.dbSet[i] = db
	}

	//如果aof选项是打开的话
	if config.Properties.AppendOnly{
		handler, err := aof.NewAofHandler(database)
		if err != nil{
			logger.Info(err)
			panic(err)
		}
		database.aofHandle = handler

		for _, db := range database.dbSet{
			temple := db //用于解决闭包，导致局部变量db逃逸到堆上问题
			temple.addAof = func(line CmdLine) {
				database.aofHandle.AddAof(temple.index, line)
			}
		}
	}



	return database
}

//set key val
//select 2
func (database *StandaloneDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		err := recover()
		if err != nil{
			logger.Info(err)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select"{
		if len(args) != 2{
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, database, args[1:])
	}

	dbIndex := client.GetDBIndex()

	return database.dbSet[dbIndex].Exec(client, args)
}

func (database *StandaloneDatabase) Close() {
}

func (database *StandaloneDatabase) AfterClientClose(c resp.Connection) {
}

//用户选择哪个数据库的指令, select 2
func execSelect(c resp.Connection, database *StandaloneDatabase, args [][]byte) resp.Reply{
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil{
		return reply.MakeStandardErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet){
		return reply.MakeStandardErrReply("ERR DB index is out of range")
	}

	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}



