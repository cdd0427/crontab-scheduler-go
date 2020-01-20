package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	//expiration of locks with lease
	//op action
	//txn transaction:if else then

	var (
		config         clientv3.Config
		client         *clientv3.Client
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		leaseId        clientv3.LeaseID
		ctx            context.Context
		cancelFunc     context.CancelFunc
		kv             clientv3.KV
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
		err            error
	)
	//client config
	config = clientv3.Config{
		Endpoints:   []string{"http://47.94.201.24:2379"},
		DialTimeout: 5 * time.Second,
	}
	//bulid link
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//set lock
	//create a lease,automatic lease renewal,catch a key with lease

	lease = clientv3.NewLease(client)
	//申请一个十秒的lease
	leaseGrantResp, err = lease.Grant(context.TODO(), 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	//拿lease ID
	leaseId = leaseGrantResp.ID

	//准备一个用于取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())

	//确保函数推出后自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	//先续租五秒之后停止续租之后还有10秒
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}
	//处理续租应答
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已失效")
					goto END
				} else {
					fmt.Println("收到自动续租应答", keepResp.ID)
				}
			}
		}
	END:
	}()

	//抢key的逻辑：if 不存在key ，then设置它，else抢锁失败

	kv = clientv3.NewKV(client)

	//创建事务
	txn = kv.Txn(context.TODO())

	//定义事务

	//如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "lock", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9")) //否则抢锁失败

	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	//判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("555,没抢到,被占用：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	//exe tansaction
	fmt.Println("处理任务嘻嘻")
	time.Sleep(5 * time.Second)
	//release lock
	//cancel automatic lease,release lease
	//上面的两个derfer会把锁释放掉
}
