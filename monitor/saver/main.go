package main

import (
	"github.com/bsm/sarama-cluster"
	"log"
	"time"
	"context"
)

func IndexName() string{ //格式化时间
	date := time.Now().Format("20060102 030405")
	return "falcon-" + date
}



func main() {
	consumer,err := cluster.NewConsumer([]string{"59.110.12.72:9092"}, //定义kafka的列表
		"falcon-saver", //consumer的guoup的名字
				[]string{"falcon"}, //topic的名字
				&cluster.Config{},  //定义默认配置
	)
	if err !=nil {
		log.Fatal(err)
	}

	esclient,err := elastic.NewClient(elastic.SetURL("http://59.110.12.72:9200")) //需要把es的地址写进去进行拨号链接。
	if err != nil {
		log.Fatal(err)
	}


	for  {
		select {
		case msg := <- consumer.Messages(): //读取"Messages"的channel
			_,err := esclient.Index().
					Index(IndexName()).
					Type("falcon").
					BodyString(string(msg.Value)).
					Do(context.TODO())
			if err != nil {
				log.Print(err)
			}
			log.Print(msg.Value)
		case err := <- consumer.Errors():  //读取"Errors"的channel
			log.Print(err)
		}
	}
}