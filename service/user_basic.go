package service

import (
	"github.com/gin-gonic/gin"
	"im/models"
	"im/tools"
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
	//_, err := models.GetUserBasicByAccountPassword(accunt, password)
	//对密码进行盐值加密处理
	ub, err := models.GetUserBasicByAccountPassword(accunt, tools.GetMd5(password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}

	//在用户登录时生成用户token信息
	token, err := tools.GenerateToken(ub.Identity, ub.Email)
	if err != nil { //若生成Token失败则返回内部错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "系统内部错误" + err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{ //成功则提示登录成功，返回生成的token信息
		"code": 200,
		"msg":  "登陆成功",
		"data": gin.H{
			"token": token,
		},
	})
}
