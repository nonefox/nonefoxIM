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
	// 用户注册
	r.POST("/register", service.Register)

	//auth分组下的路由都需要进行用户token验证
	auth := r.Group("/u", middlewares.AuthCheck())           //用户验证的分组路由
	auth.GET("/user/detail", service.UserDetail)             //获取用户详细信息
	auth.GET("user/query", service.UserQuery)                //获取用户基本信息
	auth.GET("/websocket/message", service.WebSocketMessage) //发送信息
	//r.GET("/websocket/message", service.WebSocketMessage)    //发送信息
	auth.GET("chat/list", service.ChatList)           //获取聊天记录
	auth.POST("/user/add", service.AddFriend)         // 添加用户
	auth.DELETE("/user/delete", service.DeleteFriend) // 删除好友
	return r
}
