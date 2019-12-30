# kafka 笔记

## 推荐阅读
https://zhuanlan.zhihu.com/p/79579389

## 两种消费概念
1. 队列：消费者组允许同名的消费者组成员瓜分处理
2. 发布订阅：允许广播消息给多个消费者组（不同名）
kafka的每个topic都具有这两种模式

## 每个topic有多个分区，需要对多个消费者做负载均衡
* 有很强的顺序保证
* 保证了负载均衡
* 每个partition仅由同一个消费者组中的一个消费者消费到，并确保消费者是该partition的唯一消费者，并且按照顺序消费数据
* 每个partition只对应一个消费者组中的一个消费者，并确保消费者是该partition的唯一消费者
* tips：一个partition也可以对应两个消费者，前提是这两个消费者来自不同的消费者组
相同的消费者组中不能有比分区更多的消费者，否则多出的消费者一直处于空等待，不会收到消息
* 消费者commit之后同一个组的人就不会收到这些已经commit过的消息了
* 消费者只要给了group id之后消息在这个group中都是只能被消费一次了，似乎是默认会commit的？
* 不知道只是kafka-go这个包还是所有种类的client的缺点：用了customer group之后，以下几个函数不能使用了
    1. (*Reader).SetOffset会返回错误
    2. (*Reader).Offset会返回-1
    3. (*Reader).Lag会返回-1
    4. (*Reader).ReadLag会返回错误
    5. (*Reader).Stats会返回-1
* readMessage 在给了groupID之后，会默认commit
* fetchMessage不带默认的commit

* consumer 的 fetch request：如果小于min bytes，那就不会返回消息直到有足够的消息满足min bytes这个要求或者是超过了max wait设定的时间。调整大可以稍微加大服务器吞吐量
* max bytes超过2.1G就收不到消息了
* 如果需要消息持续过来的，就把minbytes设置得相对小一点，如果没有这方面要求只对吞吐量有要求的那就可以把min bytes设置得大一点稍微增加一点吞吐量
* 一次性写入的最大size为1MB，这个似乎不能改
* min和max都小于一块消息的大小时，每个fetch request会只取一部分直到把整条消息取完整

* 边读边写：设置写的比读的快：
* lastoffset 和 firstoffset 没区别？


## 具体工作方式
### Partition 和 replica
当某个topic的replication-factor为N且N大于1时，每个Partition都会有N个副本(Replica)。kafka的replica包含leader与follower。
Replica的个数小于等于Broker的个数，也就是说，对于每个Partition而言，每个Broker上最多只会有一个Replica，因此可以使用Broker id 指定Partition的Replica。
所有Partition的Replica默认情况会均匀分布到所有Broker上
### broker
1. Leader的选举




## sarama的客户端
* 发送消息：
	1. AsyncProducer: 
	2. SyncProducer: 直到收到收到消息的ack之后才能继续

### 全局参数？
* MaxRequestSize：request消息的最大size，默认的是100 * 1024 * 1024（100MB）
* MaxResponseSize：response消息的最大size， 默认的是100 * 1024 * 1024（100MB）
  
### AsyncProducer
* 可以AsyncClose 也可以Close 必须要被使用来防止数据流失

### client interface 可以查看的参数
* config
* controller：cluster controller broker
* Brokers：current set of active brokers
* Topics: 返回当前cluster里面所有的topic
* Partition：返回对于一个topic所有的partition ID

### clusterAdmin
* CreateTopic
* ListTopics
* DescribeTopics
* DeleteTopic
* CreatePartitions
* DeleteRecords
* DescribeConfig
* AlterConfig
* CreateACL: ACL（access control lists），对于某个resource建立ACL
* ListAcls
* DeleteACL
* ListConsumerGroups
* DescribeConsumeGroups
* ListConsumerGroupOffsets
* DeleteConsumerGroup
* DescribeCluster
* Close: shuts down the admin and closes underlying client

