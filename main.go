package main

import (
	"fmt"
	"github.com/KenianShi/filestore-server/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err: %s \n", err.Error())
	}
	fmt.Println("listening :8080 ....")
}


