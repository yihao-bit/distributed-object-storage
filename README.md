# distributed-object-storage
分布式对象存储项目



1.实现接口服务和数据服务分离的分布式对象存储，提供了云存储后端的存储服务。

-  创建了一个名为apiServers的exchange，每一台数据服务节点都会持续向这个exchange发送心跳消息，每过5s发送一次心跳 
-  使用了REST和消息队列这两种不同类型的接口的原因：REST适合发送一个大数据量传输，消息队列适合群发。

2.使用 Golang 实现，使用 RESTful 架构。 

- 接口服务层提供了对外的 REST 接口，而数据服务层则提供数据的存储功能。数据服务本身也向接口服务提供REST接口。
- 使用到了Golang标准库，如 net/http ，net/url ，strings。使用了第三方开源库：github.com/klauspost/reedsolomon，github.com/streadway/amqp，github.com/klauspost/cpuid/v2。

3.接口节点和数据节点通过 RabbitMQ 消息队列进行通信；通过心跳机制，使接口节点感知数据节点的存在。 

- 当接口服务需要定位时，会创建一个临时消息队列，然后发送一条消息给dataServers exchange。内容是需要定位的对象，临时消息队列的名字是ip+端口
- 数据服务节点两个goruntine，一个发送心跳消息，一个监听定位消息
- map不是并发安全的，要用互斥锁
- 接口服务用goruntine启动一个匿名函数，用于1s后关闭临时消息队列，避免无休止等待
- 用io.Pipe创建了一对reader和writer

4.提供对象的元数据服务，使用 ElasticSearch 实现。

- 元数据指的是对象的描述信息 ，如名字，版本，散列值，大小。
- 散列值作为全局唯一标识符,散列函数是SHA-256,256的二进制数字作为散列值。
- 使用es，索引即数据库，类型即表，主分片数量固定。Digest:SHA-256=<BASE64编码>
- metadata索引使用的**映射**（mappings）结构，**映射**相当于表结构
- name属性有个额外的要求”index”：”not_analyzed“，精准匹配
- 要想获取对象当前最新版本的元数据需要使用ES搜索API

5.数据去重和数据校验功能，用hash值比对，sha-256 算法计算哈希值比对用户哈希值实现数据校验。 

- 在程序启动时扫描磁盘，一次性将全部对象缓存起来
- 接口服务实现了用户对象散列值和内容的校验，为了实现数据校验，我们将用户put的临时对象缓。存在数据服务节点。数据检验成功就转正，失败我就删除删除这个临时对象。
- 单例检查的方式实现对象去重，单份数据很危险，我们用数据冗余解决。

6.数据冗余和即时修复功能，使用 Reed-Solomon 纠错码，抵御数据丢失风险。 

- 4+2 RS 码的数据冗余策略。
- 查到任何数据发生丢失都会立即进行修复
- 6 个分片在不同的数据服务节点
- 使用第三方库 github.com/klauspost/reedsolomon

7.实现断点续传功能，并通过 Gzip 算法进行数据压缩。 

- 用 POST 上传大对象。
- 用 token 记录状态
- 客服端通过 Range 头部指定上传范围，first 必须跟 token 当前的字节数一致，不然返回 416 Range Not Satisfiable。如果上传的是最后一段数据，＜last>为空。
- 用 gzip 算法压缩对象，存储对象时压缩，取出时解压。节省磁盘资源

8.对象版本留存，定期对数据进行检验和修复。

- 留存最新的五个版本
- 清除过期对象的元数据、数据。
- 对当期对象数据的检查和修复等。 