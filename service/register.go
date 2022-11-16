package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"im/models"
	"im/tools"
	"net/http"
	"time"
)

// Register 用户注册
func Register(ctx *gin.Context) {
	//获取注册表单提交的数据
	code := ctx.PostForm("code")
	email := ctx.PostForm("email")
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if code == "" || email == "" || account == "" || password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
	}

	//判断账号是否已经被注册
	byAccount, err := models.GetUserBasicByAccount(account)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "内部错误",
		})
		return
	}
	if byAccount != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号已被注册",
		})
		return
	}

	//账号没有问题之后，判断验证码是否正确(当时存的时候，我们是把"TOKEN_"+email作为key，然后code作为value存储的)
	result, err := models.RedisClient.Get(context.Background(), "TOKEN_"+email).Result()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "内部错误",
		})
		return
	}
	if result != code {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}
	//验证码正确之后，我们就把新注册的用户信息存入用户表
	ub := &models.UserBasic{
		Identity:  tools.GetUUID(),
		Account:   account,
		Password:  tools.GetMd5(password),
		Email:     email,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	//存入用户表
	err = models.InsertOneUserBasic(ub)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册失败",
		})
		return
	}

	//新用户注册成功之后，为用户生成一个token
	token, err := tools.GenerateToken(ub.Identity, ub.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "token生成失败",
		})
		return
	}

	//反之成功则把token信息返给前端
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"token": token,
		},
	})
}
