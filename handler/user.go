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

func SignInHandler(w http.ResponseWriter,r *http.Request){
	// 检查用户名和密码
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(passwd + pwd_salt))
	pwdChecked := dblayer.UserSignIn(username,encPasswd)
	if !pwdChecked{
		w.Write([]byte("FAILED"))
		fmt.Println("Failed to check password")
		return
	}

	// 2 生成访问凭证
	token := dblayer.GenToken(username)
	UpdataToken := dblayer.UpdateToken(username,token)
	if !UpdataToken{
		w.Write([]byte("FAILED"))
		fmt.Println("Update Token Failed")
		return
	}
	fmt.Println("r.Host: ",r.Host)
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code:0,
		Msg:"OK",
		Data: struct {
			Location 	string
			Username 	string
			Token 		string
		}{
			Location:"http://" + r.Host + "/static/view/home.html",
			Username:username,
			Token:token,
		},
	}
	w.Write(resp.JSONBytes())
}

func UserInfoHandler(w http.ResponseWriter,r *http.Request){
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")

	//2.验证token是否有效
	isValidToken := IsTokenValid(token)
	if !isValidToken{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//3.查询用户信息
	user,err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fmt.Printf("%+v \n",user)
	//4.组装并且响应用户数据
	resp := util.RespMsg{
		Code:0,
		Msg:"success",
		Data:user,
	}
	w.Write(resp.JSONBytes())
}

func IsTokenValid(token string) bool{
	// 判断Token的时间有效性

	// 从数据库表tbl_user_token中查询username对应的token


	return true
}