package common

import "encoding/json"

type Job struct{
	Name string `json:"name"`//任务名
	Command string `json:"command"`//shell 命令
	CronExpr string `json:"cronExpr"`//cron表达式

}


//HTTP 接口应答
type Response struct{
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	data interface{} `json:"data"`
}
//应答方法
func BuildResponse(errno int,msg string,data interface{})(resp []byte,err error){
	//定义一个response
	var(
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.data = data
	//2 反序列化
	resp ,err = json.Marshal(response)
	return
}