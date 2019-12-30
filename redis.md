# redis 学习笔记
## 安装和启动
```shell
sudo apt install redis-server
```
```shell
redis-server
```
```shell
redis-cli
```

## 数据类型
* String: SET, GET
* Hash: HMSET, HGETALL
* List: lpush, lrange
* Set: sadd, smembers
* zset（有序集合）: zadd, zrangebyscore

## 数据库操作
* save 或者 bgsave： 数据库备份
* config get dir: 获得安装路径
* 将dump.rdb移动到安装路径并启动服务即可以恢复数据
* config set requirepass: 后面直接输入密码 设置密码
* config get requirepass: 获取密码
* redis-server --maxclients 100000: 设置服务最大连接数
* redis管道技术？？就是服务端未响应的时候客户端可以继续向服务端发送消息，然后之后一次性收到回复。反正就是变快5倍。
* 分区（partition）：优劣势
* 分区方式：
    1. 范围分区：用户id 1-10000分区1，10001-20000分区2 ...
    2. 哈希分区：先将名字hash得到一个数字，再用这个数字根据分区数量取模%


## 命令
```shell
# 本地客户端
redis-cli
# 远程客户端
redis-cli -h host -p port -a password
```
清晰的教程：
https://www.w3cschool.cn/redis/redis-hashes.html
```redis
# 各种键
get key
set key value
del key
exists key
expire key seconds
expireat key timestamp
pexpire key milliseconds
pexpireat key milliseconds-timestamp
keys pattern: 查找所有符合pattern的key
move key db：将key移动到指定db
persist key：移除key的过期时间
ttl key：查询一个key的剩余生存时间
randomkey：从当前数据库随机返回一个key
rename key newName
renamenx key newkey：当newkey不存在时，将key改名为newkey
type key：查询key所存储值的类型

# 各种string
getrange key start end
getset key value: 将给定key的值设为value，并返回key的旧值
getbit key offset: 对key所存储的字符串值，获取指定偏移量上的bit
mget key1 key2 key3...：获取给定key们的值
mset key1 value1 key2 value2 key3 value3...
msetnx key1 value1 key2 value2 key3 value3... : 同时设置多个key value，当且仅当所有给定key全都不存在
psetex key milliseconds value: 和setex命令相似，但是它以毫秒为单位设置key的生存时间，而不是像setex命令那样，以秒为单位
setbit key offset value: 对key所存储的字符串值，设置或者清除指定偏移量上的bit
setex key seconds value：将值value关联到key，并将key的过期时间设为以秒为单位
setnx key value: 只有在key不存在时设置key的值
setrange key offset value：用value参数覆盖key，从偏移量offset开始
strlen key: 返回key的len
incr key：将key中贮存的数字值+1
incrby key increment：将key中贮存的数字值+increment
incrbyfloat key increment：float type的增长
decr key
decrby key decrement
append key value：append值到key的末尾（前提key要存在）

# 各种哈希
HDEL key field1 field2 field3
HEXISTS key field: 查看key中指定字段是否存在
HGETALL key：获取这个key所有的字段
HINCRBY key field increment
HINCRBYFLOAT key field increment
HKEYS key：获取哈希表中所有的字段
HVALS key：获取哈希表中所有值
HLEN key：获取哈希表中字段的数量
HMGET key field1
HMSET key field1 value1 field2 value2
HSET key field value
HSETNX key field value
HSCAN key cursor [MATCH pattern][COUNT count]: 迭代哈希表中的键值对

# 各种列表
blpop key1 key2 timeout:pop列表第一个元素，如果列表没有元素会阻塞列表直到等待超市或发现可弹出元素为止
lpop key：pop列表第一个元素

brpop key1 key2 timeout：pop列表的最后一个元素
lindex key index
linsert key before｜after pivot value
llen key：获取列表长度
lpush key value1 value2
lpushx key value
lrange key start stop
lrem key count value：移除列表元素
lset key index value
ltrim key start stop：trim获得start - stop中间的元素
rpop key
rpoplpush source destination：移除列表的最后一个元素，并将该元素添加到另一个列表并返回
rpush key value1 value2
rpushx key value：为已经存在的列表添加值

# 各种集合
sadd set member1 member2...
scard set: 获取集合的成员数
sdiff set1 set2：返回所有集合的差集
sdiffstore destination set1 set2: 返回给定所有集合的差异并存储在destination中
sinter set1 set2：返回所有给定集合的交集
sinterstore destination set1 set2：求交集并存在destination中
sismember key member：判断member元素是否是集合key的成员
sismembers key：返回集合中所有成员
smove source destination member: 将元素从source集合mv到destination集合
spop key：移除并返回集合中的一个随机元素
srandmember key
srem key member1 member2
sunion set1 set2
sunionstore destination set1 set2
sscan key cursor [MATCH pattern][COUNT count]
```