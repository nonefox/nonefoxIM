package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Mongo 配置一个数据库对象的全局变量
var Mongo = InitMongo()

func InitMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//获取mongoDB的链接对象
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		//这里由于我们是配置的账户和密码连接，所以我们要配置相应
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://192.168.88.100:27017"))
	if err != nil {
		log.Fatal(err)
	}
	//获取连接成功，我们就可以获取数据库对象然后对数据进行操作
	db := client.Database("im")
	return db
}
