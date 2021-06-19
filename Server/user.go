package main

import (
	"encoding/json"
	"fmt"
)

//将User转化为[]byte
func StructTojson(user User) []byte {
	//将User结构体序列化操作封装成方法
	data, err := json.Marshal(user)
	//使用json.Marshal序列化User，data数据类型为[]byte
	CheckError(err)
	//error处理
	return data
	//返回序列化之后的User
}

//将user型的[]byte转化为User
func JsonTostruct(b []byte) User {
	var user User
	err := json.Unmarshal(b, &user)
	CheckError(err)
	return user
}

//保存改用户到数据库
func SetUser(user User) bool { //保存User
	//使用getUser进行查询
	_, b := GetUser(user)
	if b { //如果有此人就直接返回此人
		return false
	} else {
		rdb := NewRedis(0)                                          //连接数据库0
		err := rdb.Set(ctx, user.Name, StructTojson(user), 0).Err() //保存用户
		CheckError(err)
		rdb.Close()
		return true
	}
}

//修改密码
//为了不新增对象/结构体,这里约定客户端发来的用户信息如下:
// var user=User{
// 	Name: "用户名",
// 	Passwd: "旧密码",
// 	Status: "新密码",
// }
func SetPasswd(user User) bool {
	uu, b := GetUser(user)
	if b {
		if user.Passwd == uu.Passwd {
			//密码正确
			u := User{
				Name:   user.Name,
				Passwd: user.Status,
				Status: "",
			}
			DeletSomething(u.Name, 0)
			if SetUser(u) {
				//修改成功
				fmt.Println("此u给i成功")
				return true
			} else {
				//修改失败
				return false
			}
		} else {
			//密码错误
			return false
		}
	} else {
		//没有这个人,用户名错误
		return false
	}
}

//true-用户名密码正确/保存成功 false-密码错误 nil-用户名错误
func GetUser(user User) (User, bool) { //获取用户全部的信息，这里使用user.Name作为唯一识别号
	rdb := NewRedis(0)                           //连接数据库
	val, err := rdb.Get(ctx, user.Name).Result() //使用user.Name来进行查找
	rdb.Close()
	if nil == err {
		return JsonTostruct([]byte(val)), true //如果有此用户，则返回User，true
	}
	return User{"nil", "nil", "nil"}, false //如果没有此用户，则返回nil，false
}
