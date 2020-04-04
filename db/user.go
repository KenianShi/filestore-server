package db

import (
	mydb "github.com/KenianShi/filestore-server/db/mysql"
	"fmt"
	"time"
	"github.com/KenianShi/filestore-server/util"
)

//UserSignUp:通过用户名及密码完成User表的注册操作
func UserSignUp(username string,passwd string)bool{
	stmt,err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`) values (?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Failed to insert ,err: ",err.Error())
		return false
	}
	ret,err := stmt.Exec(username,passwd)
	if err != nil {
		fmt.Println("stmt exec error: "+err.Error())
		return false
	}
	fmt.Println("ret执行成功")
	if rowsAffected,err := ret.RowsAffected();nil != err  {
		fmt.Println("fafadsklfjasdklf")
		return false
	}else if rowsAffected <= 0 {
		fmt.Println("the username has been existed!")
		return false
	}else{
		fmt.Println("执行成功")
		return true
	}
}

func UserSignIn(username,encpasswd string)bool{
	stmt,err := mydb.DBConn().Prepare("select * from tbl_user where user_name= ? limit 1")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Failed to prepare the query sql")
		return false
	}
	rows,err := stmt.Query(username)
	if err != nil {
		fmt.Println("stmt query error")
		return false
	}else if rows == nil {
		fmt.Println("username not found")
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpasswd{
		return true
	}
	return false
}

func GenToken(username string)string{
	// token 是一位长度为40的字符，来源：md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x",time.Now().Unix())    		//将int64转成string
	tokenPrefix := util.MD5([]byte(ts))
	return tokenPrefix + ts[:8]
}

func UpdateToken(username,token string)bool{
	stmt,err := mydb.DBConn().Prepare("replace into tbl_user_token(`user_name`,`user_token`) values(?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("update token prepare failed")
		return false
	}
	_,err = stmt.Exec(username,token)
	if err != nil {
		fmt.Println("update token exec failed")
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}


func GetUserInfo(username string)(User,error){
	user := User{}
	stmt,err := mydb.DBConn().Prepare(`select user_name,signup_at from tbl_user where user_name=? limit 1`)
	defer stmt.Close()
	if err != nil {
		fmt.Println("stmt prepare error: ",err.Error())
		return user,err
	}
	err = stmt.QueryRow(username).Scan(&user.Username,&user.SignupAt)
	if err != nil {
		fmt.Println("stamt query row err: ",err.Error())
		return user,err
	}
	return user,nil
}


