package aof

import (
	"io"
	"my_redis/config"
	databaseface "my_redis/interface/database"
	"my_redis/lib/logger"
	"my_redis/lib/utils"
	"my_redis/resp/connection"
	"my_redis/resp/parser"
	"my_redis/resp/reply"
	"os"
	"strconv"
)
// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

//aof buffer管道的大小
const aofQueueSize = 1 << 16

//某个操作的信息
type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// AofHandler receive msgs from channel and write to AOF file
type AofHandler struct {
	database databaseface.Database
	aofChan  chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int //上一条aof指令写入的db号，用于实现不切换db时不额外写入select xx
}

func NewAofHandler(database databaseface.Database) (*AofHandler, error){
	handle := &AofHandler{}
	handle.aofFilename = config.Properties.AppendFilename //根据redis.conf读aof文件的名字
	handle.database = database
	//根据之前的aof文件，恢复状态
	handle.LoadAof()

	aofFile, err:= os.OpenFile(handle.aofFilename, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0600)
	if err != nil{
		return nil, err
	}

	handle.aofFile = aofFile
	//初始化channel
	handle.aofChan = make(chan *payload, aofQueueSize)
	//起协程去异步处理aof落盘
	go func() {
		handle.handleAof()
	}()


	return handle, nil
}

//把写操作放进channel里
func (a *AofHandler) AddAof(dbIndex int, cmd CmdLine){
	if config.Properties.AppendOnly == false || a.aofChan == nil{
		return
	}

	p := &payload{dbIndex: dbIndex, cmdLine: cmd}
	a.aofChan <- p
	return
}

//把channel里的数据放到文件里(落盘)
func (a *AofHandler) handleAof(){
	for p := range a.aofChan{
		//如果要切数据的话，额外加一条select xx命令
		if a.currentDB != p.dbIndex{
			args := utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))
			data := reply.MakeMultiBulkReply(args).ToBytes()
			_, err := a.aofFile.Write(data)
			if err != nil{
				logger.Error(err)
				continue
			}
			a.currentDB = p.dbIndex
		}

		//写入实际命令
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := a.aofFile.Write(data)
		if err != nil{
			logger.Error(err)
			continue
		}
	}
}

//根据之前的aof文件，恢复状态
func (a *AofHandler) LoadAof(){
	open, err:= os.OpenFile(a.aofFilename, os.O_CREATE | os.O_RDONLY, 0600)
	if err != nil{
		logger.Info(err)
		panic(err)
	}
	defer open.Close()
	conn := &connection.Connection{}

	ch := parser.ParseStream(open)
	for p := range ch{
		if p.Err != nil{
			if p.Err == io.EOF{
				break
			}
			logger.Error(err)
			continue
		}
		if p.Data == nil{
			logger.Error("empty payload")
			continue
		}
		bulkReply, ok := p.Data.(*reply.MultiBulkReply)
		if ok == false{
			logger.Error("need multi bulk")
			continue
		}
		rep := a.database.Exec(conn, bulkReply.Args)
		if reply.IsErrReply(rep){
			logger.Error("exec err", rep.ToBytes())
		}

	}
}
