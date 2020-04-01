package db

import (
	mydb "github.com/KenianShi/filestore-server/db/mysql"
	"fmt"
)

func OnFileUploadFinished(filehash,filename,fileaddr string,filesize int64)bool{
	stmt,err := mydb.DBConn().Prepare("insert ignore into tbl_file(`file_sha1`,`file_name`," +
		"`file_size`,`file_addr`,status) values(?,?,?,?,1)")
	defer stmt.Close()
	if err != nil {
		fmt.Println("Failed to prepare statment, err: "+err.Error())
		return false
	}

	ret,err := stmt.Exec(filehash,filename,filesize,fileaddr)
	if err != nil{
		fmt.Println("Exec statment err: "+err.Error())
		return false
	}
	if rf,err := ret.RowsAffected();nil == err{
		if rf <= 0{
			fmt.Printf("File with the hash:%s has been existed before \n",filehash)
		}
		return true
	}
	return false
}