package db

import (
	"database/sql"
	"fmt"
	mydb "github.com/KenianShi/filestore-server/db/mysql"
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

type TableFile struct{
	FileHash 	string
	FileName 	sql.NullString
	FileSize 	sql.NullInt64
	FileAddr 	sql.NullString
}

func GetFileMeta(fileSha1 string)(*TableFile,error){
	stmt,err := mydb.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file where file_sha1=? and status =1 limit 1 ")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	tfile := TableFile{}
	err = stmt.QueryRow(fileSha1).Scan(&tfile.FileHash,&tfile.FileAddr,&tfile.FileName,&tfile.FileSize)      //scan 接收的是指针
	if err != nil {
		if err == sql.ErrNoRows{		//查不到相应的数据
			return nil,nil
		}else{
			fmt.Println(err.Error())
			return nil,err
		}
	}

	return &tfile,nil
}

//从数据库批量获得元数据信息
func GetFileMetaList(limit int)([]TableFile,error){
	stmt,err := mydb.DBConn().Prepare("select file_name,file_sha1,file_size,file_addr from tbl_file where  status = 1 limit ?")
	defer stmt.Close()
	if err != nil {
		fmt.Println("stmt prepare err: ",err.Error())
		return nil,err
	}
	rows,err := stmt.Query(limit)
	if err != nil {
		fmt.Println("stmt exec err: ",err.Error())
		return nil,err
	}
	//cloumns,_ := rows.Columns()
	//values := make([]sql.RawBytes,len(cloumns))
	var tableFiles []TableFile
	for rows.Next(){
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileName,&tfile.FileHash,&tfile.FileSize,&tfile.FileAddr)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tableFiles = append(tableFiles,tfile)
	}
	return tableFiles,nil
}
