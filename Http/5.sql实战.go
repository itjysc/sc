package main


/*
可参考：
		1>.https://godoc.org/html/template
		2>.https://godoc.org/net/http
		3>.https://github.com/go-sql-driver/mysql
		4>.https://godoc.org/database/sql

*/

import (
	"log"
	"net/http"
	"io"
	"html/template"
	"path/filepath"
	"fmt"
	_ "github.com/go-sql-driver/mysql"  //要注意这个"_"，它的意思是倒入这个包的驱动，但是不引 用其它，因为以下代码是不会用的mysql。我们可以理解是只倒入mysql驱动。
	"database/sql"
	"crypto/md5"
)

var (
	db *sql.DB  //建立一个全局都db对象。
)

func Reader(w http.ResponseWriter,name string,data interface{})  { //用于对页面的渲染。
	path := filepath.Join("/yinzhengjie/golang/path/src/yinzhengjie/Http/template",name+ ".tpl") //定义文件的路径
	tpl,err := template.ParseFiles(path)  //建立一个名为tpl的模板对象，ParseFiles方法可以从一个文件中解析出来一个模板对象。因此我们需要传入一个文件的路径。
	if err != nil {
		http.Error(w,err.Error(),500) //注意要传3个参数进去，最后一个数字表示返回给用户的错误提示500
		return
	}
	err = tpl.Execute(w,data) //第二个参数可以传一个动态数据，比如你更新的新闻稿之类的，如果测试的话也可以传空值。
	if err != nil {
		http.Error(w,err.Error(),500) //注意要传3个参数进去，最后一个数字表示返回给用户的错误提示500
		return
	}
}

func Login(w http.ResponseWriter,r *http.Request)  {
	Reader(w,"login",nil)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) { //定义读取表单和数据的函数
	r.ParseForm()  //先解析一下表单数据。
	user := r.FormValue("user")  //读取名称为"user"的并赋值给user。
	passwd := r.FormValue("password")

	row := db.QueryRow("SELECT password FROM user WHERE name = ?",user)

	var dbpass  string
	err := row.Scan(&dbpass)
	if err != nil {
		http.Redirect(w,r,"/login",302)
		return
	}
	fmt.Fprintf(w,"%x:%s",md5.Sum([]byte(passwd)),dbpass) //将用户密码原样子返回出去。

}

func Add(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	name := r.FormValue("name")
	passwd := fmt.Sprintf("%x",md5.Sum([]byte(r.FormValue("password"))))
	note := r.FormValue("note")
/*
	stmt,err := db.Prepare("INSERT INTO user VALUES (NULL ,?,?,?,?)")  //利用Prepare可以向数据库插入多条数据。
	stmt.Exec(name,passwd,note,1)
	stmt.Exec(name,passwd,note,1)
*/

/*
	tx,err := db.Begin() //创建一个事物tx。其实tx事物我们可以当作一个db来使用哟，事物就是可以支持回滚，比如银行转账的案例。
	tx.Exec(name,passwd,note,1) //插入数据
	tx.Commit() //提交事物
	tx.Rollback() //回滚事物。
*/

	res,err := db.Exec("INSERT INTO user VALUES (NULL ,?,?,?,?)",name,passwd,note,1) //用Exec方法向数据库插入一条数据，
	if err != nil {
		http.Error(w,err.Error(),500) //500表示内部错误。
		return
	}
	log.Print(res.LastInsertId()) //用户的字段ID
	log.Print(res.RowsAffected()) //表示影响到对少行。
}


func Hello(w http.ResponseWriter, r *http.Request) {     //定义一个名称为Hello的处理函数。
	io.WriteString(w,"尹正杰来也！")
}

/*
func init()  { //用于定义初始化操作，在main函数调用之前。

}
*/


func main() {
	var err   error
	db,err = sql.Open("mysql","golang:golang@tcp(59.110.12.72:3306)/go") //链接数据库，这个mysql就是上面引入的驱动。
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

	rows,err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var  ( //以下都变量都是数据库都key。
		id int
		name string
		passwd string
		note string
		isadmin int
	)
	for rows.Next() {
		rows.Scan(&id,&name,&passwd,&note,&isadmin) //将所有对数据都获取出来。
		log.Print(id,name,passwd,note,isadmin)
	}

	http.HandleFunc("/login",Login)
	http.HandleFunc("/checkLogin",CheckLogin)
	http.HandleFunc("/hello",Hello) //给这个http服务绑定一个处理函数，用关键字"hello"来绑定Hello方法。
	http.HandleFunc("/add",Add)
	log.Fatal(http.ListenAndServe(":18080",nil)) //启动一个http服务使其监听本地所有的8080端口。
}
