package main

import (
	"fmt"
	"my_redis/config"
	"my_redis/lib/logger"
	"my_redis/resp/handler"
	"my_redis/tcp"
	"os"
)


const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "127.0.0.1",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

//测试命令
//*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
//*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n
//*2\r\n$6\r\nselect\r\n$1\r\n1\r\n

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d",
				config.Properties.Bind,
				config.Properties.Port),
		},
		handler.MakeHandler())
	if err != nil {
		logger.Error(err)
	}
}
