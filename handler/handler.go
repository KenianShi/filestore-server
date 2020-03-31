package handler

import (
	"encoding/json"
	"fmt"
	"github.com/KenianShi/filestore-server/meta"
	"github.com/KenianShi/filestore-server/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回上传的HTML页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internal server error, Failed to read the html!")
			return
		} else {
			io.WriteString(w, string(data))
		}
	} else if r.Method == "POST" {
		//接收文件流及存储到本地
		file, head, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			fmt.Printf("Failed to get data, err: %s \n", err.Error())
			return
		}

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./store/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		defer newFile.Close()
		if err != nil {
			fmt.Printf("Failed to create newFile, err: %s \n", err.Error())
			return
		}
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s \n", err.Error())
			return
		}
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished successfully.")
}

//获取文件原信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form["filehash"][0]
	fmeta := meta.GetFileMeta(fileHash)
	data, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	//fsha1 := r.Form["filehash"][0]
	fhash := r.Form.Get("filehash")
	fmeta := meta.GetFileMeta(fhash)
	f,err := os.Open(fmeta.Location)
	defer f.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data,err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("content-disposition","attachment; filename=\""+fmeta.FileName+"\"")
	w.Write(data)
}

//FileMetaUpdateHandler
func FileMetaUpdateHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	if r.Method != "POST"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fname := r.Form.Get("filename")
	fsha1 := r.Form.Get("filehash")
	op := r.Form.Get("op")
	if op != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fmeta := meta.GetFileMeta(fsha1)
	fmeta.FileName = fname
	meta.UpdateFileMeta(fmeta)
	data,err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//FileDeleteHandler 删除文件以及原信息
func FileDeleteHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fmeta := meta.GetFileMeta(fsha1)
	path := fmeta.Location
	name := fmeta.FileName
	err := os.Remove(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meta.RemoveFileMeta(fsha1)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete "+ name + " success"))
}

