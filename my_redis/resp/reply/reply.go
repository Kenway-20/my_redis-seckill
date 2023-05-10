package reply

import (
	"bytes"
	"my_redis/interface/resp"
	"strconv"
)

var(
	nullBulkReplyBytes =[]byte("$-1\r\n")
	CRLF = "\r\n"
)

type BulkReply struct {
	Arg []byte
}

type MultiBulkReply struct {
	Args [][]byte
}

type StatusReply struct {
	Status string
}

type IntReply struct {
	Code int64
}

type StandardErrReply struct {
	Status string
}

func MakeBulkReply(Arg []byte) *BulkReply{
	return &BulkReply{Arg: Arg}
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0{
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

func MakeMultiBulkReply(Args [][]byte) *MultiBulkReply{
	return &MultiBulkReply{Args: Args}
}

func (m *MultiBulkReply) ToBytes() []byte {
	argLen := len(m.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range m.Args{
		if arg == nil{
			buf.WriteString(string(nullBulkReplyBytes))
		}else{
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}

	return buf.Bytes()
}

func MakeStatusReply(Status string) *StatusReply{
	return &StatusReply{Status: Status}
}

func (s *StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

func MakeIntReply(Code int64) *IntReply{
	return &IntReply{Code: Code}
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

func MakeStandardErrReply(Status string) *StandardErrReply{
	return &StandardErrReply{Status: Status}
}

func (s *StandardErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func IsErrReply(reply resp.Reply) bool{
	return reply.ToBytes()[0] == '-'
}


type ErrorReply interface {
	Error() string
	ToBytes() []byte
}
