package mysql

import (
	"database/sql"
	"fmt"
	"github.com/KenianShi/filestore-server/config"
	"log"
	"os"
)

var db *sql.DB

func init() {
	db, _ := sql.Open("mysql", config.MySQLSource)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
