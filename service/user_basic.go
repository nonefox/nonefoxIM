package service

import (
	"github.com/gin-gonic/gin"
	"im/models"
	"net/http"
)

// Login 登录函数
func Login(ctx *gin.Context) {
	//判断前端是否传递了账户和密码
	accunt := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if accunt == "" || password == "" { //如果用户或密码为空，则直接返回提示
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}

	//通过用户名和密码来获取用户数据（后续的密码我们会是使用MD5的盐值加密进行处理，现在明文处理）
	_, err := models.GetUserBasicByAccountPassword(accunt, password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
