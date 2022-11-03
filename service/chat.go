package service

import (
	"github.com/gin-gonic/gin"
	"im/models"
	"im/tools"
	"net/http"
	"strconv"
)

// ChatList 获取聊天记录列表
func ChatList(ctx *gin.Context) {
	//依据房间ID获取房间的聊天记录
	roomId := ctx.Query("room_identity")
	if roomId == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间id不能为空",
		})
		return
	}

	//获取到对应的房间ID之后，判断用户是否属于这个房间，属于就获取该房间的聊天记录，反之则拒绝
	uc := ctx.MustGet("user_claims").(*tools.UserClaims)
	_, err := models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, roomId)
	if err != nil { //通过查询用户房间关联表判断是否属于这个房间
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户不属于该房间",
		})
		return
	}

	//从前端获取分页数据（做好分页变量），这里错误就不再处理知道有数据
	pageIndex, _ := strconv.ParseInt(ctx.Query("page_index"), 10, 32)
	pageSize, _ := strconv.ParseInt(ctx.Query("page_size"), 10, 32)
	skip := (pageIndex - 1) * pageSize //位移

	//从聊天记录表查出对应的聊天信息
	msgList, err := models.GetMessageListByRoomIdentity(roomId, &pageSize, &skip)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取聊天记录失败" + err.Error(),
		})
		return
	}
	//把聊天记录的数据返回给前端显示
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "聊天信息获取成功",
		"msgList": gin.H{ //数据单独用一个map封装
			"data": msgList,
		},
	})

}
