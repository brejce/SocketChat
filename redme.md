# 基于Socket的聊天程序设计与实现

## 摘要

本论文主要阐述了为什么使用Golang与uni-app及redis数据库。以及如何使用他们来作为Socket通信服务。
随着21世纪的信息化的快速发展与迭代。

Golang作为一款高性能的静态语言（性能仅次于C/C++），Go主要有静态语言、天生并发、内置GC、安全性高、语法简单、交叉编译和编译快速这几个方面的特性，及较低的入门门槛，无论之前学习过什么语言都能快速的上手。
总结下来就是：运行快、开发快和部署快。

uni-app 是一个使用 Vue.js 开发所有前端应用的框架，开发者编写一套代码。

可发布到iOS、Android、Web（响应式）、以及各种小程序（微信/支付宝/百度/头条/QQ/钉钉/淘宝）、快应用等多个平台。DCloud公司拥有800万开发者、数百万应用、12亿手机端月活用户、数千款uni-app插件、70+微信/qq群。阿里小程序工具官方内置uni-app（详见），腾讯课堂官方为uni-app录制培训课程（详见），开发者可以放心选择。uni-app在手，做啥都不愁。即使不跨端，uni-app也是更好的小程序开发框架（详见）、更好的App跨平台框架、更方便的H5开发框架。不管领导安排什么样的项目，你都可以快速交付，不需要转换开发思维、不需要更改开发习惯。

Redis是Remote Dictionary Server(远程数据服务)的缩写，由意大利人antirez(Salvatore Sanfilippo)开发的一款内存高速缓存数据库，该软件使用C语言编写，它的数据模型为key-value。
可持久化(一边运行，一边把数据往硬盘中备份一份，防止断电等情况导致数据丢失，等断电情况恢复之后，Redis再把硬盘中的数据恢复到内存中)，保证了数据的安全。

在本文中，由于聊天应用的数据简单，使用Redis可方便快捷的存取数据，从而快速开发。

关键词:  Golang    Vue.js    Redis    套接字 聊天应用

## Abstract

This paper mainly discusses why Golang and uni-app And Redis database. And how to use them as socket communications services.

Golang is a high-performance static language (performance is only second to C/C++), go mainly has the characteristics of static language, natural concurrency, built-in GC, high security, simple syntax, cross compilation and compilation, and low threshold for entry. No matter what language you have learned before。
it can be quickly started. The summary is: fast operation, rapid development and rapid deployment.

Uni app is a framework for developing all front-end applications using vue.js. Developers write a set of code.

Can be released to iOS, Android, Web (Responsive), and various small programs (WeChat / Alipay / Baidu / headline /QQ/ nail / Taobao), fast application and other platforms. Dcloud has 8 million developers, millions of applications, 1.2 billion mobile end-month users, thousands of uni app plug-ins, and 70+ wechat / QQ group. Alibaba applet tool officially built-in uni app (see details), Tencent classroom official recording training course uni app (see details), developers can choose at ease. Uni app is in the hand, and it doesn't worry about doing anything. Even if not cross-end, uni app is also a better framework for small program development (see details), a better cross platform framework for app, and a more convenient H5 development framework. No matter what kind of project the leader arranges, you can deliver quickly without changing the development thinking and changing the development habits.

Redis is the abbreviation of remote dictionary server (remote data service). It is a memory cache database developed by Italian anti-sanchilipo. The software is written in C language, and its data model is key-value. 
It can be persisted (while running, backup data to the hard disk to prevent data loss due to power failure, Redis will restore the data in the hard disk to memory after the power failure is restored), which ensures the data security.

In this paper, Because the data of chat application is simple, Redis can easily and quickly access data, and develop it quickly.

Key words: Golang Vue.js Redis Socket ChatApplication

## 目录

这是目录
[toc]

## 第一章 序言

### 背景分析

