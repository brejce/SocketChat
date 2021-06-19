package main

import (
	"encoding/json"
	"time"
)

//使用IdTime进行message查询
func GetMessage(idtime string) (Message, bool) { //获取单独的一条
	rdb := NewRedis(1)                        //获取redis客户端
	val, err := rdb.Get(ctx, idtime).Result() //使用IdTime获取Message
	rdb.Close()                               //关闭redis客户端
	if nil == err {
		var msg Message
		err := json.Unmarshal([]byte(val), &msg) //反序列化
		CheckError(err)
		return msg, true
		//在无err情况下返回Message，并设置状态为true
		//true表示获取成功
	} else {
		return Message{"", "", ""}, false
		//在err情况下返回空的Message，并设置状态为false
		//false表示获取失败
	}
}

//保存该message到数据库1
func SetMessage(m Message) bool {
	data, err := json.Marshal(m) //使用json.Marshal序列化，data数据类型为[]byte
	if nil == err {
		rdb := NewRedis(1)
		rdb.Set(ctx, m.IdTime, data, 24*time.Hour).Err() //保存Message,保存24小时24*time.Hour
		rdb.Close()
		return true //返回true表示存储成功
	} else {
		return false //返回false表示存储失败
	}
}

//获取全部的消息,将其打包为map
func GetAllMessageRange() []Message {
	var MessageSlice []Message
	b := GetAllKeys(1)
	for _, idtime := range b {
		m, _ := GetMessage(idtime)
		MessageSlice = append(MessageSlice, m)
	}
	return BubbleSortPro(MessageSlice)
}

//对Message进行冒泡排序,使其按照IdTime的先后顺序
func BubbleSortPro(arr []Message) []Message {
	length := len(arr)
	for i := 0; i < length; i++ {
		over := false
		for j := 0; j < length-i-1; j++ {
			if arr[j].IdTime > arr[j+1].IdTime {
				over = true
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
		if !over {
			break
		}
	}
	return arr
}
