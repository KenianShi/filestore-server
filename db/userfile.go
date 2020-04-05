package db

import (
	mydb "github.com/KenianShi/filestore-server/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	UserName 	string
	FileHash 	string
	FileName 	string
	FileSize 	int64
	UploadAt 	string
	LastUpdated	string
}

func OnUserFileUploadFinished(username,filehash,filename string,filesize int64)bool{
	stmt,err := mydb.DBConn().Prepare("insert ignore into tbl_user_file(`user_name`,`file_sha1`,`file_name`," +
		"`file_size`,`upload_at`) values(?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("stmt prepare err: ",err.Error())
		return false
	}
	_,err = stmt.Exec(username,filehash,filename,filesize,time.Now())
	if err != nil {
		fmt.Println("stmt exec err: ",err.Error())
		return false
	}
	return true
}

func QueryUserFileMetas(username string,limit int)([]UserFile,error){
	stmt,err := mydb.DBConn().Prepare("select file_sha1,file_name,file_size,upload_at,last_update " +
		" from tbl_user_file where user_name =? limit ?")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Query Userfile prepare err: ",err.Error())
		return nil,err
	}
	rows,err := stmt.Query(username,limit)
	if err != nil {
		return nil,err
	}
	var userFiles []UserFile
	for rows.Next(){
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash,&ufile.FileName,&ufile.FileSize,&ufile.UploadAt,&ufile.LastUpdated)
		if err != nil {
			fmt.Println("rows scan",err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles,nil
}