### config 主要参数
```go
KeepAlive: time.Duration, 0的话就不启用
LocalAddr: net.Addr, nil的话就自动选取
Metadata struct{
	Retry struct{
		Max: int, 在cluster处于leader election期间，重新尝试metadata request的最大次数
		Backoff: time.Duration, 两个retry之间空隙时间
		BackoffFunc： func(retries, maxRetries)，计算backOff time的方法，在更复杂的地方使用
	}
}

Producer struct{
	MaxMessageBytes: int, 消息的最大size（默认1MB），这里应该小于等于broker的message.max.bytes参数
	RequiredAcks: RequiredAcks, ack的模式，默认是WaitForLocal（其他选项：waitForAll）
	Timeout: time.Duration, 等待ack的最大时间
	Compression: CompressionCodec, 压缩方式（选项：CompressionNone, CompressionGZIP, CompressionSnappy, CompressionLZ4, CompressionZSTD）
	CompressionLevel: int, 压缩的力度
	Partitioner: PartitionerConstructor, 
	Idempotent: bool, 是否要确保一条消息只发了一次

	Flush struct{
	Bytes：flush需要的最少byte大小
	Messages：flush需要的最少message数量
	Frequency: time.Duration, flush的频率
	MaxMessages: int, 在一个producer向broker的request里面，最多能包含多少条消息
	}

	Retry struct{
		Max: int, 在cluster处于leader election期间，重新尝试metadata request的最大次数
		Backoff: time.Duration, 两个retry之间空隙时间
		BackoffFunc： func(retries, maxRetries)，计算backOff time的方法，在更复杂的地方使用
	}
}

Consumer struct{
	Group struct{
		session struct {
			Timeout: time.Duration, 心跳中timeout多久认定这个consumer失联，认定失联之后broker会去掉这个consumer然后进行rebalance，这个值必须在broker config里的"group.min.session.timeout.ms"和“group.max.session.timeout.ms”之间
		}
		Heartbeat struct {
			Interval time.Duration: 心跳的汇报间隔时间
		}
		Rebalance struct {
			Strategy: BalanceStrategy, * balance strategy（用于consumer group 分配 topic和partition）: 
                                        	1. RangeBalanceStrategyName：按照顺序给每个member分配顺序的n个partition，n为partition总量除以member的数量
	                                        2. RoundRobinBalanceStrategyName: 按照顺序轮流给每个member分配一个partition
	                                        3. StickyBalanceStrategyName：再分配的时候保证数量均衡分配的前提下尽量保持之前的分配结果（第一次分配的时候不清楚？似乎是RoundRobin）
			Timeout: time.Duration, 在rebalance开始之后，worker加入group的允许时间，超过了worker就会被从group中移除
			Retry struct {
				Max：int，rebalance的最大尝试次数
				Backoff：time.Duration，相邻两次尝试的间隔时间
			}	
		} 
		
	}

	Retry struct {
		Backoff，time.Duration：在失败读取一个partition之后，多少时间之再尝试
		BackoffFunc, func(retries int) time.Duration：计算backoff时间的方程
	}

	Fetch struct {
		Min：int32，fetch request去取的最小bytes（默认是1，0会导致consumer在没有消息可以领取的时候开始spin）
		Default：int32，每个request去取消息的大小（默认是1MB），这个参数应该比当前应用下多数message的大小大，不然consumer会花很多时间在size上
		Max：int32，每个request可以承载的最大bytes（默认是0，没有限制）。全局的那个“sarama.MaxResponseSize”依旧奏效
	}

	MaxWaitTime：time.Duration，就等待fetch request最小bytes的最长等待时间。（默认是250ms，0会使得consumer在没有event的时候开始spin）
	MaxProcessingTime：time.Duration，写入message channel的时间上限，如果超过这个时间，那个partition会停止拉取消息直到可以再次处理消息

	Return struct {
		Errors: bool
	}

	Offsets struct {
		AutoCommit struct {
			Enable: bool
			Interval: time.Duration, auto commit的频率
		}
		Initial：int64，值应该是OffsetNewest或者OffsetOldest
		Retention: time.Duration, commit过的offset的保存时间（默认是0。如果是0，则disable这个设置并使用broker的“offsets.retention.minutes”设置）

		Retry struct {
			Max：int，commit失败之后尝试的最大次数
		}
	}

	ClientID：string，（默认sarama）
	ChannelBufferSize：int，
	Version, KafkaVersion
	MetricRegistry, metrics.Registry
}
```


