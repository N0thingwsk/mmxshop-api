package main

import (
	"go.uber.org/zap"
	"mxshop-api/userWeb/initialize"
)

func main() {
	//port := 8021
	initialize.InitLogger()
	router := initialize.Routers()
	zap.S().Info("启动服务器, 端口:", 8021)
	if err := router.Run(":8021"); err != nil {
		zap.S().Panic("启动失败" + err.Error())
	}
}
