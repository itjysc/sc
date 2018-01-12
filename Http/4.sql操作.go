package main

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"  //要注意这个"_"，它的意思是倒入这个包的驱动，但是不引 用其它，因为以下代码是不会用的mysql。我们可以理解是只倒入mysql驱动。
)

func main() {
	db,err := sql.Open("mysql","golang:golang@tcp(59.110.12.72:3306)/go") //链接数据库，这个mysql就是上面引入的驱动。
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	row := db.QueryRow("SELECT CURRENT_USER() ") //对数据库进行查询。
	var user   string
	row.Scan(&user) //获取链接数据库对用户信息。
	log.Print(user)
}


