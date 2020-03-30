package handler

import (
	"net/http"
	"io/ioutil"
	"io"
	"os"
	"fmt"
)

func UploadHandler(w http.ResponseWriter,r *http.Request){
	if r.Method == "GET"{
		//返回上传的HTML页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w,"Internal server error, Failed to read the html!")
			return
		}else{
			io.WriteString(w,string(data))
		}
	}else if r.Method == "POST"{
		//接收文件流及存储到本地
		file, head, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			fmt.Printf("Failed to get data, err: %s \n",err.Error())
			return
		}

		newFile,err := os.Create("./store/"+head.Filename)
		defer newFile.Close()
		if err != nil {
			fmt.Printf("Failed to create newFile, err: %s \n",err.Error())
			return

		}

		_,err = io.Copy(newFile,file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s \n",err.Error())
			return
		}

		http.Redirect(w,r,"/file/upload/success",http.StatusFound)

	}


}

func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"Upload finished successfully.")
}
