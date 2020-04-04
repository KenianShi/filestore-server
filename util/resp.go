package util

import (
	"encoding/json"
	"log"
	"fmt"
)

type RespMsg struct {
	Code 	int 		`json:"code"`
	Msg 	string 		`json:"msg"`
	Data 	interface{} `json:"data"`
}

func NewRespMsg(code int,msg string,data interface{}) *RespMsg{
	return &RespMsg{
		Code:code,
		Msg:msg,
		Data:data,
	}
}

//将对象转成[]byte二进制数组
func (resp *RespMsg) JSONBytes() []byte{
	r,err := json.Marshal(resp)
	if err != nil {
		log.Panic(err)
	}
	return r
}

//将对象转成string数组
func (resp *RespMsg) JSONString() string{
	r,err := json.Marshal(resp)
	if err != nil {
		log.Panic(err)
	}
	return string(r)
}

//只包含code和message的响应体([]byte)
func GenSimpleRespStream(code int,msg string)[]byte{
	return []byte(fmt.Sprintf(`{"code":"%d,"msg":"%s"}`,code,msg))
}

// 只包含code和message的响应体
func GenSimpleRespString(code int,msg string) string{
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`,code,msg)
}

