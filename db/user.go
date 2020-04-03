package db

import (
	mydb "github.com/KenianShi/filestore-server/db/mysql"
	"fmt"
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
