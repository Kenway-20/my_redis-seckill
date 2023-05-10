package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"my_redis/interface/resp"
	"my_redis/resp/reply"
	"runtime/debug"
	"strconv"
	"strings"
)

type Payload struct {
	Data resp.Reply
	Err error
}

type readState struct {
	readingMultiLine bool
	expectArgsCount int
	msgType byte
	args [][]byte
	bulkLen int64
}

func (r *readState) finished() bool{
	return r.expectArgsCount > 0 && len(r.args) == r.expectArgsCount
}

func ParseStream(reader io.Reader) <-chan *Payload{
	ch := make(chan *Payload, 0)
	go parse0(reader, ch)
	return ch
}

//核心方法
func parse0(reader io.Reader, ch chan<- *Payload){
	defer func() {
		err := recover()
		if err != nil{
			fmt.Println(debug.Stack())
			return
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for true{
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		//处理错误
		if err != nil{
			if ioErr == true{
				ch <-&Payload{Err: err}
				close(ch)
				return
			}
			ch<-&Payload{Err: err}
			state = readState{}
			continue
		}

		//判断目前是否处于多行解析模式
		if state.readingMultiLine == false{ //不处于多行模式
			if msg[0] == '*'{
				err := parserMultiBulkHeader(msg, &state)
				if err != nil{
					ch <- &Payload{Err: errors.New("protocol err: " + string(msg))}
					state = readState{}
					continue
				}
				if state.expectArgsCount == 0{
					ch <- &Payload{Data: &reply.EmptyMultiBulkReply{}}
					state = readState{}
					continue
				}
			}else if msg[0] == '$'{ //$3\r\n
				err := parserBulkHeader(msg, &state)
				if err != nil{
					ch <- &Payload{Err: errors.New("protocol err: " + string(msg))}
					state = readState{}
					continue
				}
				if state.bulkLen == -1{
					ch <- &Payload{Data: &reply.EmptyMultiBulkReply{}}
					state = readState{}
					continue
				}
			}else{
				result, err := parserSingleLineReply(msg, &state)
				ch <- &Payload{Data: result, Err: err}
				state = readState{}
				continue
			}
		}else{ //处于多行模式
			err = readBody(msg, &state)
			if err != nil{
				ch <- &Payload{Err: errors.New("protocol err: " + string(msg))}
				state = readState{}
				continue
			}

			if state.finished(){
				var result resp.Reply
				if state.msgType == '*'{
					result = reply.MakeMultiBulkReply(state.args)
				}else if state.msgType == '$'{
					result = reply.MakeBulkReply(state.args[0])
				}

				ch <- &Payload{Data: result, Err: err}
				state = readState{}
			}
		}

	}
}

//从bufferIO里面读一行
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error){
	//逻辑是如果之前读到$x，那么严格读取x个字符，否则就默认按\r\n来划分
	var msg []byte
	var err error
	if state.bulkLen == 0{ //说明之前没有读到$x，这里按照\r\n切分
		msg, err = bufReader.ReadBytes('\n')
		if err != nil{
			return nil, true, err
		}

		if len(msg) == 0 || msg[len(msg)-2] != '\r'{
			return nil, false, errors.New("protocol err: " + string(msg))
		}
	}else{ //说明之前读到$x，严格读取x个字符
		msg = make([]byte, state.bulkLen + 2) //+2是因为还要多读\r\n
		_, err := io.ReadFull(bufReader, msg)
		if err != nil{
			return nil, true, err
		}

		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n'{
			return nil, false, errors.New("protocol err: " + string(msg))
		}

		state.bulkLen = 0 //已经读完了上一个x个字符，要重置状态
	}

	return msg, false, nil
}



//解析readLine里读出来的字符串，$3\r\nSET\r\n
func parserBulkHeader(msg []byte, state *readState) error{
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil{
		return errors.New("protocol err: " + string(msg))
	}
	if state.bulkLen == -1{
		return nil
	}else if state.bulkLen > 0{
		state.msgType = msg[0] //用来表示正在读什么样的数据，这里的msg[0]是'*'表示在读数组
		state.readingMultiLine = true
		state.expectArgsCount = 1
		state.args = make([][]byte, 0, 1)
	}else{
		return errors.New("protocol err: " + string(msg))
	}

	return nil
}

//解析readLine里读出来的多行数组
//*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
func parserMultiBulkHeader(msg []byte, state *readState) error{
	var err error
	var expectedLine int64 //存的是*3中的那个3
	expectedLine, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil{
		return errors.New("protocol err: " + string(msg))
	}
	if expectedLine == 0{
		state.expectArgsCount = 0
		return nil
	}else if expectedLine > 0{
		state.msgType = msg[0] //用来表示正在读什么样的数据，这里的msg[0]是'*'表示在读数组
		state.readingMultiLine = true
		state.expectArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
	}else{
		return errors.New("protocol err: " + string(msg))
	}
	return nil
}

//解析readLine里读出来的单行命令，:5\r\n、+OK\r\n、-err\r\n
func parserSingleLineReply(msg []byte, state *readState) (resp.Reply, error){
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch str[0] {
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil{
			return nil, errors.New("protocol err: " + string(msg))
		}
		result = reply.MakeIntReply(val)
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeStandardErrReply(str[1:])
	}

	return result, nil
}

//$3\r\nSET\r\n$3\r\nkey\r\n
//PING\r\n
func readBody(msg []byte, state *readState) error{
	line := msg[0: len(msg)-2]
	var err error
	////$3\r\nSET\r\n$3\r\nkey\r\n
	if line[0] == '$'{
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil{
			return errors.New("protocol err: " + string(msg))
		}
		//$0\r\n
		if state.bulkLen <= 0{
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	}else{
		state.args = append(state.args, line)
	}
	return nil
}
