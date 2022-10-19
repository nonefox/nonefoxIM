package router

import (
	"github.com/gin-gonic/gin"
	"im/service"
)

// Router 路由类
func Router() *gin.Engine {
	r := gin.Default()
	r.POST("/login", service.Login)
	return r
}
