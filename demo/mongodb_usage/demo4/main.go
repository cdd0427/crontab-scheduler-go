package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
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

//jobName过滤条件
type findByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		client *mongo.Client
		dataBase *mongo.Database
		collection *mongo.Collection
		cur mongo.Cursor
		cond *findByJobName
		record *LogRecord
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

	//4.定义查询条件
	cond=&findByJobName{JobName:"job10"}

	//5.查询
	if cur,err=collection.Find(context.TODO(),cond,findopt.Skip(0),findopt.Limit(2));err!=nil{
		fmt.Println(err)
		return
	}

	//6.遍历游标
	//反序列化
	for cur.Next(context.TODO()){
		record=&LogRecord{}
		if err=cur.Decode(record);err!=nil{
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}
}