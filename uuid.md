# Go语言下的UUID性能测试
## 简介
本测试基于google的UUID包，传送门：https://github.com/google/uuid

### UUID 的格式：
```
123e4567-e89b-12d3-a456-426655440000
xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
```
* 一共32位
* 四位数字 M表示 UUID 版本，数字 N的一至三个最高有效位表示 UUID 变体。在例子中，M 是 1 而且 N 是 a(10xx)，这意味着此 UUID 是 "变体1"、"版本1" UUID；即基于时间的 DCE/RFC 4122 UUID。

#### 例: 版本1和版本2的UUID各部分组成
|Name|Length (bytes)|Length (hex digits)|Contents|
|---|---|---|---|
|time_low|4|8|整数：低位 32 bits 时间|
|time_mid|2|4|整数：中间位 16 bits 时间|
|time_hi_and_version|2|4|最高有效位中的 4 bits“版本”，后面是高 12 bits 的时间|
|clock_seq_hi_and_res clock_seq_low|2|4|最高有效位为 1-3 bits“变体”，后跟13-15 bits 时钟序列|
|node|6|12|48 bits 节点 ID|

### 普遍UUID的版本：
* 版本1：UUID 是根据时间和节点ID（通常是生成UUID的计算机MAC地址）生成
* 版本2（DCE安全版本）：UUID 是根据标识符（通常是组或者用户ID）、时间和节点ID生成
* 版本3（MD5）和版本5（SHA1）：确定性UUID通过散列（hashing）名字空间（namespace）标识符和名称生成
* 版本4：UUID使用随机性或伪随机性生成 - NewRandom() (UUID, error)

### google/UUID的版本：
google/UUID包中UUID的版本和网上普遍描述的有一些轻微的不同
* **版本1:** NewUUID()
    * 根据clock sequence(时针序列)，节点ID（mac地址）生成。（注：这里的clock sequence可以被指定成想要的数字）

