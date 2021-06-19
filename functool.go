package main

import (
	"context"
	"log"
	"unsafe"

	"github.com/go-redis/redis/v8"
)

//Message -> Db 1
type Message struct {
	Name   string `json:"name"`
	IdTime string `json:"idtime"`
	Msge   string `json:"msge"`
}

//User ->Db 0
type User struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
	Status string `json:"status"`
}

var UserMap = make(map[int]User)

// var MessageSlice []Message

var ctx = context.Background()

//将数据库打包
func NewRedis(db int) *redis.Client { //将数据库连接操作打包为方法使用newRdis(0)方法带入数据库名调用即可
	rdb := redis.NewClient(&redis.Options{
		Addr:     "47.109.26.249:6379", //数据库默认安装在开发机，监听localhost，默认端口为6379
		Password: "169809",             // no password set
		DB:       db,                   // use default DB
	})
	return rdb //返回数据库客户端
}

//获取该数据库里所有的key
func GetAllKeys(db int) []string {
	rdb := NewRedis(db)
	defer rdb.Close()
	keys, err := rdb.Keys(ctx, "*").Result()
	CheckError(err)
	return keys
}

//error处理,可以使用客户端log处理的逻辑将错误信息收集保存到数据库,这里不在展开
func CheckError(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

//将string转为[]byte
func StrTobyte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//将[]byte转为string
func ByteTostr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//该方法可以删除Redis里任意数据库 Db 的任意 Key
func DeletSomething(key string, Db int) {
	rdb := NewRedis(Db)
	rdb.Del(ctx, key).Err()
	rdb.Close()
}

//判断是否存在该户,存在true,不存在false
func IsOurUser(name string) bool {
	bo := false
	for _, u1 := range UserMap {
		if u1.Name == name {
			bo = true
		}
	}
	return bo
}

//将用户踢出用户Map,表示该用户没有登录
func DeletLoginUser(u User) {
	for i, user := range UserMap {
		if user.Name == u.Name {
			delete(UserMap, i)
		}
	}
}
