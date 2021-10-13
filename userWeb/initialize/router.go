package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/userWeb/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	apiGroup := r.Group("/v1")
	router.InitUserRouter(apiGroup)
	return r
}
