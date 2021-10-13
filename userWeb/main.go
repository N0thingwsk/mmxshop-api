package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/userWeb/config"
	"mxshop-api/userWeb/initialize"
)

func main() {
	initialize.InitLogger()
	err := initialize.InitConfig()
	if err != nil {
		zap.S().Panic("配置文件读取失败" + err.Error())
	}
	router := initialize.Routers()
	zap.S().Info("启动服务器, 端口:", config.UserServerConfig.UserConfig.Port)
	if err := router.Run(fmt.Sprintf("%s:%d", config.UserServerConfig.UserConfig.Host, config.UserServerConfig.UserConfig.Port)); err != nil {
		zap.S().Panic("启动失败" + err.Error())
	}
}
