package main

import (
	"net/http"
	"os"
	"fmt"
	"bufio"
	"strconv"
	"bytes"
	"io/ioutil"
	"io"
	"net/url"
	"github.com/json-iterator/go"
)

func multipartUpload(filename string,targetURL string,chunkSize int) error{
	f,err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0
	ch := make(chan int)
	buf := make([]byte,chunkSize)
	for {
		n,err := bfRd.Read(buf)
		if n < 0 {
			break
		}
		index++
		bufCopied := make([]byte,5*1024*1024)
		copy(bufCopied,buf)
		go func(b []byte,curIdx int) {
			fmt.Printf("upload_size:%d \n",len(b))
			fmt.Println("index++++++++++++++++:"+strconv.Itoa(curIdx))
			resp,err := http.Post(targetURL+"&index="+strconv.Itoa(curIdx),
				"multipart/form-data",bytes.NewReader(b))
			if err != nil {
				fmt.Println(err)
			}
			body,err := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v %+v \n",string(body),err)
			resp.Body.Close()
			ch <- curIdx
		}(bufCopied[:n],index)

		if err != nil {
			if err == io.EOF{
				break
			}else{
				fmt.Println(err.Error())
			}
		}
	}
	for idx := 0;idx < index;idx++{
		select{
		case res := <-ch:
		fmt.Println(res)
		}
	}
	return nil
}



func main() {
	username := "admin"
	token := "04133679f7c1b3f425d92e47ffd2acae5e95c2c9"
	filehash := "abcdefg"
	resp,err := http.PostForm(
		"http://127.0.0.1:8080/file/mpupload/init",url.Values{
			"username":{username},
			"token":{token},
			"filehash":{filehash},
			"filesize":{"1235"},
		})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	uploadID := jsoniter.Get(body,"data").Get("UploadID").ToString()
	chunkSize := jsoniter.Get(body,"data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid: %s chunksize: %d \n",uploadID,chunkSize)
	filename := `D:\workspace.zip`
	tURL := "http://127.0.0.1:8080/file/mpupload/uppart?username=admin&token="+token+"&uploadid="+
		token+"&uploadid="+uploadID

		fmt.Printf("filename:%s,\n uploadID:%s \ntURL:%s\n",filename,uploadID,tURL)


	err = multipartUpload(filename,tURL,chunkSize)
	if err != nil{
		panic(err)
	}

	resp,err = http.PostForm(
		"http://127.0.0.1:8080/file/upload/complete",
		url.Values{
			"username":{username},
			"token":{token},
			"filehash":{filehash},
			"filename":{"noname"},
			"uploadid":{uploadID},
		})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer resp.Body.Close()
	body,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Println("complete result: %s \n",string(body))

}


