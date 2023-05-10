package reply

type PongReply struct {}
var pongBytes = []byte("+PONG\r\n")

type OkReply struct {}
var okBytes = []byte("+OK\r\n")

type NullBulkReply struct {}
var nullBulkBytes = []byte("$-1r\n")

type EmptyMultiBulkReply struct {}
var emptyMultiBulkBytes = []byte("*0\r\n")

type NoReply struct {}
var noBytes = []byte("")

func (p *PongReply) ToBytes() []byte{
	return pongBytes
}

func MakePongReply() *PongReply{
	return &PongReply{}
}

func (o *OkReply) ToBytes() []byte{
	return okBytes
}

func MakeOkReply() *OkReply{
	return &OkReply{}
}

func (n *NullBulkReply) ToBytes() []byte{
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply{
	return &NullBulkReply{}
}

func (e *EmptyMultiBulkReply) ToBytes() []byte{
	return emptyMultiBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply{
	return &EmptyMultiBulkReply{}
}

func (n *NoReply) ToBytes() []byte{
	return noBytes
}

func MakeNoReply() *NoReply{
	return &NoReply{}
}