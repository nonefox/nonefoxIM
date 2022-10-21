package middlewares

import (
	"github.com/gin-gonic/gin"
	"im/tools"
	"net/http"
)

// AuthCheck 验证用户token是否正确
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取用户token
		token := c.GetHeader("token")
		//解析用户token
		userClaims, err := tools.AnalyseToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户认证不通过",
			})
			return
		}
		//认证通过，我们就把解析出来的用户认证信息，set到context中去，方便其他服务获取
		c.Set("user_claims", userClaims)
		c.Next() //让后面的中间件继续执行
	}
}
