package service

import (
	"github.com/gin-gonic/gin"
	"im/models"
	"im/tools"
	"log"
	"net/http"
)

// Login 登录函数
func Login(ctx *gin.Context) {
	//判断前端是否传递了账户和密码
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if account == "" || password == "" { //如果用户或密码为空，则直接返回提示
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}

	//通过用户名和密码来获取用户数据（后续的密码我们会是使用MD5的盐值加密进行处理，现在明文处理）
	//_, err := models.GetUserBasicByAccountPassword(account, password)
	//对密码进行盐值加密处理
	ub, err := models.GetUserBasicByAccountPassword(account, tools.GetMd5(password))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
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

	ctx.JSON(http.StatusOK, gin.H{ //反之成功则提示登录成功，返回生成的token信息
		"code": 200,
		"msg":  "登陆成功",
		"data": gin.H{
			"token": token,
		},
	})
}

// UserDetail 获取用户详细信息
func UserDetail(ctx *gin.Context) {
	tokenUserClaims, _ := ctx.Get("user_claims") //之前authCheck中解析出来的userClaims放到了gin.context中，现在把他拿出来(是否存在不处理，我知道她一定在)
	//返回的是一个接口类型，我们需要强制转换为userClaim类型，然后拿到里面的Identity，用她来获取用户的详细信息
	userClaims := tokenUserClaims.(*tools.UserClaims)
	//通过userClaims中的Identity来获取数据
	userBasic, err := models.GetUserBasicByIdentity(userClaims.Identity)
	if err != nil {
		log.Print("DB_ERROR")
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": userBasic,
	})
}

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

	//发送验证码(后续我们会使用别人写好的接口来生成验证码)
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
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发送验证码成功",
	})

}
