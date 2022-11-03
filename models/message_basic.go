package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MessageBasic 消息基本结构
type MessageBasic struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
}

// CollectionName 获取room数据库名
func (MessageBasic) CollectionName() string {
	return "message_basic"
}

// InsertOneMessageBasic 把用户发布的每一条消息都保存起来
func InsertOneMessageBasic(mb *MessageBasic) error {
	_, err := Mongo.Collection(MessageBasic{}.CollectionName()).InsertOne(context.Background(), mb)
	if err != nil {
		log.Printf("保存消息失败：%v", err)
		return err
	}
	return nil
}

// GetMessageListByRoomIdentity 通过用户的roomIdengtity来获取对应房间的聊天记录
func GetMessageListByRoomIdentity(roomIdentity string, pSize, skip *int64) ([]*MessageBasic, error) {
	//定义一个存储消息的变量
	msgList := make([]*MessageBasic, 0)
	//查询消息表，获取所有的消息记录
	allMsg, err := Mongo.Collection(MessageBasic{}.CollectionName()).Find(context.Background(), bson.D{{"room_identity", roomIdentity}},
		&options.FindOptions{
			Limit: pSize,
			Skip:  skip,
			Sort:  bson.D{{"create_at", -1}}, //依据消息创建时间来排序
		})
	if err != nil {
		log.Printf("获取消息记录失败：%v", err)
		return nil, err
	}
	//把获取出来的消息放入我们定义好的msgList中
	for allMsg.Next(context.Background()) {
		mb := new(MessageBasic)
		err := allMsg.Decode(mb)
		if err != nil {
			log.Printf("msgList消息解析错误：%v", err)
			return nil, err
		}
		msgList = append(msgList, mb)
	}
	return msgList, nil

}
