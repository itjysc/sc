package main

import (
	"net/http"
	"log"
	"io"
	"html/template"
	"path/filepath"
	"fmt"
	//"time"
	"github.com/gorilla/sessions"
	"database/sql"
	_  "github.com/go-sql-driver/mysql"
	"crypto/md5"
	"github.com/jmoiron/sqlx"
)

var   (
	store = sessions.NewFilesystemStore("sessions")
)

func Reader(w http.ResponseWriter,name string,data interface{})  {
	path := filepath.Join("template",name+".tpl")
	tpl,err := template.ParseFiles(path)
	if err != nil {
		http.Error(w,err.Error(),500)
		return
	}
	err = tpl.Execute(w,data)
	if err != nil {
		http.Error(w,err.Error(),500)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	Reader(w,"login",nil)
}

func CheckLogin(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	user := r.FormValue("user")
	passwd := r.FormValue("password")



	if user == "admin" && passwd == "admin"{
		fmt.Fprintf(w,"登陆成功！")
		cookie := &http.Cookie{
			Name:"user",
			Value:user,
			MaxAge:10,
		}
		http.SetCookie(w,cookie)
		http.Redirect(w,r,"/yinzhengjie",302) //当登陆成功后就直接跳转到另外当一个函数去执行。
	}else {
		fmt.Fprintf(w,"user:%s,password:%s",user,passwd)
	}
}

var   (
	db m
)

func Add(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	name := r.FormValue("name")
	passwd := fmt.Sprintf("%x",md5.Sum([]byte(r.FormValue("password"))))
	note := r.FormValue("note")
	stmt,err := db.Prepare("INSERT INTO ")


	res,err := db.Exec("INSERT INTO user VALUES(NULL,?,?,?)",name,passwd,note,1)
	if err != nil {
		http.Error(w,err.Error(),500)
		return
	}
	log.Print(res.LastInsertId())
	log.Print(res.RowsAffected())
}

func Hello(w http.ResponseWriter,r *http.Request)  {
	_,err := r.Cookie("user")
	if err != nil {
		http.Redirect(w,r,"/yinzhengjie",302)
		return
	}
	io.WriteString(w,"尹正杰") //表示可以在页面返回字符串"尹正杰"
}

type User struct {
	Id int
	Name string
	Password string
	Note string
	Isadmin bool
}

func main() {
	db,err := sql.Open("mysql","golang:golang@tcp(59.110.12.72:3306)/go")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT CURRENT_USER()")
	var user  string
	row.Scan(&user)
	log.Print(user)

	rows,err := db.Query("SELECT * FROM user ") //一定要注意大小写，数据库是大写这里就得写大写，如过是小写就写小写。
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id int
		name string
		passwd string
		note string
		isadmin int
	)

	for rows.Next() {
		rows.Scan(&id,&name,&passwd,&note,&isadmin)
		log.Print(id,name,passwd,note,isadmin)
	}
	
	{
		dbx,err := sqlx.Open("mysql","golang:golang@tcp(59.110.12.72:3306)/go")
		if err !=  nil{
			log.Fatal(err)
		}
		var users []User
		err = dbx.Select(&users,"SELECT * FROM user")
		if err != nil {
			log.Fatal(err)
		}
		log.Print(users)
		err = dbx.Get(&user,"SELECT * FROM user WHERE name = ?","admin")
		if err != nil {
			log.Fatal(err)
		}
		log.Print(user)
	}
	
	http.HandleFunc("/login",Login)
	http.HandleFunc("/checkLogin",CheckLogin)
	http.HandleFunc("/yinzhengjie",Hello) //绑定一个函数，当用户在监听的端口后家关键字"／yinzhengjie"就能访问到"Helllo"这个函数啦。
	log.Fatal(http.ListenAndServe(":8090",nil)) //启动一个http服务
}
