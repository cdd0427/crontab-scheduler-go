package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		kvpair  mvccpb.KeyValue
		err     error
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

	//删除kv
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	}
	if delResp.PrevKvs != nil {
		for _, kvpair = range delResp.PrevKvs {
			fmt.Println("删除了 ", string(kvpair.Key), ":", string(kvpair.Value))
		}
	}
}
