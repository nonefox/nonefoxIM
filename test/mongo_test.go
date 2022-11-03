package test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im/models"
	"testing"
	"time"
)

/*
测试链接mongoDB
*/
func TestFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//获取mongoDB的链接对象
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		//这里由于我们是配置的账户和密码连接，所以我们要配置相应
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://192.168.88.106:27017"))
	if err != nil {
		t.Fatal(err)
	}
	//获取连接成功，我们就可以获取数据库对象然后对数据进行操作
	db := client.Database("im")
	ub := new(models.UserBasic) //实例化user_basic数据对象
	//对数据进行操作
	//db.Collection(models.UserBasic{}.CollectionName()).FindOne(context.Background(), bson.D{})//通过方法来获数据库名，方便后续程序的解耦
	err = db.Collection("user_basic").FindOne(context.Background(), bson.D{}).Decode(ub) //把从mongoDB里面查询到的数据，映射入我们定义好的结构中去
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v", ub)
}

/*
测试mongoDB查询所有数据
*/
func TestFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//获取mongoDB的链接对象
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		//这里由于我们是配置的账户和密码连接，所以我们要配置相应
		Username: "admin",
		Password: "admin",
	}).ApplyURI("mongodb://192.168.88.106:27017"))
	if err != nil {
		t.Fatal(err)
	}
	//获取连接成功，我们就可以获取数据库对象然后对数据进行操作
	db := client.Database("im")
	ub := make([]*models.UserBasic, 0) //实例化user_basic数据集对象

	allUser, err := db.Collection("user_basic").Find(context.Background(), bson.D{}) //查询出所有的用户
	for allUser.Next(context.Background()) {
		us := new(models.UserBasic) //实例化用户数据
		err := allUser.Decode(us)   //把查到的每一个用户数据映射到用户结构中
		if err != nil {
			t.Fatal(err)
		}
		ub = append(ub, us) //再把数据组装到数组结构中
	}
	for _, u := range ub {
		fmt.Printf("%v\n", u)
	}
}
