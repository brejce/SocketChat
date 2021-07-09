# Socket
## 使用方法
### 服务器端 Server 文件夹
安装 redis 参考官方网站  
修改 functool.go 内的数据库地址  
修改 main.go 内监听地址  
go run /server  
### 前端 SocketuView 文件夹
使用 hbuilderx 导入 SocketuView   
修改 app.vue内服务器地址  
运行到手机(h5端暂不支持)  
## 示例运行结果可以查看截图文件夹
## 第四章 数据库与项目的设计及搭建

### 第一节 数据库的设计与搭建

#### 项目文件结构
##### server/
- functool.go /提供公共方法
- go.mod /go mod 包管理器文件
- go.sum /go mod 包管理器文件
- client.go /为每一个客户端提供监听功能
- hub.go /每一个客户端都集中存储到hub里面
- message.go /关于 message 的方法
- user.go /关于 user 的方法
- main.go /服务器的入口

打开 VsCode 新建文件夹 sverver/ ，文件位置可以任选,本文选择在 /home/brejce/Program/server 。

在终端进入 Server/ 执行以下代码。
```Bash
go run mod init server
```
获得如下信息，表示我们已经成功使用 go mod 作为包管理器，这样我们不必手动去导包，go mod 会自动搞定。
```Bash
go: creating new go.mod: module Server
go: to add module requirements and sums:
        go mod tidy
```

打开 go.mod  可以看到。
```Golang
//go.mod
module server
go 1.16
require github.com/go-redis/redis/v8 v8.8.0

```
redis/v8 这是我们使用 Golang 操作 Redis 数据库关键的 包。
#### 数据结构

Redis数据库主要存储User及Message结构体，代码如下：

```golang
//functool.go


//Message -> Db 1
type Message struct {
    Name   string `json:"name"`
    IdTime string `json:"idtime"`
    Msge   string `json:"msge"`
}
//User ->Db 0
type User struct {
    Name   string `json:"name"`
    Passwd string `json:"passwd"`
    Status string `json:"status"`
}

//已登录用户
var UserMap = make(map[int]User)
```

#### User 管理

##### 序列化

因在数据库直接存储结构体会导致，存储后输出的结果无法反序列化，所以这里再存储之前会对数据进行序列化。
当然有序列化必然有反序列化，以便后期从Redis获取User。
        

```Golang
//user.go


//将User转化为[]byte
func StructTojson(user User) []byte { 
    //将User结构体序列化操作封装成方法
    data, err := json.Marshal(user) 
    //使用json.Marshal序列化User，data数据类型为[]byte
    CheckError(err)                 
    //error处理
    return data                     
    //返回序列化之后的User
}
//将user型的[]byte转化为User
func JsonTostruct(b []byte) User {
    var user User
    err := json.Unmarshal(b, &user)
    CheckError(err)
    return user
}


```

##### 链接Redis数据库
本文已将链接数据库操作打包为一个方法，直接调用 NewRedis() 即可返回数据库客户端。
```Golang
//functool.go

func NewRedis(db int) *redis.Client { //将数据库连接操作打包为方法使用newRdis(0)方法带入数据库名调用即可
    rdb := redis.NewClient(&redis.Options{
        Addr:     "pi4:6379", //数据库默认安装在开发机，监听localhost，默认端口为6379
        Password: "passwd",              // 使用设置好的redis密码
        DB:       db,                    // 设置要使用的数据库名
    })
    return rdb //返回数据库客户端
}

```
这里数据库地址为pi4，是因为树莓派4B的主机名（参考上文Manjaro-arm的安装）叫pi4，所以使用pi4这个地址即可。

```Golang
//functool.go

func TestlinkRedis() ({ //获取单独的一条
    rdb := NewRedis(1)                      //获取redis客户端，并连接到数据库1
    val, err := rdb.Get(ctx, "key").Result() //使用"key"来进行查询，key是一个字符串
    if nil == err {
         fmt.Println("你得到了数据：",val)
    }else{
        fmt.Println("获取数据失败")
    }
    rdb.Close()                             //关闭redis客户端
}

```
根据TestlinkRedis()方法我们可以得知，获取数据库客户端后可以获得一个 * redis.Client 对象，使用它的Get方法可以获取 "key"对应的"value" 。

