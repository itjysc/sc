package main

import (
	"github.com/BurntSushi/toml"
	"flag"
	"log"
)

type UserScriptConfig struct {
	Path string //定义自定义脚步的路径。
	Step int	//定义调用周期。
}

type SenderConfig struct {
	TransAddr string `toml:"trans_addr"`      		//注意首字母要大写，定义转发的地址，并定义用"trans_addr"关键字指定相应的IP地址。
	FlushInterval int `toml:"flush_interval"`		//定义刷新时间。
	MaxSleepTime int `toml:"max_sleep_time"`		//最大等待时间。
}


type config struct {
	Sender SenderConfig
	UserScript []UserScriptConfig `toml:"user_script"`
}

var   (
	configPath = flag.String("config","/yinzhengjie/golang/path/src/yinzhengjie/monitor/toml/config.toml","config path") /*
	定义默认的配置文件"config.toml"，然后它会在里面获取相应的参数信息。
	*/
	gcfg config
)

func main() {
	flag.Parse()
	_,err := toml.DecodeFile(*configPath,&gcfg) //将configPath配置文件的内容传给名为gcfg的config结构体。
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v",gcfg)
}