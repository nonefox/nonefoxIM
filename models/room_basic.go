package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// RoomBasic 房间基本信息结构
type RoomBasic struct {
	Identity     string `bson:"identity"`
	Number       string `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	UserIdentity string `bson:"user_identity"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (RoomBasic) CollectionName() string {
	return "room_basic"
}

// InsertOneRoomBasic 插入新的房间基本信息
func InsertOneRoomBasic(roomBasic *RoomBasic) error {
	_, err := Mongo.Collection(RoomBasic{}.CollectionName()).InsertOne(context.Background(), roomBasic)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRoomBasic(roomIdentity string) error {
	_, err := Mongo.Collection(RoomBasic{}.CollectionName()).DeleteOne(context.Background(), bson.D{{"identity", roomIdentity}})
	if err != nil {
		log.Printf("删除房间基本信息失败:%v", err)
		return err
	}
	return nil
}
