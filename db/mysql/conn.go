package mysql

import (
	"database/sql"
	"fmt"
	"github.com/KenianShi/filestore-server/config"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
)
//注意此处需要import github.com/go-sql-driver/mysql

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", config.MySQLSource)        //重要，注意此处一定不能用:=，不然会新构建一个局部的db，生命周期仅在本方法内，不能给全局变量赋值
	if err != nil {
		fmt.Printf("OPEN DB error: %s \n",err.Error())
		os.Exit(1)
	}

	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to DB success")
}

func DBConn() *sql.DB {
	return db
}
//
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns,_ := rows.Columns()
	scanArgs := make([]interface{},len(columns))
	values := make([]interface{},len(columns))
	for j := range values{
		scanArgs[j] = &values[j]				//是可以对切片元素取地址的，但是不能对map元素取地址，因为hash桶的位置不明确
	}
	record := make(map[string]interface{})
	records := make([]map[string]interface{},0)
	if rows.Next(){
		err := rows.Scan(scanArgs...)
		checkErr(err)
		for i,col := range values{
			if col != nil {
				record[columns[i]] = col
			}
		}
	records = append(records,record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