##### 将用户User添加到数据库 0 里面
建立 SetUser(user User) bool 方法，我们将一个 User 传入，先判断是否存在该用户 ，如果不存在就保存该用户，存在就不保存，使用 booll 可以帮我们判断用户是否保存成功。
```Golang
//user.go

//保存该用户到数据库
func SetUser(user User) bool { //保存User
    //使用getUser进行查询
    _, b := GetUser(user)
    if b { //如果有此人就直接返回此人
        return false
    } else {
        rdb := NewRedis(0)                                          //连接数据库0
        err := rdb.Set(ctx, user.Name, StructTojson(user), 0).Err() //保存用户
        CheckError(err)
        rdb.Close()
        return true
    }
}
```
其中可以看到先使用了一个 GetUser() 方法来判断该用户是否存在过，目的是避免同一用户多次存储，造成新用户挤掉老用户的问题
其中还使用到上文所说的 StructTojson(User) 方法。
##### 从数据库 0 获取User
建立 GetUser(user User) (User, bool) 方法，可以看出，使用此方法需要一个 User 对象，因为我们使用 User.Name 作为唯一识别码，所以这里只需要 User.Name 带值即可。
```Golang
//user.go

//查询该用户，并返回相应数据
func GetUser(user User) (User, bool) { //获取用户全部的信息，这里使用user.Name作为唯一识别号
    rdb := NewRedis(0)                           //连接数据库
    val, err := rdb.Get(ctx, user.Name).Result() //使用user.Name来进行查找
    rdb.Close()
    if nil == err {
        return JsonTostruct([]byte(val)), true //如果有此用户，则返回User，true
    }
    return User{"nil", "nil", "nil"}, false //如果没有此用户，则返回nil，false
}
```
在上述代码中，从数据库 0 中获取回来的值 val 是一个 string 类型的数据，需要先转化为[]byte类型，再使用 JsonTostruct([]byte) 方法进行反序列化。
GetUser() 方法会返回 User，bool ，可以从bool的值来判断用户获取是否成功。
##### 代码测试
接下来使用一段代码来对上述内容进行测试。
```Golang
//functool.go

func TestsaveUser() { 
    user := User{
        Name:   "bill",
        Passwd: "123412sdd",
        Status: true,
    }
    b := SetUser(user)
    if b {
        fmt.Println("oh yeah we saved this user!")
    } else {
        fmt.Println("opps! we alredy have this user!dont save agn!")
    }
    return
}
func TestgetUser() {
    user := User{
        Name:   "bill",
        Passwd: "",
    }
    u, msg := GetUser(user)
    if msg {
        fmt.Println("yeah ~~ we have this user :", u)
    } else {
        fmt.Println("opps we dont have this user!")
    }
    return
}

//该方法可以删除Redis里任意数据库 Db 的任意 Key
func DeletSomething(key string, Db int) {
    rdb := NewRedis(Db)
    rdb.Del(ctx, key).Err()
    rdb.Close()
}
```
###### 运行 TestsaveUser() 方法

```Golang
func main() {
    TestsaveUser()
}
```
得到保存成功的信息。
```Bash
go run server.go functool.go
oh yeah we saved this user!
```
如再次运行该方法则会得到 我们已经有这个数据了，不要再来保存的提示。
```Bash
go run server.go functool.go
opps! we alredy have this user!dont save agn！
```
###### 运行 TestgetUser() 方法
```Golang
func main() {
    TestgetUser()
}
```
可以获取到该用户的全部信息。
```Bash
go run server.go functool.go
yeah ~~ we have this user : {bill 123412sdd true}
```
###### 运行DeletSomething()
```Golang
//functool.go
func main() {
    DeletSomething("bill",0)
}
```
没有返回信息，根据  rdb.Del(ctx, key).Err() 方法，在没有出错的情况下是不会有返回值的。
```Bash
go run server.go functool.go
```
可以再次运行 TestgetUser() 方法来查询名为 "bill" 用户。
```Bash
go run server.go functool.go
opps we dont have this user!
```
可以看出该用户已经不存在。

户管理功能，成功

* * *

#### Message 管理
##### 序列化 Message
```Golang
//messaage.go
data, err := json.Marshal(m)
```
使用 json.Marshal() 方法实现对 Message 的序列化。

##### 保存Message到数据库 1

