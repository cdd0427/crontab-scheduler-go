package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)


type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime int64	`bson:"endTime"`
}

type LogRecord struct {
	JobName string	`bson:"jobName"`
	Command string	`bson:"command"`
	Err string	`bson:"err"`
	Content string	`bson:"content"`
	TimePoint	`bson:"timePoint"`
}

type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client *mongo.Client
		dataBase *mongo.Database
		collection *mongo.Collection
		delCond *DeleteCond
		delResult *mongo.DeleteResult
		err error
	)
	//1.建立链接
	client,err=mongo.Connect(context.TODO(),"mongodb://47.94.201.24:27017",clientopt.ConnectTimeout(5*time.Second))
	if err!=nil{
		fmt.Println(err)
		return
	}
	//2.选择数据库
	dataBase=client.Database("cron")
	//3.选择表
	collection=dataBase.Collection("log")

	//4.删除开始时间早于当前时间的所有日志
	//delete({"timepoint.startTime":{"$lt":当前时间}})

	delCond=&DeleteCond{beforeCond:TimeBeforeCond{Before:time.Now().Unix()}}

	//执行删除
	if delResult,err=collection.DeleteMany(context.TODO(),delCond);err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println(delResult.DeletedCount)
}