美国著名心理学家亚伯拉罕·马斯洛将人类的需求划分为五大层次，它包括生理的需要、安全的需要、社交的需要、尊重的需要和自我实现的需要。我认为他的划分在分类上是比较严谨的，但是在通俗上来说只需要：生理需求--精神需求，两个大类.生理需求是一个人最基本的生存生活需求，主要有：食、性、衣、住、行五个方面。不用多说。精神需求主要包括：安全需求、社交需求、荣誉需求、倾诉需求、认同需求、艺术需求等等。通常，在满足生理需求的基础上，人们才产生精神需求，但有时会与生理需求交叉。总体上，精神需求层次要高于生理需求。人类目前的生存环境，物质需求基本已经满足，所以现在的人类都在最求精神上的满足，目前最热门的精神需求莫过于网络聊天/网络直播。网上聊天的真正魅力应该来源于文字对情感的渗透，我们眼睛里看到的世界是有限的，而心灵感受到的思想才是无限的永恒的，现实生活中我们了解一个人，首先是看见，认识和沟通，网上则不然，我们用眼睛看不见对方，也无法理性的去认识你对面的人，大家相互所感受到的，是通过文字，通过思想情感的交流，循序渐进的去挖掘去了解他或者她。

在21世纪里，国内主要是QQ/微信/微博。国外主要是FaceBook/Twiter/ Instagram。
在本文中，就是讲如何制作一个简单的聊天程序。

### 使用技术、设备与面向平台

#### 使用语言、开发工具

Golang、Vue.js、HTML、js等语言。HbuilderX、VsCode等开发工具。Windows、Linux、Android等系统。

#### 使用设备


搭载Windows10系统的台式机、搭载Manjaro-arm系统的树莓派4B，搭载Android系统的荣耀、红米等手机
#### 资料库
百度（www.baidu.com），Github（hithub.com），ArchWiki（wiki.archlinux.org），uni-app（uniapp.dcloud.io/README），Vue（cn.vuejs.org/v2/guide/），菜鸟教程（www.runoob.com）

#### 面向平台


服务端：该聊天程序服务端运行与电脑上，系统上无特别需求，Winsows/Linux/Mac均可作为服务端的运行环境。

客户端：客户端主要以Vue.js为核心技术的uni-app通过HbuilderX开发工具，可发布于各平台（iOS、Android、Web（响应式）、以及各种小程序（微信/支付宝/百度/头条/QQ/钉钉/淘宝）、快应用等多个平台）以获得跨平台的聊天体验，

在本文中客户端运行在安卓平台上。


## 第二章 技术简介

### 第一节 前端框架介绍

#### HTML 简介


 HTML的全称为超文本标记语言，是一种标记语言。它包括一系列标签．通过这些标签可以将网络上的文档格式统一，使分散的Internet资源连接为一个逻辑整体。

HTML文本是由HTML命令组成的描述性文本，HTML命令可以说明文字，图形、动画、声音、表格、链接等。超文本是一种组织信息的方式，它通过超级链接方法将文本中的文字、图表与其他信息媒体相关联。这些相互关联的信息媒体可能在同一文本中，也可能是其他文件，或是地理位置相距遥远的某台计算机上的文件。这种组织信息方式将分布在不同位置的信息资源用随机方式进行连接，为人们查找，检索信息提供方便。


#### Javascript 简介


JavaScript（简称“JS”） 是一种具有函数优先的轻量级，解释型或即时编译型的编程语言。虽然它是作为开发Web页面的脚本语言而出名，但是它也被用到了很多非浏览器环境中，JavaScript 基于原型编程、多范式的动态脚本语言，并且支持面向对象、命令式和声明式（如函数式编程）风格。

JavaScript在1995年由Netscape公司的Brendan Eich，在网景导航者浏览器上首次设计实现而成。因为Netscape与Sun合作，Netscape管理层希望它外观看起来像Java，因此取名为JavaScript。但实际上它的语法风格与Self及Scheme较为接近。JavaScript的标准是ECMAScript 。截至 2012 年，所有浏览器都完整的支持ECMAScript 5.1，旧版本的浏览器至少支持ECMAScript 3 标准。2015年6月17日，ECMA国际组织发布了ECMAScript的第六版，该版本正式名称为 ECMAScript 2015，但通常被称为ECMAScript 6 或者ES2015。


#### Vue 简介


Vue (读音 /vjuː/，类似于 view) 是一套用于构建用户界面的渐进式框架。

与其它大型框架不同的是，Vue 被设计为可以自底向上逐层应用。Vue 的核心库只关注视图层，不仅易于上手，还便于与第三方库或既有项目整合。另一方面，当与现代化的工具链以及各种支持类库结合使用时，Vue 也完全能够为复杂的单页应用提供驱动。


