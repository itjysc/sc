package main


/*
可参考：
		1>.https://godoc.org/html/template
		2>.https://godoc.org/net/http
		3>.https://github.com/go-sql-driver/mysql
		4>.https://godoc.org/database/sql
		5>.https://github.com/gorilla/handlers  //日志接口

*/

import (
	"log"
	"net/http"
	"io"
	"html/template"
	"path/filepath"
	"fmt"
	_ "github.com/go-sql-driver/mysql"  //要注意这个"_"，它的意思是倒入这个包的驱动，但是不引 用其它，因为以下代码是不会用的mysql。我们可以理解是只倒入mysql驱动。
	"crypto/md5"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB  //建立一个全局都db对象。
)



type User struct {
	Id int
	Name string
	Password string
	Note string
	Isadmin bool
}


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
	name := r.FormValue("user")  //读取名称为"user"的并赋值给user。
	passwd := r.FormValue("password")

	var user User
	err := db.Get(&user,"SELECT password FROM user WHERE name = ?",name)

	if err != nil {
		//http.Redirect(w,r,"/login",302)
		Reader(w,"login","user not found!")
		return
	}
	//fmt.Fprintf(w,"%x:%s",md5.Sum([]byte(passwd)),user.Password) //将用户密码原样子返回出去。
	if fmt.Sprintf("%x",md5.Sum([]byte(passwd))) == user.Password {  //登陆成功就跳转到这个界面上来。
		http.SetCookie(w,&http.Cookie{
			Name:"user",
			Value:name,
			MaxAge:10, //定义超时时间，意思是当超过10s后自动会弹出"登陆时间过期！"的提示。
		})
		http.Redirect(w,r,"/list",302)
	}else { //登陆失败就跳转到这个界面上来。
		Reader(w,"login","bad password!")
	}
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

func List(w http.ResponseWriter,r *http.Request)  {


	var users []User
	err := db.Select(&users,"SELECT * FROM user") //将查询到的内容赋值给users。可以获取到多条记录
	if err != nil {
		http.Error(w,err.Error(),500)
		return
	}
	Reader(w,"list.html",users)
}

func NeedLogin(h http.HandlerFunc)http.HandlerFunc  { //定义一个装饰器，我们也可以叫做定义一个中间件。
	return func(w http.ResponseWriter, r *http.Request) {
		_,err := r.Cookie("user") //验证用户是否登陆
		if err != nil { //如果没有登陆返回登陆界面
			Reader(w,"login","登陆时间过期！")
			return
		}
		h(w,r) //如果用户登陆就透传模式。不做任何到操作。
	}
}

type counter struct { //定义一个定时器
	count int
}

func (c *counter) ServeHTTP(w http.ResponseWriter,r *http.Request) { //注意方法名称必须为"ServeHTTP"，不然会编译报错哟！
	c.count++
	fmt.Fprintf(w,"%d",c.count)
}

func main() {
	var err   error
	db,err = sqlx.Open("mysql","golang:golang@tcp(59.110.12.72:3306)/go") //链接数据库，这个mysql就是上面引入的驱动。
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

	//http.HandleFunc("/login",Login)
	//http.HandleFunc("/checkLogin",CheckLogin)
	//http.HandleFunc("/hello",NeedLogin(Hello)) //给hello函数绑定一个装饰器（我们也可以叫做中间件）
	//http.HandleFunc("/add",NeedLogin(Add))
	//http.HandleFunc("/list",NeedLogin(List))
	c := new(counter)
	//http.Handle("/counter",c) //调用计数器。
	//h := handlers.LoggingHandler(os.Stderr,http.DefaultServeMux) //定义一个有日志输出的对象。它其实也是一个中间件（装饰器）。

	mux := http.NewServeMux()  //手动构造一个handle
	mux.HandleFunc("login",Login)
	mux.Handle("/counter",c)
	log.Fatal(http.ListenAndServe(":8088",mux)) //启动一个http服务使其监听本地所有的8080端口。并将访问记录传给h，你可以理解所有的流量都要经过h。
}
