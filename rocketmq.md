# Rocket mq 性能测试
## 同步复制和异步复制
**同步复制**
1. CommitLog.putMessage 将消息写入内存缓存之中
2. 调用 handleDiskFlush进行同步/异步刷盘
3. handleHA 进行主从复制处理

**异步复制**
1. CommitLog.putMessage 将消息写入内存缓存之中
2. 调用 handleDiskFlush进行同步/异步刷盘
3. handleHA不会进行任何操作，也不管slave broker的复制进度，复制完全是由后台HAConnection.WriteSocketService服务在监听到有从Broker的链接可写时，向其写等待复制的数据。每个从broker发送进度则是由从broker定时汇报的自身当前已复制进度控制。该汇报由HAConnection.ReadSocketService负责处理

* 个人觉得 异步复制 + 同步刷盘 比较好一点。同步刷盘已经保证了可靠性


## 性能测试
### 配置环境
性能测试在云服务器上完成，其配置为：
* CPU：4 vCPU，双核，2.8GHz
* RAM：16GB
* OS：Ubuntu 18.04
* RocketMQ: 4.6.0
<br>

### 实验方式
* 采用一个Master一个Slave的组合，分别采用 异步复制+异步刷盘，异步复制+同步刷盘，同步复制+异步刷盘，同步复制+同步刷盘 四种集群方式来进行ROCKETMQ的性能测试

**异步复制 异步刷盘**

|集群模式|并发数|消息大小|TPS|AverageRT|
|----|-----|-----|-----|-----|
|异步复制 异步刷盘|1|2048|2781|0.36s|
|异步复制 异步刷盘|8|2048|4800|1.47s|
|异步复制 异步刷盘|16|2048|5527|2.88s|
|异步复制 异步刷盘|32|2048|6000|5.1s|
<br>

**异步复制 同步刷盘**

|集群模式|并发数|消息大小|TPS|AverageRT|
|----|-----|-----|-----|-----|
|异步复制 同步刷盘|1|2048|387|2.58s|
|异步复制 同步刷盘|8|2048|414|19.37s|
|异步复制 同步刷盘|16|2048|445|34.1s|
|异步复制 同步刷盘|32|2048|427|69.8s|
<br>

**同步复制 异步刷盘**
55% cpu
|集群模式|并发数|消息大小|TPS|AverageRT|
|----|-----|-----|-----|-----|
|同步复制 异步刷盘|1|2048|1156|0.87s|
|同步复制 异步刷盘|8|2048|2510|3.1s| 
|同步复制 异步刷盘|16|2048|2465|5.7s|
|同步复制 异步刷盘|32|2048|2675|11.3s|
<br>

**同步复制 同步刷盘**

|集群模式|并发数|消息大小|TPS|AverageRT|
|----|-----|-----|-----|-----|
|同步复制 同步刷盘|1|2048|352|2.5s|
|同步复制 同步刷盘|8|2048|401|18.55s|
|同步复制 同步刷盘|16|2048|403|37.4s|
|同步复制 同步刷盘|32|2048|410|75.4s|
<br>

* 当刷盘模式改为同步刷盘时候，TPS减少到了400左右。此时的磁盘利用率为100%，然而CPU利用率只有6% ~ 10%。并发数的增加也对TPS毫无增益。由此可见，在涉及同步刷盘的集群模式中，其TPS的瓶颈在于磁盘的写入速度
* 当主从的复制模式从异步复制改为同步复制时，TPS收到的影响很小





### 个性化配置文件
```
#所属集群名字
brokerClusterName=rocketmq-cluster
#broker名字，注意此处不同的配置文件填写的不一样  例如：在a.properties 文件中写 broker-a  在b.properties 文件中写 broker-b
brokerName=broker-a
#0 表示 Master，>0 表示 Slave
brokerId=0
#nameServer地址，这里nameserver是单台，如果nameserver是多台集群的话，就用分号分割（即namesrvAddr=ip1:port1;ip2:port2;ip3:port3或者修改hosts文件不使用IP而是域名）
namesrvAddr=rocketmq-nameserver1:9876;rocketmq-nameserver2:9876
#在发送消息时，自动创建服务器不存在的topic，默认创建的队列数。由于是4个broker节点，所以设置为4
defaultTopicQueueNums=4
#是否允许 Broker 自动创建Topic，建议线下开启，线上关闭
autoCreateTopicEnable=true
#是否允许 Broker 自动创建订阅组，建议线下开启，线上关闭
autoCreateSubscriptionGroup=true
#Broker 对外服务的监听端口
listenPort=10911
#删除文件时间点，默认凌晨 4点
deleteWhen=04
#文件保留时间，默认 48 小时
fileReservedTime=120
#commitLog每个文件的大小默认1G
mapedFileSizeCommitLog=1073741824
#ConsumeQueue每个文件默认存30W条，根据业务情况调整
mapedFileSizeConsumeQueue=300000
#destroyMapedFileIntervalForcibly=120000
#redeleteHangedFileInterval=120000
#检测物理文件磁盘空间
diskMaxUsedSpaceRatio=88
#存储路径
storePathRootDir=/usr/local/rocketmq/store
#commitLog 存储路径
storePathCommitLog=/usr/local/rocketmq/store/commitlog
#消费队列存储路径存储路径
storePathConsumeQueue=/usr/local/rocketmq/store/consumequeue
#消息索引存储路径
storePathIndex=/usr/local/rocketmq/store/index
#checkpoint 文件存储路径
storeCheckpoint=/usr/local/rocketmq/store/checkpoint
#abort 文件存储路径
abortFile=/usr/local/rocketmq/store/abort
#限制的消息大小
maxMessageSize=65536
#flushCommitLogLeastPages=4
#flushConsumeQueueLeastPages=2
#flushCommitLogThoroughInterval=10000
#flushConsumeQueueThoroughInterval=60000
#Broker 的角色
#- ASYNC_MASTER 异步复制Master
#- SYNC_MASTER 同步双写Master
#- SLAVE
brokerRole=ASYNC_MASTER
#刷盘方式
#- ASYNC_FLUSH 异步刷盘
#- SYNC_FLUSH 同步刷盘
flushDiskType=ASYNC_FLUSH
#checkTransactionMessageEnable=false
#发消息线程池数量
#sendMessageThreadPoolNums=128
#拉消息线程池数量
#pullMessageThreadPoolNums=128
```