```Golang
//message.go

//保存该message到数据库1
func SetMessage(m Message) bool {
    data, err := json.Marshal(m) //使用json.Marshal序列化，data数据类型为[]byte
    if nil == err {
        rdb := NewRedis(1)
        rdb.Set(ctx, m.IdTime, data, 24*time.Hour).Err() //保存Message,保存24小时24*time.Hour
        rdb.Close()
        return true //返回true表示存储成功
    } else {
        return false //返回false表示存储失败
    }
}
```
使用 SetMessage() 方法实现 Message 的存储。
其中 rdb.Set(ctx, msg.Id, data, 24 * time.Hour).Err() 这里的 24 * time.Hour 用来控制 Message 存在时间，服务端 Message 存储时限设置为24小时，过期自动删除，在后面的测试中应使用更短的时间例如 10 * time.Second 超过10秒后自动删除。

##### 从数据库 1 中获取 Message
```Golang
//functool.go

//使用IdTime进行message查询
func GetMessage(idtime string) (Message, bool) { //获取单独的一条
    rdb := NewRedis(1)                        //获取redis客户端
    val, err := rdb.Get(ctx, idtime).Result() //使用IdTime获取Message
    rdb.Close()                               //关闭redis客户端
    if nil == err {
        var msg Message
        err := json.Unmarshal([]byte(val), &msg) //反序列化
        CheckError(err)
        return msg, true
        //在无err情况下返回Message，并设置状态为true
        //true表示获取成功
    } else {
        return Message{"", "", ""}, false
        //在err情况下返回空的Message，并设置状态为false
        //false表示获取失败
    }
}
```
这个是获取单个 Message 的测试，在实际应用中， 服务端收到来自客户端的 Message 后将会把这条 Message 群发给所有在线的用户，随即保存到 Redis 1 里。
这里有一个用户场景，例如用户 bill 不在线，所以 bill 的状态是不在线，当 bill 重新上线后，需要历史消息，就会向客户端发起请求，这里就会用到 GetAllMessageRange() ，将所有的历史消息打包给用户 bill 。
```Golang
//functool.go

//获取该数据库里所有的key
func GetAllKeys(db int) []string {
    rdb := NewRedis(db)
    defer rdb.Close()
    keys, err := rdb.Keys(ctx, "*").Result()
    CheckError(err)
    return keys
}

//message.go


//获取全部的消息,将其打包为map
func GetAllMessageRange() []Message {
    var MessageSlice []Message
    b := GetAllKeys(1)
    for _, idtime := range b {
        m, _ := GetMessage(idtime)
        MessageSlice = append(MessageSlice, m)
    }
    return BubbleSortPro(MessageSlice)
}

//对Message进行冒泡排序,使其按照IdTime的先后顺序
func BubbleSortPro(arr []Message) []Message {
    length := len(arr)
    for i := 0; i < length; i++ {
        over := false
        for j := 0; j < length-i-1; j++ {
            if arr[j].IdTime > arr[j+1].IdTime {
                over = true
                arr[j], arr[j+1] = arr[j+1], arr[j]
            }
        }
        if !over {
            break
        }
    }
    return arr
}
```
以上代码中我们使用到了冒泡排序，这主要是因为 Redis 是无序数据库，所以需要特别的进行顺序排列，按时间戳 IdTime 进行排序 。有利于后期的 Message 遍历。
##### 删除
在上文中我们提到了， Message 在Redis数据库里超过24小时，将会自动删除，所以这里并不特别需要删除 Message 。
如需删除指定 Message 可以使用上文【用户管理】中提到的 DeletSomething(key string, Db int) 进行删除。

##### 测试

