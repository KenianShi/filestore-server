package main

import (
	"net/http"
	"github.com/KenianShi/filestore-server/handler"
	"fmt"
)

func main() {
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/success",handler.UploadSucHandler)

	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		fmt.Printf("Failed to start server, err: %s \n",err.Error())
	}
	fmt.Println("listening :8080 ....")
}
