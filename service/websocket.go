package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im/define"
	"im/models"
	"im/tools"
	"log"
	"net/http"
	"time"
)

//定义一个默认的upgrader链接（这个包可以把我们的http链接升级为websocket链接）
var upgrader = websocket.Upgrader{}

//定义存放所有的链接的用户map，方便把数据发送给所有用户（自动去重）
var allConn = make(map[string]*websocket.Conn)

// WebSocketMessage 聊天室发布信息
func WebSocketMessage(ctx *gin.Context) {
	//建立一个websocket链接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "链接聊天室失败" + err.Error(),
		})
		return
	}
	defer conn.Close()

	//反之连接成功
	uc := ctx.MustGet("user_claim").(*tools.UserClaims) //获取用户声明信息（要发消息首先用户需要先登录）
	//把用户声明中的identity与用户的连接绑定，存入一个map数据中
	allConn[uc.Identity] = conn
	for {
		ms := new(define.MessageStruct) //定义一个消息结构
		//把读到的数据以json串的格式映射到我们的message结构中去（message中定义好了对应json标签）
		err := conn.ReadJSON(ms)
		if err != nil {
			log.Printf("读取信息失败:%v\n", err)
			return
		}

		//在发消息和收消息之前都要先确认该用户是否属于这个房间，如果不是这个房间就不能发送和接受信息
		_, err = models.GetUserRoomByUserIdentityRoomIdentity(uc.Identity, ms.RoomIdentity)
		if err != nil {
			log.Printf("用户ID：%v不存在房间ID：%v中", uc.Identity, ms.RoomIdentity)
			return
		}

		//保存消息（方便后续发送）
		mb := &models.MessageBasic{
			UserIdentity: uc.Identity,
			RoomIdentity: ms.RoomIdentity,
			Data:         ms.Message,
			CreateAt:     time.Now().Unix(),
			UpdateAt:     time.Now().Unix(),
		}

		//在把消息发送给关联用户之前先把消息存储起来
		err = models.InsertOneMessageBasic(mb)
		if err != nil {
			log.Printf("消息保存失败：%v", err)
			return
		}

		//获取特定房间的用户（然后把保存好的消息发送给该房间的所有用户）
		urs, err := models.GetUserRoomByRoomIdentity(ms.RoomIdentity)
		if err != nil {
			log.Printf("获取关联用户数据失败：%v", err)
			return
		}
		//把信息发送给所有的用户（我们把用户连接存到了一个map中）
		for _, userInRoom := range urs {
			cc := allConn[userInRoom.Identity] //存入map的时候UID作为key，conn作为value（用用户ID获取连接对象）
			err := cc.WriteMessage(websocket.TextMessage, []byte(ms.Message))
			if err != nil {
				log.Printf("写入信息失败:%v\n", err)
				return
			}
		}

	}
}