```Golang
//main.go


func main() {
    testSetMssage()
    for i, msg := range GetAllMessageRange() {
        u, err := strconv.ParseInt(msg.IdTime, 10, 64)
        CheckError(err)
        d := time.Unix(u/1e9, 0)
        fmt.Println("ID:", i, "时间：", d, "名字", msg.Name, "Message：", msg.Msge)
    }
}
func testSetMssage() {
    for j := 0; j < 10; j++ {
        s := strconv.Itoa(j)
        time.Sleep(1 * time.Second)
        user := User{
            "bill" + s,
            "",
            "",
        }
        
        t :=time.Now().UnixNano()
        msg := Message{
            user.Name,
            string(t),
            "message " + s,
        }
        SetMessage(msg)
    }
    fmt.Println("数据保存完毕！")
}


```
用以上代码进行测试，我们可以在控制台看到打印出来的10个 Message 。
```Bash
go run Server.go functool.go
数据保存完毕！
ID: 2 时间： 2021-04-23 13:48:51 +0800 CST 名字 bill2 Message： message 2
ID: 4 时间： 2021-04-23 13:48:53 +0800 CST 名字 bill4 Message： message 4
ID: 5 时间： 2021-04-23 13:48:54 +0800 CST 名字 bill5 Message： message 5
ID: 8 时间： 2021-04-23 13:48:57 +0800 CST 名字 bill8 Message： message 8
ID: 0 时间： 2021-04-23 13:48:49 +0800 CST 名字 bill0 Message： message 0
ID: 1 时间： 2021-04-23 13:48:50 +0800 CST 名字 bill1 Message： message 1
ID: 7 时间： 2021-04-23 13:48:56 +0800 CST 名字 bill7 Message： message 7
ID: 9 时间： 2021-04-23 13:48:58 +0800 CST 名字 bill9 Message： message 9
ID: 3 时间： 2021-04-23 13:48:52 +0800 CST 名字 bill3 Message： message 3
ID: 6 时间： 2021-04-23 13:48:55 +0800 CST 名字 bill6 Message： message 6
```
其中我们可以看见， map 里的数据是没有顺序的，因为 Redis 数据库是一个无序数据库，所以打印出来不是顺序排列的，但是我们可以简单的通过遍历将数据顺序排列出来。
在线用户收到消息必然是顺序排列的，因 Message 收到时先群发给所有在线用户，再使用 SetMessage() ，这个事件是有时间顺序的。
新上线用户，会收到历史 Message ，也就是使用上文说的冒泡排序，客户端收到后进行遍历，这样的顺序就是正常的了。

### 第二节 服务器的设计与搭建
#### http服务器的搭建
这里使用 Golang 自带的 "net/http" 包来搭建http服务器，http服务器将监听所有来自客户端的请求并返回相应数据。
```Golang
 //main.go
 
 //设置服务器监听地址与端口
var addr = flag.String("addr", ":8989", "http service address")
///服务器入口
func main() {
    flag.Parse()
    hub := newHub() //实例化hub
    go hub.run() //开启 hub
    http.HandleFunc("/message", func(rw http.ResponseWriter, r *http.Request) {
        message(hub, rw, r)
    })//监听 message聊天功能
    http.HandleFunc("/register", func(rw http.ResponseWriter, r *http.Request) {
        Register(rw, r)
    })//监听注册功能
    http.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
        Login(rw, r)
    })//监听登录功能
    http.HandleFunc("/dislogin", func(rw http.ResponseWriter, r *http.Request) {
        DisLogin(rw, r)
    })监听注销登录功能
    http.HandleFunc("/chanagepasswd", func(rw http.ResponseWriter, r *http.Request) {
        ChanagePasswd(rw, r)
    })//监听修改密码功能
    http.HandleFunc("/getallmessage", func(rw http.ResponseWriter, r *http.Request) {
        GetAllMsg(rw, r)
    })//监听获取24小时消息列表的功能
    err := http.ListenAndServe(*addr, nil)//服务器开始监听
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```
##### 注册功能的实现
注册功能，前端对数据进行简单校验后发送给服务器，服务器收到请求后进行相应处理。
```Golang
//client.go


func Register(rw http.ResponseWriter, r *http.Request) {
    s := ""
    data := User{}
    json.NewDecoder(r.Body).Decode(&data)
    b := SetUser(data)
    if b { //成功保存此用户到数据库,用户注册成功
        s = "yes"
    } else { //已存在,用户注册失败
        s = "no"
    }
    rw.Write([]byte(s))
}
```
将收到的数据解码为 User 调用 SetUser() 方法，进行用户的注册，返回相应信息给客户端。
##### 登录功能的实现

