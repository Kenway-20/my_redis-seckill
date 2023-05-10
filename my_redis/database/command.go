package database

import "strings"

var cmdTable = make(map[string]*command) //string表示指令，cmdTable是用来存整个系统中已经注册的指令，每个指令对应一个command结构体

type command struct {
	exector ExecFunc
	arity int
}

func RegisterCommand(name string, exector ExecFunc, arity int){
	name = strings.ToLower(name)
	cmdTable[name] = &command{exector: exector, arity: arity}
}
