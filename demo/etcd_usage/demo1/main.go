package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
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
	client = client
}