```Golang
//client.go

func Login(rw http.ResponseWriter, r *http.Request) {
    s := "defalut"
    data := User{}
    json.NewDecoder(r.Body).Decode(&data)
    user, b := GetUser(data)
    if b { //用户名正确
        if data.Passwd == user.Passwd { //密码正确,加入登录列表,登录成功
            s = "yes"
            if !IsOurUser(data.Name) {
                u1 := user
                u1.Name = data.Name //只保存用户名
                var l = len(UserMap)
                UserMap[l+1] = u1
            }
        } else { //密码错误,登录失败
            s = "no"
        }
    } else { //用户名错误
        s = "no"
    }
    rw.Write([]byte(s))
}

//functool.go

//判断是否存在该户,存在true,不存在false
func IsOurUser(name string) bool {
    bo := false
    for _, u1 := range UserMap {
        if u1.Name == name {
            bo = true
        }
    }
    return bo
}
```
将收到的 User 对其在数据库内进行查询。如果在数据库查询到该用户，用户名正确，进行密码的验证。如密码正确，将对客户端返回 "yes"，表示用户成功登录，反之返回 "no" 表示用户登录失败， 用户名或者密码错误。在成功登录后还会觉得是否加入登录列表，浙江决定用户是否能进行正常的聊天操作。确保用户没有登录的情况下是不能进行聊天操作的。
#####  注销登录功能的实现
```Golang
//client .go

func DisLogin(rw http.ResponseWriter, r *http.Request) {
    data := User{}
    json.NewDecoder(r.Body).Decode(&data)
    DeletLoginUser(data)
    rw.Write([]byte("yes"))
}

//functool.go

//将用户踢出用户Map,表示该用户没有登录
func DeletLoginUser(u User) {
    for i, user := range UserMap {
        if user.Name == u.Name {
            delete(UserMap, i)
        }
    }
}
```
服务器收到注销请求后，将用户踢出登录列表，在此功能里，如果用户在登录列表里将会成功注销，服务器返回 "yes" 。如果没有成功踢出登录列表，表示用户没有在登录列表却要注销，这种极端情况，我们也对客户端返回 "yes" ，因该用户并不存在于登录列表 UserMap里，所以不会对其进行删除操作。
##### 密码修改功能的实现
```Golang
//client.go

func ChanagePasswd(rw http.ResponseWriter, r *http.Request) {
    s := "defalut"
    data := User{}
    json.NewDecoder(r.Body).Decode(&data)
    if SetPasswd(data) {
        //修改成功
        s = "yes"
    } else {
        //修改失败,检查账号和密码,如还是失败联系管理员
        s = "no"
    }
    rw.Write([]byte(s))
}

//user.go

//修改密码
//为了不新增对象/结构体,这里约定客户端发来的用户信息如下:
// var user=User{
//  Name: "用户名",
//  Passwd: "旧密码",
//  Status: "新密码",
// }
func SetPasswd(user User) bool {
    uu, b := GetUser(user)
    if b {
        if user.Passwd == uu.Passwd {
            //密码正确
            u := User{
                Name:   user.Name,
                Passwd: user.Status,
                Status: "",
            }
            DeletSomething(u.Name, 0)
            if SetUser(u) {
                //修改成功
                fmt.Println("此u给i成功")
                return true
            } else {
                //修改失败
                return false
            }
        } else {
            //密码错误
            return false
        }
    } else {
        //没有这个人,用户名错误
        return false
    }
    // return true
}
```
服务端将用 User 解码，我们与客户端约定好用户修改密码的数据用以上注释内容的 user 定义，这样将不会产生新的数据结构。
用户信息到达后先对其进行数据可查询，用户名与旧密码正确后，才进行修改操作，将客户端发来的用户数据 User.Status 填到 User.Passwd ，先删除该用户再保存该用户，以免调用 SetUser() 方法报错。
如修改成功，服务器会发送 "yes" ，如果失败将发送 "no" ，表示用户提交的信息有误。
##### 获取消息列表
```Golang
//client.go

func GetAllMsg(rw http.ResponseWriter, r *http.Request) {
    //  IsOurUser()
    s := GetAllMessageRange()
    data, e := json.Marshal(s)
    CheckError(e)
    rw.Write(data)
}
```
该方法会将所有24小时内（因上文中的 Message 数据可的设计，保存时限限制为24小时）的消息打包为 Slice （在 Golang 里 Slice 即是数组）然后发送给客户端，客户端将其进行解析，生成历史消息列表。
#### WebSokcet服务器的搭建

