package router

import (
	"github.com/gin-gonic/gin"
	"im/middlewares"
	"im/service"
)

// Router 路由类
func Router() *gin.Engine {
	r := gin.Default()
	r.POST("/login", service.Login) //用户登录

	// 发送验证码
	r.POST("/send/code", service.SendCode)

	//auth分组下的路由都需要进行用户token验证
	auth := r.Group("/u", middlewares.AuthCheck()) //用户验证的分组路由
	auth.GET("/user/detail", service.UserDetail)
	return r
}
