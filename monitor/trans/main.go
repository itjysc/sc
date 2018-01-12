package main

import (
	"net"
	"log"
	"bufio"
	"github.com/Shopify/sarama"
)

func Handle(conn net.Conn,ch chan<- *sarama.ProducerMessage) {  //注意，"chan<-"表示只写channel。"<-chan"表示只读channel。
	defer conn.Close()
	r := bufio.NewReader(conn)
	for  {
		line,err := r.ReadString('\n')
		if err != nil {
			log.Print(err)
		}
		if len(line) == 0 {
			continue
		}
		line = line[:len(line)-1]
		message := &sarama.ProducerMessage{
			Topic:"falcon",
			Key:nil,
			Value:sarama.StringEncoder(line),
		}
		ch <- message  //将数据丢给kafka。


	}
}
func main()  {
	listener,err := net.Listen("tcp","59.110.12.72:6001")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	producer,err := sarama.NewAsyncProducer([]string("59.110.12.72:9092"),nil)  //定义一个生产者，需要获取kafka的地址和参数。
	if err != nil {
		log.Fatal(err)
	}
	ch := producer.Input()  //这个生产者可以产生一个channel。

	for  {
		conn,err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go Handle(conn,ch)
	}

}

