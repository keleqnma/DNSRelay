[TOC]

## 1 DNS协议

### 1.1 DNS packet组成

![img](http://www.tcpipguide.com/free/diagrams/dnsgenformat.png)



#### 1.1.1 Header

![img](http://www.tcpipguide.com/free/diagrams/dnsheaderformat.png)

1. **会话标识（2字节）：**是DNS报文的ID标识，对于请求报文和其对应的应答报文，这个字段是相同的，通过它可以区分DNS应答报文是哪个请求的响应

2. 标志（2字节）：

   | QR（1bit）     | 查询/响应标志，0为查询，1为响应                              |
   | -------------- | ------------------------------------------------------------ |
   | opcode（4bit） | 0表示标准查询，1表示反向查询，2表示服务器状态请求            |
   | AA（1bit）     | 表示授权回答                                                 |
   | TC（1bit）     | 表示可截断的                                                 |
   | RD（1bit）     | 表示期望递归                                                 |
   | RA（1bit）     | 表示可用递归                                                 |
   | rcode（4bit）  | 表示返回码，0表示没有差错，3表示名字差错，2表示服务器错误（Server Failure） |

3. **数量字段（总共8字节）：**Questions、Answer RRs、Authority RRs、Additional RRs 各自表示后面的四个区域的数目。Questions表示查询问题区域节的数量，Answers表示回答区域的数量，Authoritative namesversers表示授权区域的数量，Additional recoreds表示附加区域的数量

#### 1.1.2 Question

![img](http://www.tcpipguide.com/free/diagrams/dnsquestionformat.png)

1. **查询名：**长度不固定，且不使用填充字节，一般该字段表示的就是需要查询的域名

2. **查询类型:**

| 类型 | 助记符 | 说明                                     |
| ---- | ------ | ---------------------------------------- |
| 1    | A      | 由域名获得IPv4地址（本中继服务支持类型） |
| 2    | NS     | 查询域名服务器                           |
| 5    | CNAME  | 查询规范名称                             |
| 12   | PTR    | 把IP地址转换成域名                       |
| 28   | AAAA   | 由域名获得IPv6地址                       |

3. **查询类：**通常为1，表明是Internet数据。

#### 1.1.3 Answer

![img](http://www.tcpipguide.com/free/diagrams/dnsrrformat.png)

**1. 域名（2字节或不定长）：**它的格式和Question区域的查询名字字段是一样的。

**2.  查询类型：**表明资源纪录的类型，见1.2节的查询类型表格所示 

**3. 查询类：**对于Internet信息，总是IN

**4.  生存时间（TTL）：**以秒为单位，表示的是资源记录的生命周期，

**5. 资源数据：**该字段是一个可变长字段，表示按照查询段的要求返回的相关资源记录的数据。可以是Address（表明查询报文想要的回应是一个IP地址）或者CNAME（表明查询报文想要的回应是一个规范主机名）等。在本中继服务中，本地返回数据一律为Ipv4格式的地址，所以Answer长度固定为16字节。

### 1.2 实例

使用 **tcpdump** 监听本机53端口上的dns查询操作，使用 **nslookup** 查询 **baidu.com**。

可以看到nslookup指令向本机53端口上发送了一条长度为27 bytes *(Header length 12 bytes, Question length 15 bytes，其中question name长度11（结束符一个byte，打包baidu.com域名10个byte），question type和question class 4 byte)* 的查询域名数据。

随后本机53端口回复了一条长度为59 bytes *(Header length 12 bytes, Question length 15 bytes, 两个长度为16byte的answer)* 的结果数据。

![image-20200914182431542](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200914182431542.png)

![image-20200914183320154](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200914183320154.png)



## 2 DNS 中继程序

### 2.1 整体架构

![image-20200916144926577](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916144926577.png)

>**Feature:**
>
>1. 多携程安全。
>2. 转发查询后本地加载Redis缓存，Bench测试表明查询响应速度优于DNS默认查询1.5倍，优于转发查询2.5倍。
>3. 支持自定义屏蔽IP
>4. 支持多条查询IP返回

#### 2.1.1 工作逻辑

这里没有用到互斥锁或者线程池，因为Goroutine属于两级线程，创建和调度都不需要进入内核，分配回收代价极低，几百几千QPS的情况下完全够用了。频繁请求互斥锁会极大增多系统上下文切换，UDP请求这种非流式协议就不用互斥锁了。

这里用到Redis除了性能的考虑，也是因为域名-IP地址这种关系非常适合key-value这种数据。

![未命名文件](/Users/chenqiqi/Downloads/未命名文件.png)

#### 2.1.2 Bench测试

**Bench测试结果：**

> 查询DNS中继服务 查询( 10000 )次 花费时间( 101 )秒
> 查询转发服务114.114.114.114 查询( 10000 )次 花费时间( 253 )秒
> 查询默认DNS服务 查询( 10000 )次 花费时间( 142 )秒

Bench测试参数说明：

> ./bench.sh 10000 127.0.0.1 114.114.114.114

第一个参数是查询的次数，第二个参数是DNS中继服务的IP，第三个参数是DNS转发服务的IP。

### 2.2 本地查询

> **Feature：**
>
> 1. 本地查询仅支持Ipv4格式。
> 2. 支持多条查询结果。

![image-20200916115152839](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916115152839.png)

使用Nslookup 指令在127.0.0.1服务器上查询baidu.com域名，收到回复。

![image-20200916115217137](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916115217137.png)

使用tcpdump监听本机53端口上的dns查询操作，可以看到Nslookup指令向本机53端口上发送了一条长度为27的查询域名数据，随后本机53端口回复了一条长度为59的结果数据。

![image-20200916115246039](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916115246039.png)

下面是本DNS中继服务的运行日志，可以看到本地搜索到该域名共对应两条ip结果，随后中继服务将查询结果打包发回。

### 2.3 屏蔽查询

> **Feature：**
>
> 1. 屏蔽IP可以自由设置
>
> 2. 转发查询和本地查询的结果IP都可以屏蔽

![image-20200916115732742](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916115732742.png)

![image-20200916115646332](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916115646332.png)

屏蔽IP写入配置文件，中继服务初始化时会读取配置文件并输出日志。

#### 2.3.1 转发查询屏蔽

![image-20200916121302062](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916121302062.png)

使用Nslookup 指令在127.0.0.1服务器上查询alimama.com域名。

![image-20200916121322559](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916121322559.png)

根据中继服务服务日志可以看到，服务在本地没有搜索到alimama.com域名的解析结果，于是转发查询，解析转发查询结果发现其中有一个查询结果ip在本地的屏蔽列表中，于是对这条查询进行屏蔽处理。

![image-20200916121348102](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916121348102.png)



#### 2.3.2 本地查询屏蔽

![image-20200916121428061](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916121428061.png)

使用Nslookup 指令在127.0.0.1服务器上查询test.com域名。

![image-20200916121442117](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916121442117.png)

根据中继服务服务日志可以看到，服务在本地查询到了test.com域名的解析结果，发现其中有一个查询结果ip在本地的屏蔽列表中，于是对这条查询进行屏蔽处理。

![image-20200916121455045](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916121455045.png)



### 2.4 转发查询

> **Feature:**
>
> 转发查询后自动将查询结果加载到本地Redis缓存（只支持Ipv4类型）

**第一次查询：**

![image-20200916124949242](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916124949242.png)

使用Nslookup 指令在127.0.0.1服务器上查询youtube.com域名。

![image-20200916125003267](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916125003267.png)

根据中继服务服务日志可以看到，服务在本地没有搜索到youtube.com域名的解析结果，于是转发查询，将查询结果加载到本地Redis数缓存。

![image-20200916125734074](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916125734074.png)

加载缓存后，中继服务返回查询数据。

![image-20200916125019489](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916125019489.png)

**第二次查询：**

![image-20200916125041605](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916125041605.png)

使用Nslookup 指令在127.0.0.1服务器上查询youtube.com域名。

![image-20200916125521598](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916125521598.png)

根据中继服务服务日志可以看到，服务在本地没有查询到到了youtube.com域名的解析结果（上次加载进来的），于是打包返回。

![image-20200916125448416](/Users/chenqiqi/Library/Application Support/typora-user-images/image-20200916125448416.png)