这里参考 [Github/gorilla/websocket/examples/chat]( [下载地址](https://manjaro.org/download/)) 示例。
本文讲的是 基于 Socket 的聊天程序设计与实现，所以不会对 Socket 进行重复造轮子，本文会基于上述示例，进行理解与修改，使其符合我们的功能需求。

##### 消息队列

```Golang
//client.go

// 将连接升级为 webSocket 并将其加入消息队列
func message(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upGrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client
    go client.writePump()
    go client.readPump()
}


var upGrader = websocket.Upgrader{
    ReadBufferSize:  1024,//读取数据大小
    WriteBufferSize: 1024,//写入数据大小
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
```
<u>本文使用 Socket 技术作为聊天功能实现技术，那么先将请求升级为 WebSocket ，使用 websocket.Upgrader 进行升级。</u>
升级后，实例化客户端 client ，并对其添加数据。使用管道 channel 将客户端加入 hub.register，进入到队列里。
分别开启该客户端的读、写 线程。至此，该用户处于等待收发消息的状态。

##### client.readPump 客户端消息读取

```Golang
//client.go

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        mt, message, err := c.conn.ReadMessage() //读取客户端发送的数据
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        var msg Message
        e := json.Unmarshal(message, &msg)
        CheckError(e)
        if !IsOurUser(msg.Name) {
            SendToClient(c.conn, mt, "You are not Login Users")
            break
        }
        SetMessage(msg)
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
        c.hub.broadcast <- message
    }
}
```
使用 defer func() 在该线程运行最后执行客户端 conn.close 关闭，并从 hub中踢出 消息队列 。
在开始该 client 客户端的线程之前，先设置消息长度限制、等待时间等。
开始线程后监听 ReadMessage ，当客户端发送消息后，这里将消息读取出来，message 的类型是 []byte 型。
错误处理：如果发生了错误，那么打印错误信息，并打断 for 循环，后 defer 开始进行收尾处理，关闭 Socket 连接。
如果一切正常，下一步就是将消息解析出来，使用 IsOurUser(）来判断是否是登录客户，如果返回 false 代表，该用户为登录，或者信息有误，导致判断不是登录 User，简而言之 返回 false 代表该 message 不是我们约定好的，那么将打断 for 循环，关闭该连接。
如过是我们登录的 User 那么将该条 message 保存到数据库。
继续，该条消息加入到 hub.go 第42行 。
hub.go 第42行 会对其客户端进行循环遍历，直到收到 message   ，收到后将其加入到 client.send里，而这时轮到 client.writePump 出场了。
##### client.writePump 客户端消息写入
```Golang
//client.go

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)//设置定时器
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // The hub closed the channel.
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)
            // Add queued chat messages to the current websocket message.
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write(newline)
                w.Write(<-c.send)
            }
            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}


```
使用 defer func() 在该线程运行最后执行客户端 conn.close 关闭，定时器关闭，并从 hub中踢出 消息队列 。
在 for 循环内 select {} 它将会从第一个可运行的 case 开始运行，直到没有可以运行的 case ，它就开始阻塞， default 总会运行。
如果收到了消息，那么第一个 case 将运行，第一个 case 运行的时候，会将客户端发送的 message 发送该客户端。
因该线程是每一个客户端都会有的，所以接收到 message 后，每一个在线的客户端都会收到该条 message。
如果没有收到 message 那么将会运行 定时器，截至时间已过，将会重新等待下一个 message 的到来。
至此，消息收发功能完成。

### 第三节 前端的设计与搭建

#### 项目文件结构
##### SocketuView/
- .hbuilderx
- pages
- Changeasswordp
    - Changeasswordp.vue /密码修改页面
- login
    - login.vue /登录页面
- message
    - message.vue /消息，聊天页面
- register
    - register.vue /注册页面
- setting
    - setting.vue /设置页面
- log
    - log.vue /log页面
- static
- uni_modules 
- unpackage 
- uview
- eslintignore
- App.vue /app 的入口
- LICENSE
- main.js
- manifest.json
- pages.json
- uni.scss
- vue.config.js

可以看见的就是主要结构及文件，其余不在我们这个项目的讨论范围内
#### 建立连接

在前端中将会使用 http请求 以及 websocket 与服务端进行连接。
##### http请求
```js
uni.request({
    url:getApp().globalData.serverAddr+'dislogin',
    method:'POST',
    data:JSON.stringify(user),
    success: (res) => {
        getApp().globalData.SetLog('dis login success',res)
    },
    fail: (res) => {
        getApp().globalData.SetLog('dis login dail',res)
    }
})
```
http请求使用 uni-app 的 uni.request() API ，使用非常简单，填入请求地址，方法类型，要发送的数据。
请求成功将会在 success 里面接收数据，进行下一步处理，例如登录功能的请求成功方法如下：
```js
switch(res.data){
    case 'yes':									
        try{
            uni.setStorageSync('user',this.user)
        }catch(e){
            getApp().globalData.SetLog('login get user fail',e)
        }
        uni.reLaunch({
            url:'../message/message'
        });
        uni.showToast({
            title: '登录成功!',
            duration: 2000
        });
        getApp().globalData.SetLog('login success','yes')
    break;
    case 'no':
        uni.showToast({
            title: '信息有误!',
            duration: 2000
        });
        getApp().globalData.SetLog('login fail','Wrong user name or password ')
    break;
    default:
        getApp().globalData.SetLog('login fail',res)
        uni.showToast({
            title: '网络开小差了!',
            duration: 2000
        });
}
```
请求成功后对收到的数据 res进行解析，而 res 接收到的数据是 json 类型的，不过 uni.request 如果发送的数据类型是 json 那么它会自动尝试解析 JSON.parse 。
在服务端，login 请求的返回值只有 "yes" or "no" 两种情况，所以写两个 case 加一个 default 即可，如果服务器返回 "yes" 那么将使用 uni.setStorageSync 同步缓存 ，把 这个 this.user 保存到本地缓存中，这样下一次打开 APP 将会自动填充账号。

##### websocket连接
这里使用 uni-app 的 uni.connectSocket() API ,因 uni-app 的 socket 是单一 且唯一的连接，所以我们在 App.vue 创建全局 socket 连接，这样在 app 内任意页面都可以管理这个 socket 连接，使用  getApp().globalData.SocketTask 获取 SocketTask 对象。
```js
//App.vue

globalData:{
    SocketTask:null,
    initSocketTask(){
        getApp().globalData.SocketTask = uni.connectSocket({
        url:getApp().globalData.serverAddr+'message',
        success:(res)=>{
            getApp().globalData.SetLog('message init task success',res)
        }
        });
    },
}
```
例如在 message.vue 使用：

```js
initTask(){
    getApp().globalData.initSocketTask()//初始化socket
    this.Task = getApp().globalData.SocketTask//获取socketTask
    this.Task.onError((res)=>{
        getApp().globalData.SetLog('message init task onErro',res)
    });
    this.Task.onClose((res)=>{
        getApp().globalData.SetLog('message init task onClose',res)
    });
    this.Task.onOpen((res)=>{
        getApp().globalData.SetLog('message init task onOpen',res)
    });
    this.Task.onMessage((res)=>{//监听服务武器返回消息
        getApp().globalData.SetLog('message init task onMessage','成功')
        if(null == this.list){//防止消息列表为空，服务器里没有任何一条 message 导致出错
            var s = [JSON.parse(res.data)]//新建一个数组来保存这条 message
            this.list = s
        }else{//如果消息列表不为空
            this.list.push(JSON.parse(res.data))//使用push()将这条 messgae 放入数组最后，如果数组为空将不能使用push()
        }
    });
},
```
以上只提到了接收 message 那么以下就是发送 message 的实例：
```js
sendGo(){
    if(''!=this.msg){//判断发送的 message 是否为空空
        let d = new Date().getTime().toString()//获取时间
        var msgData ={ //拼接 message 对象
            name:this.user.name,
            idtime:d,
            msge:this.msg
        }
        this.Task.send({
            data:JSON.stringify(msgData),
            success:(res)=>{
                this.msg = ''
            },
            fail:(res)=>{
                getApp().globalData.SetLog('message send msg fail',res)
            }
        });
    }else{
        uni.showToast({
            title:'消息不能为空',
            duration:2000
            })
    }	
},
```
#### 前端页面的实现
##### 为什么使用uView

在21世纪的今天，前端页面不需要开发者一个个 css 去写，优秀、好看、易用的ui框架层出不穷，而本文中使用 uniapp 作为前端app的开发框架，而 [uView](uviewui.com) 正是 uniapp 里的ui框架，所以结论显而易见。
##### 代码部分
具体代码这里不过多阐述，整个项目，前端代码，包括后台服务器代码将上传 GitHub ，如有需要请移步 [github/brejce/](https://github.com/brejce/SocketChat)
