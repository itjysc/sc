package main


/*
可参考：
		1>.https://godoc.org/html/template
		2>.https://godoc.org/net/http
		3>.https://github.com/gorilla/sessions
*/

import (
	"log"
	"net/http"
	"io"
	"html/template"
	"path/filepath"
	"fmt"
)

func Hello(w http.ResponseWriter, r *http.Request) {     //定义一个名称为Hello的处理函数。
	_,err := r.Cookie("user") //获取Cookie方法，如果获取到用户的"user"信息，就是不需要重新登陆，如果获取不到，就需要重新登陆。
	if err != nil {
		http.Redirect(w,r,"/login",302)
		return  //跳转之后需要返回，
	}
	io.WriteString(w,"尹正杰来也！")
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
	user := r.FormValue("user")  //读取名称为"user"的并赋值给user。
	passwd := r.FormValue("password")
	if user == "admin" && passwd == "admin" {
		//fmt.Fprintf(w,"login seccussful!")
		cookie := &http.Cookie{
			Name:"user",
			Value:user,
			MaxAge:10,  //定义超时时间
		}
		http.SetCookie(w,cookie) //配置cookie
		http.Redirect(w,r,"/hello",302)  //如果用admin用户登陆成功就跳转到我们想要的页面上去。Redirect方法可以实现跳转。
	}else {
		fmt.Fprintf(w,"用户名：%s,密码：%s",user,passwd) //将用户密码原样子返回出去。
	}
}

func main() {
	http.HandleFunc("/login",Login)
	http.HandleFunc("/checkLogin",CheckLogin)
	http.HandleFunc("/hello",Hello) //给这个http服务绑定一个处理函数，用关键字"hello"来绑定Hello方法。
	log.Fatal(http.ListenAndServe(":8080",nil)) //启动一个http服务使其监听本地所有的8080端口。
}
