package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

func main() {
	var (
		client *mongo.Client
		dataBase *mongo.Database
		collection *mongo.Collection
		err error
	)
	//1.建立链接
	client,err=mongo.Connect(context.TODO(),"mongodb://47.94.201.24:27017",clientopt.ConnectTimeout(5*time.Second))
	if err!=nil{
		fmt.Println(err)
		return
	}
	//2.选择数据库
	dataBase=client.Database("my_db")
	//3.选择表
	collection=dataBase.Collection("my_collection")

	collection=collection
}
