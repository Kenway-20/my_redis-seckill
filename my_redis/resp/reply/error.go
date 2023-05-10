package reply

type UnknownErrReply struct {}
var unknownBytes = []byte("-Err unknown\r\n")

type ArgNumErrReply struct {
	cmd string
}

type SyntaxErrReply struct {}
var syntaxErrReplyBytes = []byte("-Err syntax error\r\n")

type WrongTypeErrReply struct {}
var wrongTypeErrReplyBytes = []byte("-Err wrong type error\r\n")

type ProtocolErrReply struct {
	msg string
}
var protocolErrReplyBytes = []byte("-Err proto col  error\r\n")


func (u *UnknownErrReply) Error() string {
	return "Err unknown"
}

func (u *UnknownErrReply) ToBytes() []byte {
	return unknownBytes
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-Err wrong number of arguments for '" + a.cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply{
	return &ArgNumErrReply{cmd: cmd}
}

func MakeSyntaxErrReply(cmd string) *SyntaxErrReply{
	return &SyntaxErrReply{}
}

func (s *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrReplyBytes
}

func (p *ProtocolErrReply) Error() string {
	return "-Err proto col error: '" + p.msg +"' \r\n"
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-Err proto col error: '" + p.msg +"' \r\n")
}

