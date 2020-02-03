package master

import (
	"context"
	"encoding/json"
	"github.com/yinger-cheng/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//任务管理器
type JobMgr struct{
	client * clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var(
	//单例
	G_jobMgr *JobMgr
)
func InitJobMgr()(err error){
	var(
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	//初始化配置
	config = clientv3.Config{
		Endpoints: G_config.EtcdEndpoints,//集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout)*time.Millisecond,//连接超时
	}
	//建立连接
	if client,err = clientv3.New(config);err != nil{
		return
	}

	//得到KV和lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	//赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

func (JobMgr *JobMgr) SaveJob(job *common.Job)(oldJob *common.Job,err error){
	//把任务保存到/cron/jobs/任务名 ->json
	var (
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
		oldJobObj common.Job
	)
	jobKey = "/cron/jobs/"+job.Name
	if jobValue,err = json.Marshal(job);err != nil{
		return
	}
	//保存到etcd
	if putResp, err = JobMgr.kv.Put(context.TODO(),jobKey,string(jobValue),clientv3.WithPrevKV());err != nil{
		return
	}
	//如果是更新 那么返回旧值
	if putResp.PrevKv != nil{
		//对旧值做一个反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value,&oldJobObj);err != nil{
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}