package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// UserRoom 用户与房间的关联结构
type UserRoom struct {
	Identity     string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType     int    `bson:"room_type"` //[0:表示单独聊天，1:表示群体聊天]
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

// GetUserRoomByUserIdentityRoomIdentity 通过userIdentity和roomIdentity来获取用户对应的的房间
func GetUserRoomByUserIdentityRoomIdentity(userIdentity, roomIdentity string) (*UserRoom, error) {
	ur := new(UserRoom)
	//Mongo我们在Init包中定义的mmongoDB数据库的全局变量（所以这里就直接使用）
	err := Mongo.Collection(UserRoom{}.CollectionName()).FindOne(context.Background(),
		bson.D{{"user_identity", userIdentity}, {"room_identity", roomIdentity}}).Decode(ur) //把从mongoDB里面查询到的数据，映射入我们定义好的结构中去
	if err != nil {
		log.Fatal(err)
	}
	return ur, err
}

// GetUserRoomByRoomIdentity 通过房间ID查询user_room表，获取出该房间关联的所有用户关联信息
func GetUserRoomByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	//依据房间ID查询出先关联的所有的用户信息
	aur, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), bson.D{{"room_identity", roomIdentity}})
	if err != nil {
		log.Printf("获取房间关联用户信息错误：%v", err)
		return nil, err
	}
	allUsers := make([]*UserRoom, 0)     //用来存储关联用户信息的数组
	for aur.Next(context.Background()) { //依次拿出关联用户信息
		ur := new(UserRoom)
		err := aur.Decode(ur)
		if err != nil {
			return nil, err
		}
		//放入数组中
		allUsers = append(allUsers, ur)
	}
	return allUsers, nil
}

// IsFriend 判断两个用户是否是单独聊天的好友（只有拥有单聊房间才被称为好友，共同在一个群聊中不算好友，只算群友 ^_^ ）
func IsFriend(userIdentity1, userIdentity2 string) bool {
	//首先我们先查出当前用户所有参与单独聊天房间的roomIdentity
	finds, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{"user_identity", userIdentity1}, {"room_type", 0}})
	roomIds := make([]string, 0) //存放查出来的roomIdentity
	if err != nil {
		log.Printf("查询单聊房间失败：%v", err)
		return false
	}

	//把roomIdentity提取出来
	for finds.Next(context.Background()) {
		ur := new(UserRoom)
		//把查到的当前用户的单独聊天房间信息，映射到我们的userRoom结构中
		err := finds.Decode(ur)
		if err != nil {
			log.Printf("映射单独聊天userRoom数据事变：%v", err)
			return false
		}
		roomIds = append(roomIds, ur.RoomIdentity) //拿到所有的当前用户所有参与单独聊天房间的roomIdentity
	}

	//然后我们依据上面拿到的roomIds，作为条件来查询用户2是否包含这些房间identity（如果查出结果，那么他们就在一个单独聊天房间里，就是好友）
	num, err := Mongo.Collection(UserRoom{}.CollectionName()).
		CountDocuments(context.Background(), bson.M{"user_identity": userIdentity2, "room_type": 0, "room_identity": bson.M{"$in": roomIds}})
	if err != nil {
		log.Printf("查询关联单独聊天数据失败：%v", err)
		return false
	}
	if num > 0 {
		return true
	}
	return false
}

// InsertOneUserRoom 插入新的用户房间关系
func InsertOneUserRoom(userRoom *UserRoom) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).InsertOne(context.Background(), userRoom)
	if err != nil {
		log.Printf("插入户房间关系失败:%v", err)
		return err
	}
	return nil
}

// GetUserRoomIdentity 依据两个互为好友的identity，查询他们对应的单独聊天房间ID
func GetUserRoomIdentity(userIdentity1, userIdentity2 string) (string, error) {
	//首先获取user1在用户房间关系表中的单聊房间关系的房间ID
	results, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{"user_identity", userIdentity1}, {"room_type", 0}})
	if err != nil {
		log.Printf("user1查询户单独聊天房间ID失败:%v", err)
		return "", err
	}
	//定义一个切片存储查询出来的房间ID
	roomIds := make([]string, 0)
	for results.Next(context.Background()) {
		ur := new(UserRoom)
		err := results.Decode(ur)
		if err != nil {
			log.Printf("Decode用户房间关系失败:%v", err)
			return "", err
		}
		roomIds = append(roomIds, ur.RoomIdentity)
	}

	//把上面查出的user1的单独聊天房间id作为条件，查询user2的单独聊天房间id，如果有匹配的，那这个房间id对应的数据，就是我们要删除的好友关系
	ur := new(UserRoom)
	err = Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.M{"user_identity": userIdentity2, "room_type": 0, "room_identity": bson.M{"$in": roomIds}}).Decode(ur)
	if err != nil {
		log.Printf("user2查询户单独聊天房间ID失败:%v", err)
		return "", err
	}

	return ur.RoomIdentity, nil
}

// DeleteUserRoom 删除用户房间关系
func DeleteUserRoom(roomIdentity string) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).DeleteOne(context.Background(), roomIdentity)
	if err != nil {
		log.Printf("删除用户房间关系失败:%v", err)
		return err
	}
	return nil
}