#### uni-app 简介


uni-app 是一个使用 Vue.js 开发所有前端应用的框架，开发者编写一套代码。

可发布到iOS、Android、Web（响应式）、以及各种小程序（微信/支付宝/百度/头条/QQ/钉钉/淘宝）、快应用等多个平台。DCloud公司拥有800万开发者、数百万应用、12亿手机端月活用户、数千款uni-app插件、70+微信/qq群。阿里小程序工具官方内置uni-app（详见），腾讯课堂官方为uni-app录制培训课程（详见），开发者可以放心选择。uni-app在手，做啥都不愁。即使不跨端，uni-app也是更好的小程序开发框架（详见）、更好的App跨平台框架、更方便的H5开发框架。不管领导安排什么样的项目，你都可以快速交付，不需要转换开发思维、不需要更改开发习惯。


### 第二节 后端框架介绍

#### Golang 简介


Go 编程语言是一个开源项目，它使程序员更具生产力。

Go 语言具有很强的表达能力，它简洁、清晰而高效。得益于其并发机制，用它编写的程序能够非常有效地利用多核与联网的计算机，其新颖的类型系统则使程序结构变得灵活而模块化。 Go 代码编译成机器码不仅非常迅速，还具有方便的垃圾收集机制和强大的运行时反射机制。它是一个快速的、静态类型的编译型语言，感觉却像动态类型的解释型语言。


### 第三节 数据库介绍

#### Redis数据结构服务器简介


REmote DIctionary Server(Redis) 是一个由 Salvatore Sanfilippo 写的 key-value 存储系统，是跨平台的非关系型数据库。Redis 是一个开源的使用 ANSI C 语言编写、遵守 BSD 协议、支持网络、可基于内存、分布式、可选持久性的键值对(Key-Value)存储数据库，并提供多种语言的 API。

Redis 通常被称为数据结构服务器，因为值（value）可以是字符串(String)、哈希(Hash)、列表(list)、集合(sets)和有序集合(sorted sets)等类型。


### 第四节 开发工具及开发设备介绍

#### VsCode 简介


Visual Studio Code是一个轻量级但强大的源代码编辑器，运行在您的桌面上，可用于窗口，macOS和Linux。它配备了对 JavaScript、TypeScript 和节点的内置支持.js并且具有丰富的其他语言扩展生态系统（如C++、C#、Java、Python、PHP、Go）和运行时间（如 。NET 和 Unity）。


#### HBuilderX 简介


HBuilderX，H是HTML的首字母，Builder是构造者，X是HBuilder的一代版本。我们也简称HX。 HX是轻如编辑器、强如IDE的合体版本。


1.    轻巧 仅10余M的绿色发行包(不含插件)。
2.    极速 不管是启动速度、大文档打开速度、编码提示，都极速响应 C++的架构性能远超Java或Electron架构。
3.    vue开发强化 HX对vue做了大量优化投入，开发体验远超其他开发工具。
4.    小程序支持 国外开发工具没有对中国的小程序开发优化，HX可新建uni-app或小程序、快应用等项目，为国人提供更高效工具。
5.    markdown利器 HX是唯一一个新建文件默认类型是markdown的编辑器，也是对md支持最强的编辑器 HX为md强化了众多功能，请务必点击【菜单-帮助-markdown语法示例】，快速掌握md及HX的强化技巧！
6.    清爽护眼 HX的界面比其他工具更清爽简洁，绿柔主题经过科学的脑疲劳测试，是最适合人眼长期观看的主题界面 。
7.    强大的语法提示 HX是中国唯一一家拥有自主IDE语法分析引擎的公司，对前端语言提供准确的代码提示和转到定义(Alt+鼠标左键)。
8.    高效极客工具 更强大的多光标、智能双击...让字处理的效率大幅提升。
更强的json支持 现代js开发中大量json结构的写法，HX提供了比其他工具更高效的操作。

本文中，会使用 HBuilderX 作为前端开发工具。

#### Manjaro-arm 简介


Manjaro是一款基于Arch Linux、对用户友好、全球排名第1的Linux发行版（排名数据源于DistroWatch，统计日期2018.03.02，时间段3个月。）

在Linux圈，Arch的确是一个异常强大的发行版。它有3个无与伦比的优势：