## kafka-go
1. Message
```go
{
    Partition int
    Offset    int64
	Key       []byte
	Value     []byte
    Headers   []Header
    Time      time.Time
}
```
2. write config
```go
//简约版
{
    Brokers []string
    Topic string
    Dialer *Dialer
    //指定一个分发消息至各个partition的balancer
    Balancer Balancer
    //发送时的最大尝试次数
    MaxAttempts int
    //writer 自身 queue的大小
    QueueCapacity int
    //一次request里可以存储的最大消息数量
    BatchSize int
    //request的最大size
    BatchBytes int
    //不完整batch发射的频率
    BatchTimeout time.Duration
    ReadTimeout time.Duration
    WriteTimeout time.Duration
    //刷新partition列表的频率
    RebalanceInterval time.Duration
    IdleConnTimeout time.Duration
    //在收到一个response之前需要收到几个ack
    RequiredAcks int
    //针对WriteMessage这个方法的异步
    Async bool
    //指定用于压缩消息的codec
    CompressionCodec
    Logger Logger
    ErrorLogger Logger

	newPartitionWriter func(partition int, config WriterConfig, stats *writerStats) partitionWriter 
}

//详细版本
{
    // The list of brokers used to discover the partitions available on the
	// kafka cluster.
	//
	// This field is required, attempting to create a writer with an empty list
	// of brokers will panic.
	Brokers []string

	// The topic that the writer will produce messages to.
	//
	// This field is required, attempting to create a writer with an empty topic
	// will panic.
	Topic string

	// The dialer used by the writer to establish connections to the kafka
	// cluster.
	//
	// If nil, the default dialer is used instead.
	Dialer *Dialer

	// The balancer used to distribute messages across partitions.
	//
	// The default is to use a round-robin distribution.
	Balancer Balancer

	// Limit on how many attempts will be made to deliver a message.
	//
	// The default is to try at most 10 times.
	MaxAttempts int

	// A hint on the capacity of the writer's internal message queue.
	//
	// The default is to use a queue capacity of 100 messages.
	QueueCapacity int

	// Limit on how many messages will be buffered before being sent to a
	// partition.
	//
	// The default is to use a target batch size of 100 messages.
	BatchSize int

	// Limit the maximum size of a request in bytes before being sent to
	// a partition.
	//
	// The default is to use a kafka default value of 1048576.
	BatchBytes int

	// Time limit on how often incomplete message batches will be flushed to
	// kafka.
	//
	// The default is to flush at least every second.
	BatchTimeout time.Duration

	// Timeout for read operations performed by the Writer.
	//
	// Defaults to 10 seconds.
	ReadTimeout time.Duration

	// Timeout for write operation performed by the Writer.
	//
	// Defaults to 10 seconds.
	WriteTimeout time.Duration

	// This interval defines how often the list of partitions is refreshed from
	// kafka. It allows the writer to automatically handle when new partitions
	// are added to a topic.
	//
	// The default is to refresh partitions every 15 seconds.
	RebalanceInterval time.Duration

	// Connections that were idle for this duration will not be reused.
	//
	// Defaults to 9 minutes.
	IdleConnTimeout time.Duration

	// Number of acknowledges from partition replicas required before receiving
	// a response to a produce request (default to -1, which means to wait for
	// all replicas).
	RequiredAcks int

	// Setting this flag to true causes the WriteMessages method to never block.
	// It also means that errors are ignored since the caller will not receive
	// the returned value. Use this only if you don't care about guarantees of
	// whether the messages were written to kafka.
	Async bool

	// CompressionCodec set the codec to be used to compress Kafka messages.
	// Note that messages are allowed to overwrite the compression codec individually.
	CompressionCodec

	// If not nil, specifies a logger used to report internal changes within the
	// writer.
	Logger Logger

	// ErrorLogger is the logger used to report errors. If nil, the writer falls
	// back to using Logger instead.
	ErrorLogger Logger

	newPartitionWriter func(partition int, config WriterConfig, stats *writerStats) partitionWriter 
}
```
3. Reader config
```go
//简约版
{
    Brokers []string
    GroupID string
    Topic string
    Partition int
    Dialer *Dialer
    QueueCapacity int
    //min and max bytes to fetch in each request
    MinBytes int
    MaxBytes int
    //等待新batch消息的MaxWait
    MaxWait time.Duration
    ReadLagInterval time.Duration
    //consumer group的优先列表
    GroupBalancers []GroupBalancer
    HeartbeatInterval time.Duration
    CommitInterval time.Duration
    //只有当有group id且watchPartitionChanges字段为true的时候才会被激活
    PartitionWatchInterval time.Duration
    WatchPartitionChanges bool
    // SessionTimeout optionally sets the length of time that may pass without a heartbeat
	// before the coordinator considers the consumer dead and initiates a rebalance.
    SessionTimeout time.Duration
    //等人的时间，等更多的人进来一起被balance
    RebalanceTimeout time.Duration
    //consumer group发生error之后等多久再次尝试连接
    JoinGroupBackoff time.Duration
    //多久执行save一次到broker
    RetentionTime time.Duration
    //从头还是从尾开始读取
    StartOffset int64

    ReadBackoffMin time.Duration
    ReadBackoffMax time.Duration
    Logger Logger
    ErrorLogger Logger
    IsolationLevel IsolationLevel
    MaxAttempts int
}
```
```go
{
    // The list of broker addresses used to connect to the kafka cluster.
	Brokers []string

	// GroupID holds the optional consumer group id.  If GroupID is specified, then
	// Partition should NOT be specified e.g. 0
	GroupID string

	// The topic to read messages from.
	Topic string

	// Partition to read messages from.  Either Partition or GroupID may
	// be assigned, but not both
	Partition int

	// An dialer used to open connections to the kafka server. This field is
	// optional, if nil, the default dialer is used instead.
	Dialer *Dialer

	// The capacity of the internal message queue, defaults to 100 if none is
	// set.
	QueueCapacity int

	// Min and max number of bytes to fetch from kafka in each request.
	MinBytes int
	MaxBytes int

	// Maximum amount of time to wait for new data to come when fetching batches
	// of messages from kafka.
	MaxWait time.Duration

	// ReadLagInterval sets the frequency at which the reader lag is updated.
	// Setting this field to a negative value disables lag reporting.
	ReadLagInterval time.Duration

	// GroupBalancers is the priority-ordered list of client-side consumer group
	// balancing strategies that will be offered to the coordinator.  The first
	// strategy that all group members support will be chosen by the leader.
	//
	// Default: [Range, RoundRobin]
	//
	// Only used when GroupID is set
	GroupBalancers []GroupBalancer

	// HeartbeatInterval sets the optional frequency at which the reader sends the consumer
	// group heartbeat update.
	//
	// Default: 3s
	//
	// Only used when GroupID is set
	HeartbeatInterval time.Duration

	// CommitInterval indicates the interval at which offsets are committed to
	// the broker.  If 0, commits will be handled synchronously.
	//
	// Default: 0
	//
	// Only used when GroupID is set
	CommitInterval time.Duration

	// PartitionWatchInterval indicates how often a reader checks for partition changes.
	// If a reader sees a partition change (such as a partition add) it will rebalance the group
	// picking up new partitions.
	//
	// Default: 5s
	//
	// Only used when GroupID is set and WatchPartitionChanges is set.
	PartitionWatchInterval time.Duration

	// WatchForPartitionChanges is used to inform kafka-go that a consumer group should be
	// polling the brokers and rebalancing if any partition changes happen to the topic.
	WatchPartitionChanges bool

	// SessionTimeout optionally sets the length of time that may pass without a heartbeat
	// before the coordinator considers the consumer dead and initiates a rebalance.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	SessionTimeout time.Duration

	// RebalanceTimeout optionally sets the length of time the coordinator will wait
	// for members to join as part of a rebalance.  For kafka servers under higher
	// load, it may be useful to set this value higher.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	RebalanceTimeout time.Duration

	// JoinGroupBackoff optionally sets the length of time to wait between re-joining
	// the consumer group after an error.
	//
	// Default: 5s
	JoinGroupBackoff time.Duration

	// RetentionTime optionally sets the length of time the consumer group will be saved
	// by the broker
	//
	// Default: 24h
	//
	// Only used when GroupID is set
	RetentionTime time.Duration

	// StartOffset determines from whence the consumer group should begin
	// consuming when it finds a partition without a committed offset.  If
	// non-zero, it must be set to one of FirstOffset or LastOffset.
	//
	// Default: FirstOffset
	//
	// Only used when GroupID is set
	StartOffset int64

	// BackoffDelayMin optionally sets the smallest amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 100ms
	ReadBackoffMin time.Duration

	// BackoffDelayMax optionally sets the maximum amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 1s
	ReadBackoffMax time.Duration

	// If not nil, specifies a logger used to report internal changes within the
	// reader.
	Logger Logger

	// ErrorLogger is the logger used to report errors. If nil, the reader falls
	// back to using Logger instead.
	ErrorLogger Logger

	// IsolationLevel controls the visibility of transactional records.
	// ReadUncommitted makes all records visible. With ReadCommitted only
	// non-transactional and committed records are visible.
	IsolationLevel IsolationLevel

	// Limit of how many attempts will be made before delivering the error.
	//
	// The default is to try 3 times.
	MaxAttempts int
}
```