* **版本2:** google/UUID这个包中的版本2有三种不同的方式
    1. NewDCESecurity(Group, uint32(os.Getgid())）
        * 通过始终序列,节点ID,Group和一个group id生成，其中Group和Group id需要被输入。这里的Group是一个domain，可以通过uuid.Domain(n)生成获得（n为任意数字）。
        * 也可以通过NewDCEGroup()生成
        * 是DCE（Distributed Computing Environment）安全的版本。
    2. NewDCESecurity(Person, uint32(os.Getuid()))
        * 与第一种方法类似。通过Person和一个uid生成，这里的Person和上面方式中的Group一样是一个domain，可以通过uuid.Domain(n)函数生成获得。
        * 也可以通过NewDCEPerson()生成
        * 是DCE（Distributed Computing Environment）安全的版本。
    3. NewDCESecurity(domain Domain, id uint32)
        * 与上面的方式类似，只是将第一个参数从指定的改成了任意的domain，第二个参数为在这个domain下的id
        * 是DCE（Distributed Computing Environment）安全的版本。
  
* **版本3:** NewMD5(space UUID, data []byte)
  * 第一个参数是一个命名空间
  * 第二个参数是任意的数据
  * 也可以通过NewHash(md5.New(), space, data, 3)生成，参数相同

* **版本4:** NewRandom()
    * 版本4的UUID直接通过随机生成，无需任何参数
    * 碰撞概率是十亿分之一

* **版本5:** NewSHA1(space UUID, data []byte)
    * 版本5与版本3类似，同样的输入参数，只是由哈希方式换成了SHA1加密方式生成
    * 也可以通过NewHash(sha1.New(), space, data, 5)方式生成

### 碰撞可能：
* 版本1，版本2：当多次生成相同的 UUID 并将其分配给不同的指示对象时，就会发生冲突。对于使用来自网卡的唯一MAC地址的标准 版本1和2 的 UUID ，只有当实施与标准不同时，无论是无意还是故意，都可能发生冲突
* 版本3，版本4，版本5：碰撞概率可以由生日问题获得，计算方式在下面

#### 碰撞概率计算
生日悖论描述的是在n个学生中间，有两个或两个以上学生生日一样的概率是多少。<br>
计算方式如下（n为学生的数量）：

$$P(有至少两个人生日一样) = \frac{1 - 365!}{365^n * (365 - n)!}$$

在计算UUID碰撞概率的情景中，365应该改为UUID的总量，也就是：
$$ (26 + 10)^{32} \approx 6 \times 10^{40}$$
(26+10)中的26是指 a-z 26个字母，UUID只生成小写字母所以不用考虑大写<br>
(26+10)中的10是指 0-9 10个数字<br>
所以UUID中每一位都有26 + 10 = 36种可能<br>
UUID总共32位，所以最后结果通过上面的公式计算得到UUID的总量大概为6e40<br>
n在计算UUID的情景下代表生成n个uuid<br>
所以最后的公式如下：
$$P(UUID 碰撞) = \frac{1 - 6e40!}{6e40^n * (6e40 - n)!}$$
不过由于此应用下数字极大且涉及阶乘，通过这种方式直接在电脑上计算会导致非常大的数字从而导致溢出错误<br>
根据wiki上的计算， 在103万亿个版本4的UUID中找到重复UUID的概率是十亿分之一 

#### 碰撞概率模拟
由于涉及到阶乘，通过生日悖论的方法进行计算会计算出非常大的数字从而导致溢出错误<br>
所以这里写了一段python代码来对碰撞直接进行模拟

```python
import random, time
from tqdm import tqdm

# 准备工作, 参数设置:
# UUID的总量
pool_size = 6e40
# 公式中的n, 模拟生成UUID的数量
test_num = 1e8
# 重复模拟的次数
rounds = 900
# 初始化计数
count = 0

# 开始模拟:
start = time.perf_counter()
print("Simulation started...")
for i in tqdm(range(rounds)):
    # 列表推导式：随机抽取从(1,366)抽取23次，组成列表
    bds = [random.randint(1, pool_size) for x in range(int(test_num))]
    # 利用集合的无重复性来判断是否有人生日一致
    if len(bds) != len(set(bds)):
        count += 1
print("Simulation finished...")
end = time.perf_counter()

# 结束模拟, 打印结果:
timecost = round(end - start, 2)
print("Rounds: ", rounds)
print("count: ", count)
print("Rate: ", count/rounds)
print("time: ", timecost)
```

由于每次测试用时几到几十小时和内存上的限制，总共只进行了2次模拟
1. 生成千万的uuid，总共模拟1000次，没有发生碰撞
2. 生成一亿的uuid，总共模拟1000次，没有发生碰撞


### 性能测试：
#### 计算机配置
* CPU：3.4GHz，4核8线程
* RAM: 16G

#### 测试方式
1. 利用单协程生成
2. 利用100个协程同时生成

#### 测试结果
1. 单协程下的测试结果:
    |UUID版本|生成速度|
    |-----|-----|
    |版本1|14021408 UUID/s|
    |版本2（通过Group ID生成）|12491274 UUID/s|
    |版本2（通过Person ID生成）|12113973 UUID/s|
    |版本2（通过自定义的命名空间+自定义的ID生成）|11959340 UUID/s|
    |版本3|4412772 UUID/s|
    |版本4|3938385 UUID/s|
    |版本5|3684399 UUID/s|


2. 100个多协程下的测试结果:
    |UUID版本|生成速度|
    |-----|-----|
    |版本1|6934835 UUID/s|
    |版本2（通过Group ID生成）|6227007 UUID/s|
    |版本2（通过Person ID生成）|6217744 UUID/s|
    |版本2（通过自定义的命名空间+自定义的ID生成）|6281429 UUID/s|
    |版本3|17116305 UUID/s|
    |版本4|5279727 UUID/s|
    |版本5|12414673 UUID/s|


#### 测试结论
* 单协程下版本1与版本2有最快的生成方式, 是其余方式的三到四倍
* 多协程下很明显版本3有最快的生成速度，版本5次之。分别是其余版本的二到三倍。
* 版本1和版本2在单协程下速度更快, 是多协程下的两倍, 其最主要的原因是他们生成的函数中涉及了互斥锁
* 版本3, 4, 5在多协程下的速度更快一些, 尤其是版本3和版本5, 速度提升是原来的四倍
