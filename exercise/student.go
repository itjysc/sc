package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	cmd       string
	name      string
	id        int
	line      string
	file_name string
)

type Student struct {
	ID   int
	Name string
}

func main() {
	f := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(" 请输入>>> ")
		line, _ = f.ReadString('\n')
		fmt.Sscan(line, &cmd)
		if len(line) == 1 {
			continue
		}
		switch cmd {
		case "list":
			list()
		case "add":
			add()
		case "save":
			save()
		case "load":
			load()
		case "stop":
			os.Exit(0)
		default:
			fmt.Println("您输出的命令无效")
		}
	}
}

func list() {
	f, err := os.Open("student_info.json") //打开一个文件，如果这个文件不存在的话就会报错。
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(f) //取出文件的内容
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Print(line)
	}
	f.Close()
}

func add() {
	fmt.Sscan(line, &cmd, &id, &name)
	f, err := os.Open("student_info.json") //打开一个文件，如果这个文件不存在的话就会报错。
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(f) //取出文件的内容
	flag := 0               //定义一个标志位，当输入的ID和name相同时，就将其的值改为1.
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		s := Student{
			ID:   id,
			Name: name,
		}
		buf, err := json.Marshal(s)
		if err != nil {
			log.Fatal("序列化报错是：%s", err)
		}
		line = strings.Replace(line, "\n", "", -1) //将换行符替换为空，你可以理解是删除了换行符。
		if line == string(buf) {
			fmt.Println("对不起，您输入的用户或者ID已经存在了，请重新输入！")
			flag = 1
			break

		}
	}
	if flag == 0 {
		s := Student{
			ID:   id,
			Name: name,
		}
		buf, err := json.Marshal(s) //序列化一个结构体，
		if err != nil {
			log.Fatal("序列化报错是：%s", err)
		}
		fmt.Println(string(buf))
		f, err := os.OpenFile("student_info.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		f.WriteString(string(buf))
		f.WriteString("\n")
		f.Close()
		fmt.Println("写入成功")

	}
}
func save() {
	fmt.Sscan(line, &cmd, &file_name)
	f, err := os.Open("student_info.json") //打开一个文件，如果这个文件不存在的话就会报错。
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(f)                                                    //取出文件的内容
	f2, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644) //打开一个新文件
	if err != nil {
		log.Fatal(err)
	}
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		f2.WriteString(line) //将student_info.json内容写到指定的文件中去。
	}
	f2.Close()
	f.Close()
}

func load() {
	fmt.Sscan(line, &cmd, &file_name)
	f, err := os.Open(file_name) //打开一个文件，如果这个文件不存在的话就会报错。
	if err != nil {
		fmt.Println("对不起!系统没用找到该文件！")
		return
	}
	r := bufio.NewReader(f) //取出文件的内容
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Print(line)
	}
	f.Close()
}
