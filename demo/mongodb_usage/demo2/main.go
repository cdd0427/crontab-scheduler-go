package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
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

func main() {
	var (
		client *mongo.Client
		dataBase *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		result *mongo.InsertOneResult
		docId objectid.ObjectID
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

	record=&LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{
			StartTime:time.Now().Unix(),
			EndTime:time.Now().Unix()+10,
		},
	}
	//插入记录{bson}
	if result,err=collection.InsertOne(context.TODO(),record);err!=nil{
		fmt.Println(err)
		return
	}
	//不传id就会默认生成一个全局id：12字节的二进制
	docId=result.InsertedID.(objectid.ObjectID)
	//转到16进制输出
	fmt.Println("自增ID：",docId.Hex())

}