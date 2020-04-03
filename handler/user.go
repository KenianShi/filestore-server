package handler

import (
	"net/http"
	"io/ioutil"
	"github.com/KenianShi/filestore-server/util"
	dblayer "github.com/KenianShi/filestore-server/db"
	"fmt"
)

const pwd_salt  = "*#890"

//处理用户注册请求
func SignupHandler(w http.ResponseWriter,r *http.Request){
	if r.Method == http.MethodGet {
		data,err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	fmt.Printf("received username, passwd:%s, %s \n",username,passwd)
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("invalid parameters"))
		return
	}
	encPasswd := util.Sha1([]byte(passwd+pwd_salt))
	suc := dblayer.UserSignUp(username,encPasswd)
	if suc {
		fmt.Println("chuangjian chenggong")
		w.Write([]byte("SUCCESS"))
	}else{
		fmt.Println("chuangjian shibai")
		w.Write([]byte("Sign up Failed!"))
	}





}