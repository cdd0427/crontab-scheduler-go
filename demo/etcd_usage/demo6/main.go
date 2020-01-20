package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		kv             clientv3.KV
		err            error
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
	//申请一个lease
	lease = clientv3.NewLease(client)
	//申请一个十秒的lease
	leaseGrantResp, err = lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	//拿lease ID
	leaseId = leaseGrantResp.ID

	////自动续租
	//ctx,_:=context.WithTimeout(context.TODO(),5*time.Second)

	//先续租五秒之后停止续租之后还有10秒
	if keepRespChan,err=lease.KeepAlive(context.TODO(),leaseId);err!=nil{
		fmt.Println(err)
		return
	}
	//处理续租应答
	go func() {
		for  {
			select {
			case keepResp=<-keepRespChan:
				if keepRespChan==nil{
					fmt.Println("租约已失效")
					goto END
				}else{
					fmt.Println("收到自动续租应答",keepResp.ID)
				}
			}
		}
		END:
	}()
	//拿kv对象
	kv = clientv3.NewKV(client)
	//put一个kv，与lease关联，从而实现一定时间后自动过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "none", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功", putResp.Header.Revision)
	//定时查看key是否过期
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("已过期")
			break
		}
		fmt.Println(time.Now(), "还没过期")
		time.Sleep(2 * time.Second)
	}
}
