# RocketMQ 简介与事务消息的性能测试
## RocketMQ 简介
作为背后默默支撑双十一的消息队列中间件，RocketMQ具有接近Kafka的消息快速处理速度。RocketMQ是基于Kafka开发的一个消息队列，对Kafka进行了一系列的优化。现在有很多项目中都利用到了RocketMQ消息队列。而且从4.3.0版本开始，开源的RocketMQ开始支持事务消息。鉴于网上并没有很多对于事务消息的压力测试，本文在这里进行了一次。

### 主要优点
* 严格保证消息顺序
* 上亿消息堆积能力
* 接近Kafka的速度表现
* 相对与Kafka，RocketMQ的一个亮点是支持事务消息

### 架构与主要组成部分
![avatar](http://rocketmq.apache.org/assets/images/rmq-basic-arc.png)

* Producer Cluster（消息生产者集群）：职责为将消息发送给Broker（消息代理）。有多种消息发送的模式，包括同步，异步，和单向。其中同步和异步为常用的发送模式。生产者可以有Producer Group（消息生产者组）将具有相同角色的消息生产者结合在一起。
* Consumer Cluster（消息消费者集群）：职责是从Broker（消息代理）中获取消息。消费者有两种，分别是
    1. Pull Consumer - 主动从Broker中拉取消息，一旦消息被拉出，用户应用便启动消息
    2. Push Consumer - 封装拉动消息，消费进度和其他维护工作。留下一个回调接口给最终用户实现，这个接口在消息到达时被执行。
   
   消费者也可以有Customer Group（消息消费者组）将具有相同角色的消息消费者组合在一起，相同角色可以理解为消费同一类消息
* NameServer Cluster：只提供路由信息
* Broker Cluster（消息代理集群）：主要组件，一般是以一个集群的方式存在。职责为接受producer发来的消息,存储这些消息等消费者来获取。Broker中分为Master broker和Slave broker两种。一个master可以有多个slave，但是一个salve只能有一个master。一个broker集群中可以有多个master。另外，Producer发送消息到Broker群时只能发送到Master，而Consumer接受消息时既可以从master获取消息，也可以从slave获取消息。
* 同步复制和异步复制
    * 同步复制
        1. CommitLog.putMessage 将消息写入内存缓存之中
        2. 调用 handleDiskFlush进行同步/异步刷盘
        3. handleHA 进行主从复制处理
    * 异步复制
        1. CommitLog.putMessage 将消息写入内存缓存之中
        2. 调用 handleDiskFlush进行同步/异步刷盘
        3. handleHA不会进行任何操作，也不管slave broker的复制进度，复制完全是由后台HAConnection.WriteSocketService服务在监听到有从Broker的链接可写时，向其写等待复制的数据。每个从broker发送进度则是由从broker定时汇报的自身当前已复制进度控制。该汇报由HAConnection.ReadSocketService负责处理

### RocketMQ中的一些消息概念
#### 普通消息
* Message（消息）: 一条消息包括Topic，Tag，消息内容等
* Topic（消息主题）: 消息主题，一般只有一到三四个
* Tag（消息标签）: 消息主题里的细分，一般推荐使用Tag来进行消息各类型的区分而不是Topic
* Message Queue（消息队列）：针对每个topic（消息主题），都可以有多个消息队列

#### 事务消息
除了包含普通信息的所有东西外，为了确保分布式事务的最终一致性和消息的原子性，RocketMQ的事务消息应用了两阶段提交理论（two-phase commit）在事务的前后分别发送两条消息。
![avatar](http://static.silence.work/mq_transcation.png)


## 性能测试
### 配置环境
性能测试在云服务器上完成，其配置为：
* CPU: 4 vCPU，每个vCPU双核2.8GHz
* RAM: 16GB
* OS: Ubuntu 18.04
* RocketMQ版本: 4.6.0
* 硬盘写入速度: 约 125 MB/s

### 测试方式
采用一个Master一个Slave的组合，分别采用 异步复制+异步刷盘，异步复制+同步刷盘，同步复制+异步刷盘，同步复制+同步刷盘 四种集群方式来进行ROCKETMQ的性能测试

#### 准备工作
1. 修改集群模式</br>
在"conf/2m-2s-async/"文件夹下有模板
```shell
# 所属于集群的名字
brokerClusterName=DefaultCluster
# 主从的brokerName保持一致来让他们找到自己的主或者从
brokerName=broker-a
# Master的ID为0，slave的ID为1到正无穷
brokerId=0
# 删除文件时间，默认4点
deleteWhen=04
# 文件保留时间，默认48小时
fileReservedTime=48
#主从设置和复制方式
#- ASYNC_MASTER 主从复制方式：异步复制，角色:Master
#- SYNC_MASTER 主从复制方式：同步复制， 角色：Master
#- SLAVE 角色：Slave (主从复制方式根据Master的设置来)
brokerRole=ASYNC_MASTER
# 刷盘方式
#- ASYNC_FLUSH 异步刷盘
#- SYNC_FLUSH 同步刷盘
flushDiskType=ASYNC_FLUSH
```

2. 启动nameServer
```shell
# nameServer的默认端口是9876
nohup sh bin/mqnamesrv &
```

3. 启动master broker
```shell
# -c 为配置文件路径（自定义）
# -n 为nameServer的地址和端口
nohup sh bin/mqbroker -c conf/2m-2s-async -n localhost:9876
```

4. 在另外一个云服务器上启动slave broker
```shell
nohup sh bin/mqbroker -c conf/2m-2s-async -n 10.128.0.3:9876
```

#### 执行测试
* 集群模式通过修改conf文件夹下
* 测试使用了RocketMQ中benchmark文件夹下自带的tproducer.sh文件
  ```shell
  sh benchmark/tproducer.sh --w=8 --s=2048
  ```
  这里的w参数为线程数（默认32），s参数为发送消息的大小（事务消息默认2048）

### 测试结果

**异步复制 异步刷盘**
* CPU利用率：100%，磁盘利用率70% ~ 80%

|并发数|消息大小|TPS|AverageRT|
|-----|-----|-----|-----|
|1|2048|2781|0.36ms|
|8|2048|4800|1.47ms|
|16|2048|5527|2.88ms|
|32|2048|6000|5.1ms|

**异步复制 同步刷盘**
* CPU利用率：6% ~ 11%，磁盘利用率100%

|并发数|消息大小|TPS|AverageRT|
|-----|-----|-----|-----|
|1|2048|387|2.58ms|
|8|2048|414|19.37ms|
|16|2048|445|34.1ms|
|32|2048|427|69.8ms|

**同步复制 异步刷盘**
* CPU利用率：50% ~ 66%，磁盘利用率80% ~ 90%

|并发数|消息大小|TPS|AverageRT|
|-----|-----|-----|-----|
|1|2048|1156|0.87ms|
|8|2048|2510|3.1ms| 
|16|2048|2465|5.7ms|
|32|2048|2675|11.3ms|

**同步复制 同步刷盘**
* CPU利用率：6% ~ 11&，磁盘利用率100%

|并发数|消息大小|TPS|AverageRT|
|-----|-----|-----|-----|
|1|2048|352|2.5ms|
|8|2048|401|18.55ms|
|16|2048|403|37.4ms|
|32|2048|410|75.4ms|


### 测试结论
* 在异步复制+异步刷盘模式，也就是最快的模式下，32线程可以达到6000左右的TPS
* 在具备更高可靠性的同步刷盘模式下，消息的速度大大的受到了硬盘写入速度的限制。TPS直接减少了近90%。在此次测试中的硬盘写入速度比较有限，大约为125MB/s，若能有读写更快的硬盘使用，TPS可以提升的空间还是很大的
* 当主从的复制模式从异步复制改为同步复制时，TPS大概减少了一半



### 常见的集群方式
* 单master：单点故障就瘫痪
* 多master：无单点故障，线上生产常用模式
* 单master多slave：slave瘫痪没事，master瘫痪后消费者还可以从slave获取消息但是生产者不能在发送消息
* 多master多slave：无单点故障


同步复制 -> 同步双写
