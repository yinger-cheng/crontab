package master

import (
	"encoding/json"
	"github.com/yinger-cheng/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

//任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
	err error
}

var (
	//单例对象
	G_apiServer *ApiServer
)
//保存任务接口
// POST job={"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(resp http.ResponseWriter,req *http.Request){
	var (
		err error
		postJob string
		job common.Job
		oldJob *common.Job
		bytes []byte
	)
	//1 解析post表单
	if err = req.ParseForm();err!= nil{
		goto ERR
	}
	//2 取表单中的job字段
	postJob = req.PostForm.Get("job")
	//3 反序列化job
	if err = json.Unmarshal([]byte(postJob),&job);err!=nil{
		goto ERR
	}
	//4 保存到etcd
	if oldJob,err = G_jobMgr.SaveJob(&job);err!= nil{
		goto ERR
	}
	//5 返回正常应答({"error":0,"msg":"；"data":{……}})
	if bytes,err = common.BuildResponse(0,"success",oldJob);err !=nil{
		resp.Write(bytes)
	}
	return
ERR:
	//6 返回异常应答
	if bytes,err = common.BuildResponse(-1,err.Error(),nil);err !=nil{
		resp.Write(bytes)
	}
}

func InitApiServer()(err error) {
	var (
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save",handleJobSave)
	//启动TCP监听
	if listener,err = net.Listen("tcp",":"+strconv.Itoa(G_config.APiPort));err!=nil{
		return
	}
	//创建一个http服务
	httpServer = &http.Server{
		ReadHeaderTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:mux,
	}
	//赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}
	//启动了服务端
	go httpServer.Serve(listener)

	return
}
