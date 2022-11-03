package service

import (
	"github.com/gin-gonic/gin"
	"im/models"
	"im/tools"
	"log"
	"net/http"
)

// UserDetail 获取用户详细信息（获取自己的）
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

// UserQuery 获取用户基本信息（获取他人）
func UserQuery(ctx *gin.Context) {
	acc := ctx.Query("account")
	if acc == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	//走的和UserDetail一个接口，后面我们从新做一个数据返回
	uc, err := models.GetUserBasicByAccount(acc)
	if err != nil {
		log.Printf("获取用户数据失败：%v", err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}

	//拿到当前用户的user_claims，里面有我们需要的user_identity
	uClaims := ctx.MustGet("user_claims").(*tools.UserClaims)
	//拿到UserBasic之后,把需要的信息从新做成一个UserQueryResult结构，然后返回
	data := models.UserQueryResult{
		Nickname: uc.NickName,
		Sex:      uc.Sex,
		Email:    uc.Email,
		Avatar:   uc.Avatar,
		IsFriend: false,
	}

	//调用判断是否是好友
	isFriend := models.IsFriend(uClaims.Identity, uc.Identity)
	if isFriend {
		data.IsFriend = isFriend
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": data,
	})
}
