package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct{
	APiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"ApiReadTimeout"`
	ApiWriteTimeout int `json:"ApiWriteTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
}

var (
	//单例
	G_config *Config
)

func InitConfig(filename string) (err error){
	var (
		content []byte
		conf Config
	)
	//1, 把配置文件读进来
	if content , err = ioutil.ReadFile(filename);err != nil{
		return
	}
	//2 做json反序列化
	if err = json.Unmarshal(content,&conf);err != nil{
		return
	}

	//3 单例赋值
	G_config = &conf
	return
}