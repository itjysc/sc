package main

import (
	"github.com/jmoiron/sqlx"
	"log"
	_ "github.com/go-sql-driver/mysql"  //要注意这个"_"，它的意思是倒入这个包的驱动，但是不引 用其它，因为以下代码是不会用的mysql。我们可以理解是只倒入mysql驱动。
)

type User struct {
	Id int
	Name string
	Password string
	Note string
	Isadmin bool
}


func main() {
	db,err := sqlx.Open("mysql","golang:golang@tcp(59.110.12.72:3306)/go") //链接数据库，这个mysql就是上面引入的驱动。
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	var users []User
	db.Select(&users,"SELECT * FROM user") //将查询到的内容赋值给users。可以获取到多条记录
	if err != nil {
		log.Fatal(err)
	}
	log.Print(users)

	var   user User
	err = db.Get(&user,"SELECT * FROM user WHERE name = ?","admin") //可以获取到一条记录
	if err != nil {
		log.Fatal(err)
	}
	log.Print(user)
}
