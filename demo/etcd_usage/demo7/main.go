package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main(){
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watcher clientv3.Watcher
		watchRespChan <-chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event
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

	//创建一个kv
	kv= clientv3.NewKV(client)

	go func() {
		for{
			_, _ = kv.Put(context.TODO(), "/cron/jobs/job7", "7")
			_, _ = kv.Delete(context.TODO(), "/cron/jobs/job7")
			time.Sleep(1*time.Second)
		}
	}()

	//先GET到当前值在监听后续的变化
	if getResp,err=kv.Get(context.TODO(),"/cron/jobs/job7");err!=nil{
		fmt.Println(err)
		return
	}
	//若k存在
	if len(getResp.Kvs)!=0{
		fmt.Println("当前值为",getResp.Kvs[0].Value)
	}
	//监听后续变化
	watchStartRevision=getResp.Header.Revision+1

	//创建监听器
	watcher=clientv3.NewWatcher(client)

	//启动监听
	fmt.Println("从该版本开始监听：",watchStartRevision)

	watchRespChan=watcher.Watch(context.TODO(),"/cron/jobs/job7",clientv3.WithRev(watchStartRevision))

	//处理kv变化事件
	for watchResp=range watchRespChan{
		for _,event= range watchResp.Events{
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：",string(event.Kv.Value),"Revision:",event.Kv.CreateRevision," ",event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了Revision:",event.Kv.ModRevision)

			}
		}
	}


}