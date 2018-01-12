package main

import (
	"flag"
	//"os"
	"runtime"
	"time"
	"yinzhengjie/monitor/common"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"strings"
	"os/exec"
	"log"
	"bufio"
	"strconv"
	"fmt"
	"github.com/BurntSushi/toml"
)

//var (
//	transAddr = flag.String("trans", "59.110.12.72:6000", "transfer address")
//)

func NewMetric(metric string, value float64) *common.Metric {  //构造出一个辅助函数，可以简化我们的代码。
	//hostname, _ := os.Hostname()
	hostname := "尹正杰"
	return &common.Metric{
		Metric:    metric,
		Endpoint:  hostname,
		Value:     value,
		Tag:       []string{runtime.GOOS},
		Timestamp: time.Now().Unix(),
	}
}

func CpuMetric() []*common.Metric {
	var ret []*common.Metric
	cpus, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}
	metric := NewMetric("cpu.usage", cpus[0])  //调用cpu的指标
	ret = append(ret, metric)

	cpuload, err := load.Avg()
	if err == nil {
		metric = NewMetric("cpu.load1", cpuload.Load1)   //采集一分钟内cpu的Load指标，当然只会采集到linux的load，在windows上是没有load指标的。
		ret = append(ret, metric)
		metric = NewMetric("cpu.load5", cpuload.Load5)  //采集五分钟内cpu的Load指标，当然只会采集到linux的load，在windows上是没有load指标的。
		ret = append(ret, metric)
		metric = NewMetric("cpu.load15", cpuload.Load15)  //采集十五分钟内cpu的Load指标，当然只会采集到linux的load，在windows上是没有load指标的。
		ret = append(ret, metric)
	}
	return ret
}

func GetUserMetrics(cmdstr string) ([]*common.Metric,error) { //采集到用户[]*common.Metric参数。但是还多了一个error尾巴，我们需要甩掉这个尾巴。
	var ret []*common.Metric
	cmd := exec.Command("bash","-c",cmdstr) //cmd就是我们拿到的bash命令。这一步是在构建命令。
	stdout,_ := cmd.StdoutPipe() //获取标准输出
	err := cmd.Start() //启动命令，不建议用Start().
	if err != nil {
		return nil,err
	}
	r := bufio.NewReader(stdout) //将输出的内容进行包装起来，便于我们操作。
	for  {
		line,err := r.ReadString('\n') //按照换行接受参数。
		if err != nil {
			log.Printf("read %s error:%v", cmdstr, err)
			break
		}
		log.Printf("line:%v", line)
		line = strings.TrimSpace(line) //脱去空格和换行符
		fields := strings.Fields(line) //获取到key和value。
		if len(fields)!= 2 { //如果获取到的长度不是2就跳过本次循环。
			continue
		}
		log.Printf("fields:%v", fields)
		key,value := fields[0],fields[1] //将获取到的参数赋值给key和value.
		n,err := strconv.ParseFloat(value,64)   //将字符串转换成浮点型。
		if err != nil {
			log.Print(err)
			continue
		}
		metric := NewMetric(key,n) //将key到采集指标传给metric
		fmt.Println("1111",metric)
		ret = append(ret,metric)
	}
	fmt.Println("22222",ret)
	return ret,nil
}

func NewUserMetric(cmdstr string)MetricFunc  { //利用闭包函数，将GetUserMetrics到小尾巴error处理掉，进行一下格式化，哈哈。
	return func() []*common.Metric {
		metrics,err := GetUserMetrics(cmdstr)
		if err != nil {
			log.Print(err)
			return []*common.Metric{} //如果GetUserMetrics函数出错，我们就返回空的MetricFunc。
		}
		//fmt.Println(metrics)
		return metrics  //如果GetUserMetrics没有出错，就只返回MetricFunc函数。
	}
}



func main() {
	flag.Parse()   //解析命令行参数
	_,err := toml.DecodeFile(*configPath,&gcfg)  //将configPath配置文件的内容传给名为gcfg的config结构体。
	if err != nil {
		log.Fatal(err)
	}
	sender := NewSender(gcfg.Sender.TransAddr)
	ch := sender.Channel()

	sched := NewSched(ch)   //构造一个调度器

	go  sched.AddMetric(CpuMetric, time.Second*2) //定义调度器的调度周期，表示两秒钟采集一次数据。
	fmt.Println(gcfg.UserScript)
	for _,ucfg := range gcfg.UserScript	{
		fmt.Println(ucfg.Path,ucfg.Step)
		go sched.AddMetric(NewUserMetric(ucfg.Path),time.Duration(ucfg.Step)*time.Second)
	}

	// memory, time.Second * 3
	// disk, time.Minute
	sender.Start()
}
