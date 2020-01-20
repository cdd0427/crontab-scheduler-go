package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main(){
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		putOp clientv3.Op
		getOp clientv3.Op
		opResp clientv3.OpResponse
		err    error
	)
	//客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"http://47.94.201.24:2379"},
		DialTimeout: 5 * time.Second,
	}
	//建立链接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv=clientv3.KV(client)

	//建立op：operation
	putOp=clientv3.OpPut("/corn/jobs/job8","op")

	//exec op
	if opResp,err=kv.Do(context.TODO(),putOp);err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println("wrting revision:",opResp.Put().Header.Revision)

	//create get op
	getOp=clientv3.OpGet("/corn/jobs/job8")

	//exec op
	if opResp,err=kv.Do(context.TODO(),getOp);err!=nil{
		fmt.Println(err)
		return
	}
	//print
	fmt.Println("data revision",opResp.Get().Kvs[0].ModRevision)
	fmt.Println("data value", string(opResp.Get().Kvs[0].Value))
}
