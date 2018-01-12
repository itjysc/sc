/*
	这个包是专门处理网络数据发送的模块。
*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"yinzhengjie/monitor/common"
)

type Sender struct { //定义一个Sender结构体，接收Channel发来的数据，同时通过网络IP将数据发送出去。
	addr string              //定义需要发送数据到远端的IP地址
	ch   chan *common.Metric //接收到的数据。
}

func NewSender(addr string) *Sender { //Sender的构造方法。
	return &Sender{
		addr: addr,
		ch:   make(chan *common.Metric, 1000), //定义channel包的大小。
	}
}

func (s *Sender) connect() net.Conn { //定义断线重连的函数
	n := 100 * time.Millisecond
	for {
		conn, err := net.Dial("tcp", s.addr) //建立连接
		if err != nil {
			log.Print(err)
			time.Sleep(n)
			n = n * 2
			if n > time.Second*30 {
				n = time.Second * 30 //重试拨号，最大间隔时间不能超过30秒。
			}
			continue
		}
		return conn //如果拨号成功就把conn返回，终止循环。
	}
}

func (s *Sender) Start() { //用于从Channel中读取数据。
	conn := s.connect() //可以判断连接是否断开，如果断开就进行重连，如果没有断开连接就拿到这个连接。
	log.Printf("本端地址：%v,对端地址：%v\n", conn.LocalAddr(), conn.RemoteAddr())
	w := bufio.NewWriter(conn)                //引用bufio就是为了实现定时定量的发送数据。使得代码具有抗压性。
	ticker := time.NewTicker(time.Second * 5) //定义一个定时器
	for {
		select {
		case metric := <-s.ch: //如果有数据，就从从channle读取数据。
			buf, _ := json.Marshal(metric)
			_, err := fmt.Fprintf(w, "%s\n", buf) //将数据发送给conn
			if err != nil {
				conn.Close()
				conn = s.connect()
				w = bufio.NewWriter(conn)
				log.Print(conn.LocalAddr())
			}
		case <-ticker.C: //定时发送数据，我们的定时器是5秒钟，所以以下的代码会每5秒运行一次。
			err := w.Flush() /*将数据强行发送给conn，（w默认存储大小是4K,只要数据不达到4k，它就一直不会发送
			数据，会一直缓存到4k后才会发送。因此，我们这是强制性要求不管多少数据都需要发送出去。）*/
			if err != nil {
				conn.Close()
				conn = s.connect()
				w = bufio.NewWriter(conn)
				log.Print(conn.LocalAddr())
			}
		}
	}
}

func (s *Sender) Channel() chan *common.Metric { //把自己的Channle暴露出去，让别人可以给它发送数据。
	return s.ch
}
