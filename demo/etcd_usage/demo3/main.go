package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		err     error
		getResp *clientv3.GetResponse
	)
	//生成配置
	config = clientv3.Config{
		Endpoints:   []string{"47.94.201.24:2379"},
		DialTimeout: 5 * time.Second,
	}
	//建立一个client
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	//生成用于读写的键值对
	kv = clientv3.NewKV(client)
	//
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Count)
	}

}
