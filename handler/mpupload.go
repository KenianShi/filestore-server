package handler

import (
	"net/http"
	"strconv"
	"github.com/KenianShi/filestore-server/util"
	"github.com/KenianShi/filestore-server/cache"
	"time"
	"math"
	"os"
	"path"
	"github.com/gomodule/redigo/redis"
	"strings"
	"github.com/KenianShi/filestore-server/db"
	"fmt"
)

type MultipartUploadInfo struct {
	FileHash 		string
	FileSize 		int
	UploadID		string
	ChunkSize 		int
	ChunkCount 		int
}

//初始化分块信息
func InitialMultipartUploadHandler(w http.ResponseWriter,r *http.Request){
	//1. 解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize,err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1,"params invalid",nil).JSONBytes())
		return
	}
	fmt.Printf("username:%s ,filehash:%s,filesize:%s\n",username,filehash,filesize)
	//2. 获得redis的一个连接
	rConn := cache.RedisPool().Get()
	defer rConn.Close()
	//3. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:filehash,
		FileSize:filesize,
		UploadID:username+fmt.Sprintf("%x",time.Now().UnixNano()),
		ChunkSize:5*1024*1024,
		ChunkCount:int(math.Ceil(float64(filesize)/(5*1024*1024))),
	}

	//4. 将初始化信息写入redis
	rConn.Do("HSET","MP_"+upInfo.UploadID,"chunkcount",upInfo.ChunkCount)
	rConn.Do("HSET","MP_"+upInfo.UploadID,"filehash",upInfo.FileHash)
	rConn.Do("HSET","MP_"+upInfo.UploadID,"filesize",upInfo.FileSize)
	fmt.Println(upInfo)
	//5. 将响应初始化数据返回给客户端
	w.Write(util.NewRespMsg(0,"OK",upInfo).JSONBytes())
}

//上传分块信息
func UploadPartHandler(w http.ResponseWriter,r *http.Request){
	//1. 解析用户请求参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")
	fmt.Println("获取到的index：",chunkIndex)
	//2. 获取redis连接池
	rConn := cache.RedisPool().Get()
	defer rConn.Close()

	//3. 获得文件句柄，用于存储分块内容
	fpath := "./data/"+ uploadID + "/" +chunkIndex
	os.MkdirAll(path.Dir(fpath),0744)
	fd,err := os.Create(fpath)
	if err != nil {
		fmt.Println(err.Error())
		w.Write(util.NewRespMsg(-1,"Upload part failed",nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte,1024*1024)
	for{
		n,err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	//4. 更新redis缓存
	rConn.Do("HSET","MP_"+uploadID,"chkidx_"+chunkIndex,1)

	//5. 返回处理结果到客户端
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}

//通知上传合并
func CompleteUploadHandler(w http.ResponseWriter,r *http.Request){
	//1. 解析用户参数
	r.ParseForm()
	username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	filehash := r.Form.Get("filehash")
	filesize,_ := strconv.Atoi(r.Form.Get("filesize"))
	filename := r.Form.Get("filename")

	//2. 获得redis链接
	rConn := cache.RedisPool().Get()
	defer rConn.Close()

	//3. 通过uploadid查询redis，判断所有分块是否上传完成
	data,err := redis.Values(rConn.Do("HGETALL","MP_"+uploadID))
	if err != nil {
		w.Write(util.NewRespMsg(-1,"complete upload failed",nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkedCount := 0
	for i := 0;i < len(data); i += 2{
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount"{
			totalCount,_ = strconv.Atoi(v)
		}else if strings.HasPrefix(k,"chkidx_") && v == "1" {
			chunkedCount++
		}
	}
	if totalCount != chunkedCount {
		w.Write(util.NewRespMsg(-2,"invalid request",nil).JSONBytes())
		return
	}
	
	// 4. TODO 合并分块

	// 5. 更新唯一文件表以及用户文件表
	db.OnFileUploadFinished(filehash,filename,"",int64(filesize))
	db.OnUserFileUploadFinished(username,filehash,filename,int64(filesize))
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}


