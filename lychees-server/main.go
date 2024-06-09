package main

import (
	"lychees-server/logs"
	"lychees-server/middlewares"
	"lychees-server/router"
)

func main() {
	logs.Logger.Infof("Version:%s ", middlewares.VERSION)

	router.Start()
}
