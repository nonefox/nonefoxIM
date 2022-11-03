package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"im/models"
	"im/tools"
	"log"
	"net/http"
	"time"
)

// SendCode 发送验证码
func SendCode(ctx *gin.Context) {
	//先获取用户邮箱
	email := ctx.PostForm("email")
	if email == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱不能为空",
		})
	}
	//通过用户的邮箱查询用户账户个数，如果超过0则，提示用户此邮箱已被人占用
	countByEmail, err := models.GetUserBasicCountByEmail(email)
	if err != nil {
		log.Printf("数据库错误：%v", err)
		return
	}
	if countByEmail > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱已被占用",
		})
	}

	//发送验证码(后续我们会使用写好的接口来生成验证码)
	//err = tools.SendCode(email, "123456")
	//发送验证码（使用generateCode方法生成code）
	code := tools.GenerateCode()
	err = tools.SendCode(email, code)
	if err != nil {
		log.Printf("发送验证码错误：%v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发送验证码失败",
		})
		return
	}

	//获取redis的连接，并且把验码和邮箱等信息发送到redis中，设置好过期时间
	err = models.RedisClient.Set(context.Background(), "TOKEN_"+email, code, time.Second*300).Err()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "设置验证信息失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发送验证码成功",
	})

}
