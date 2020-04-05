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
	"github.com/KenianShi/filestore-server/db"
	"strconv"
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
		//meta.UpdateFileMeta(fileMeta)
		//UpdateFileMeteDB
		if !meta.UpdateFileMetaDB(fileMeta) {
			w.Write([]byte("UpdateFileMetaDB error"))
			return
		}

		//TODO 更新用户文件表记录
		r.ParseForm()
		username := r.Form.Get("username")
		suc := db.OnUserFileUploadFinished(username,fileMeta.FileSha1,fileMeta.FileName,fileMeta.FileSize)
		if suc{
			http.Redirect(w,r,"/static/view/home.html",http.StatusFound)
		}else{
			w.Write([]byte("Update UserFile db error"))
		}
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished successfully.")
}

//获取文件原信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("filehash")
	fmeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.Write([]byte("查询失败"))
		return
	}
	data,err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func FileQueryHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	limitCnt,_ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMetas,err := meta.GetLatestFileMetasDB(limitCnt)
	fileMetas,err := db.QueryUserFileMetas(username,limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(fileMetas)
	data,err := json.Marshal(fileMetas)
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
	if !meta.UpdateFileMetaDB(fmeta){
		w.Write([]byte("update file meta to DB error"))
		return
	}


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

func TryFastUploadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	// 1. 解析请求参数
	username := r.Form.Get("username")
	fileHash := r.Form.Get("filehash")
	fileName := r.Form.Get("filename")
	filesize,_ := strconv.Atoi(r.Form.Get("filesize"))

	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if fileMeta == nil {
		resp := util.RespMsg{
			Code:-1,
			Msg:"秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	suc := db.OnUserFileUploadFinished(username,fileHash,fileName,int64(filesize))
	resp := util.RespMsg{}

	if suc{
		resp.Code = 0
		resp.Msg = "秒传成功"
	}else{
		//resp := util.RespMsg{
		//	Code:-2,
		//	Msg:"秒传失败，请稍后再试",
		//}
		resp.Code = -2
		resp.Msg = "秒传失败，请稍后再试"
	}

	w.Write(resp.JSONBytes())
	return
	// 2. 从文件表中查询相同Hash的文件
	// 3. 查找不到则返回秒传失败
	// 4. 找到的话，则将文件信息写入用户文件信息表
}

