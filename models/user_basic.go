package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

//聊天室用户基本信息与数据库结构对应

type UserBasic struct {
	Identity string `bson:"_id"`
	Account  string `bson:"account"`
	Password string `bson:"password"`
	NickName string `bson:"nickName"`
	Sex      int    `bson:"sex"`
	Email    string `bson:"email"`
	Avatar   string `bson:"avatar"`
	CreateAt int64  `bson:"create_at"`
	UpdateAt int64  `bson:"update_at"`
}

// CollectionName 获取当前数据库名称
func (UserBasic) CollectionName() string {
	//这里我们直接写死在了方法里面，后续可以通过读取配置文件信息来获取
	return "user_basic"
}

// GetUserBasicByAccountPassword 通过前端提交的账户和密码，查询mongo中的数据
func GetUserBasicByAccountPassword(account, password string) (*UserBasic, error) {
	ub := new(UserBasic)
	//Mongo我们在Init包中定义的mmongoDB数据库的全局变量（所以这里就直接使用）
	err := Mongo.Collection("user_basic").FindOne(context.Background(),
		bson.D{{"account", account}, {"password", password}}).Decode(ub) //把从mongoDB里面查询到的数据，映射入我们定义好的结构中去
	if err != nil {
		log.Fatal(err)
	}
	return ub, err
}

// GetUserBasicByIdentity 通过用户的Identity来获取用户
func GetUserBasicByIdentity(identity string) (*UserBasic, error) {
	ub := new(UserBasic)
	//Mongo我们在Init包中定义的mmongoDB数据库的全局变量（所以这里就直接使用）
	err := Mongo.Collection("user_basic").FindOne(context.Background(),
		bson.D{{"identity", identity}}).Decode(ub) //把从mongoDB里面查询到的数据，映射入我们定义好的结构中去
	if err != nil {
		log.Fatal(err)
	}
	return ub, err
}

// GetUserBasicByAccount 通过用户的account来获取用户
func GetUserBasicByAccount(account string) (*UserBasic, error) {
	ub := new(UserBasic)
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{"account", account}}).
		Decode(ub)
	return ub, err
}

// GetUserBasicCountByEmail 通过email来获取该邮箱注册用户的数量
func GetUserBasicCountByEmail(email string) (int64, error) {
	count, err := Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{{"email", email}})
	return count, err
}
