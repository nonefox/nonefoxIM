package service

import (
	"github.com/gin-gonic/gin"
	"im/models"
	"im/tools"
	"net/http"
	"time"
)

// AddFriend 添加好友
func AddFriend(ctx *gin.Context) {
	//首先通过前端传过来的account信息获取，该用户的详细信息(方便后面更新用户关系)
	account := ctx.PostForm("account")
	if account == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数为空",
		})
		return
	}

	//获取用户详细信息
	ub, err := models.GetUserBasicByAccount(account)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取用户详细信息失败",
		})
		return
	}

	//获取当前用户的声明信息
	uc := ctx.MustGet("user_claims").(*tools.UserClaims)
	//先判断两个用户是否原本就是好友
	isFriend := models.IsFriend(ub.Identity, uc.Identity)
	if isFriend {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "互为好友，不可重复添加",
		})
		return
	}

	//不是好友，我们就先把房间记录保存起来（因为判断好友关系是看是否有单独聊天房间关系）
	rb := &models.RoomBasic{
		Identity:     tools.GetUUID(),
		UserIdentity: uc.Identity, //当前用户加好友，所以创建的房间是当前用户的identity
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = models.InsertOneRoomBasic(rb)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "插入房间基本信息失败",
		})
		return
	}

	//然后分别保存两个用户关于这个房间的关系，这样当两个用户互相查询对方的好友关系时就都可以查到
	ur1 := &models.UserRoom{
		Identity:     ub.Identity,
		RoomIdentity: rb.Identity,
		RoomType:     0,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = models.InsertOneUserRoom(ur1)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "插入用户房间关系信息失败",
		})
		return
	}
	ur2 := &models.UserRoom{
		Identity:     uc.Identity,
		RoomIdentity: rb.Identity,
		RoomType:     0,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	err = models.InsertOneUserRoom(ur2)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "插入用户房间关系信息失败",
		})
		return
	}

	//添加成功
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
	})
}

// DeleteFriend 删除好友
func DeleteFriend(ctx *gin.Context) {
	//获取前端带过来的好友的identity
	identity := ctx.Query("identity")
	if identity == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	//获取当前用户的用户声明信息
	uc := ctx.MustGet("user_claims").(tools.UserClaims)
	//还是要判断一下是否真的具有好友关系（防止前端传错 ^_^）
	isFriend := models.IsFriend(identity, uc.Identity)
	if !isFriend {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "不是好友，无需删除",
		})
		return
	}

	//根据两个好友用户的identity信息，查询用户关系表，获取对应的单独聊天房间ID
	roomIdentity, err := models.GetUserRoomIdentity(identity, uc.Identity)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查找单独聊天房间ID失败",
		})
		return
	}

	//删除对应用户房间关系
	err = models.DeleteUserRoom(roomIdentity)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除好友关系失败",
		})
		return
	}

	//删除对应好友关系之后，再删除对应房间基本信息
	err = models.DeleteRoomBasic(roomIdentity)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除房间基本信息失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除好友成功",
	})
}