1. 
   滚动更新可以使软件保持最新。
   
2. 
   AUR软件仓库有着世界上最齐全的Linux软件。
   
3. 
   丰富的wiki和活跃的社区让所有问题都可以快速得到满意的答案。
   


本文使用的树莓派4B就是搭载Manjaro-arm，支持ARM架构的Manjaro系统。


#### Windows系统 简介


Microsoft Windows操作系统是美国微软公司研发的一套操作系统，它问世于1985年，起初仅仅是MS-DOS模拟环境，后续的系统版本由于微软不断的更新升级，不但易用，也成为了当前应用最广泛的操作系统。

Windows采用了图形用户界面（GUI），比起从前的MS-DOS需要输入指令使用的方式更为人性化。

随着计算机硬件和软件的不断升级，Windows也在不断升级，从架构的16位、32位再到64位，系统版本从最初的Windows 1.0到大家熟知的Windows 95、Windows 98、Windows 2000、Windows XP、Windows Vista、Windows 7、Windows 8、Windows 8.1、Windows 10和Windows Server服务器企业级操作系统，微软一直在致力于Windows操作系统的开发和完善。


#### 树莓派4B 简介


树莓派由注册于英国的慈善组织“Raspberry Pi 基金会”开发，Eben·Upton/埃·厄普顿为项目带头人。

2012年3月，英国剑桥大学埃本·阿普顿（Eben Epton）正式发售世界上最小的台式机，又称卡片式电脑，外形只有信用卡大小，却具有电脑的所有基本功能，这就是Raspberry Pi电脑板，中文译名"树莓派"。自问世以来，受众多计算机发烧友和创客的追捧，曾经一“派”难求。

别看其外表“娇小”，内“心”却很强大，视频、音频等功能通通皆有，可谓是“麻雀虽小，五脏俱全”。自从树莓派问世以来，经历了A型、A+型、B型、B+型、2B型、3B型、3B+型、4B型等型号的演进。2019年6月25日，树莓派基金会宣布树莓派4B版本发布。

在本文中，会使用该设备进行服务端的设计与搭建。


## 第三章 开发环境介绍与搭建


### Windows 开发环境


关于Windows10系统的安装本文不做阐述，网络上有大量的教程文章。

Windows是普通人群使用最多的系统，熟悉度非常高，毕竟如今单机、网游、直播盛行，生态环境良好。本文中Windows系统将作为客户端uni-app的开发环境。


#### HBuilderX开发工具安装


HBuilderX官网下载地址：[下载地址](https://www.dcloud.io/hbuilderx.html)     。

下载之后直接解压到    D:\HBuilderX  。
进入其中运行<u>HBuilderX.exe</u>即可。



### Linux 开发环境 


本文中使用的Linux系统的发行版是 Manjaro-arm        [下载地址](https://manjaro.org/download/)        ，使用的设备为树莓派4B，在树莓派4B上安装 Manjaro-arm 需要下载特定版本，打开下载地址依次点击        Editions->ARM->Raspberry Pi 4接着选择喜欢的桌面环境版本，本文使用的是 Raspberry Pi 4 KDE Plasma21.02 版本。

系统安装过程 [可参考此地址](https://chenkaihua.com/2020/08/28/pi4-install-manjaro/) 进行安装配置。

系统安装、系统升级完成后需要安装 Golang、VsCode、Redis。


```Bash

# 使用以下代码进行安装即可
sudo pacman -S code
sudo pacman -S Golang
sudo pacman -S Redis

```

本文中 Linux 平台将作为 Golang 服务端开发环境及 Redis 数据库运行环境。


#### Golang语言开发环境安装


Golang 语言安装  [参考](https://www.runoob.com/go/go-environment.html)  。

#### VsCode开发工具的安装


详细安装过程 [参考](https://code.visualstudio.com/docs/setup/windows) 。
安装完成打开 VsCode 需要配置插件。

Chinese (Simplified) Language Pack for Visual Studio Code （汉化VsCode）。
Go （开启VsCode对Golang的支持）。

打开 VsCode 开始 Golang 之旅。
#### Redis数据库的安装

Redis数据库安装参考 [Github](https://github.com/redis/redis) 安装后运行即可。

本文中的开发环境搭建仅仅只做参考，网络上有大量的安装、配置教程，这里就不逐一展开。


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
