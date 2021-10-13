package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/userWeb/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("user")
	{
		UserRouter.GET("list", api.GetUserList)
	}
}
