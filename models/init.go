package models

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Mongo 配置一个MongoDB数据库对象的全局变量
var Mongo = InitMongo()

// RedisClient 配置一个Redis数据库连接的全局变量
var RedisClient = InitRedis()

// InitMongo 获取mongoDB的连接对象
func InitMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//获取mongoDB的链接对象
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		//这里由于我们是配置的账户和密码连接，所以我们要配置相应
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://192.168.92.174:27017"))
	if err != nil {
		log.Fatal(err)
	}
	//获取连接成功，我们就可以获取数据库对象然后对数据进行操作
	db := client.Database("im")
	return db
}

// InitRedis 获取redis的连接对象
func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "192.168.92.174:6379",
	})
}
