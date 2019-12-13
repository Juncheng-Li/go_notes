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

|Name|Length (bytes)|Length (hex digits)|Contents|
|---|---|---|---|
|time_low|4|8|整数：低位 32 bits 时间|
|time_mid|2|4|整数：中间位 16 bits 时间|
|time_hi_and_version|2|4|最高有效位中的 4 bits“版本”，后面是高 12 bits 的时间|
|clock_seq_hi_and_res clock_seq_low|2|4|最高有效位为 1-3 bits“变体”，后跟13-15 bits 时钟序列|
|node|6|12|48 bits 节点 ID|

### UUID的版本：
* 版本1：UUID 是根据时间和节点ID（通常是生成UUID的计算机MAC地址）生成
* 版本2（DCE安全版本）：UUID 是根据标识符（通常是组或者用户ID）、时间和节点ID生成
* 版本3（MD5）和版本5（SHA1）：确定性UUID通过散列（hashing）名字空间（namespace）标识符和名称生成
* 版本4：UUID使用随机性或伪随机性生成 - NewRandom() (UUID, error)

### google/UUID的版本：
google/UUID包中UUID的版本和网上普遍描述的有一些轻微的不同
* 版本1：根据

### 碰撞可能：
* 版本1，版本2：当多次生成相同的 UUID 并将其分配给不同的指示对象时，就会发生冲突。对于使用来自网卡的唯一MAC地址的标准 版本1和2 的 UUID ，只有当实施与标准不同时，无论是无意还是故意，都可能发生冲突
* 版本3，版本4，版本5：碰撞概率可以由生日问题获得，计算方式在下面

#### 碰撞概率计算
生日悖论描述的是在n个学生中间，有两个或两个以上学生生日一样的概率是多少。<br>
生日悖论中，相同生日碰撞概率的计算方式如下（n为学生的数量）：

$$P(有至少两个人生日一样) = \frac{1 - 365!}{365^n * (365 - n)!}$$

在计算UUID碰撞概率的情景中，365应该改为UUID的总量，也就是：
$$ (26 + 10)^{32} \approx 6 \times 10^{40}$$
这里的26是指 a-z 26个字母，UUID只生成小写字母所以不用考虑大写<br>
这里的10是指 0-9 10个数字<br>
所以UUID中每一位都有26 + 10 = 36种可能<br>
UUID总共32位，所以最后结果通过上面的公式计算<br>
<br>
n在这里的情景下代表生成n个uuid<br>
所以我们的公式如下：

$$P(UUID 碰撞) = \frac{1 - 6e40!}{6e40^n * (6e40 - n)!}$$
根据wiki上的计算， 在103万亿个版本4的UUID中找到重复UUID的概率是十亿分之一 

#### 碰撞概率模拟
由于涉及到阶乘，通过生日悖论的方法进行计算会计算出非常大的数字从而导致溢出错误<br>
所以这里写了一段python的代码来对碰撞直接进行模拟

```python
import random, time
from tqdm import tqdm

# Prepare
# UUID的总量
pool_size = 6e40
# 模拟生成n个UUID，test_num就是n
test_num = 1e8
# 一共要模拟多少轮
# 通常 >= 800 轮来获得一个比较准确的数字
rounds = 1
# 初始化计数，记下所有轮中有多少轮是发生了UUID的碰撞
count = 500

# Running
start = time.perf_counter()
print("Simulation started...")
for i in tqdm(range(rounds)):
    # 随机抽取从1到6e40抽取n次，组成list
    bds = [random.randint(1, pool_size) for x in range(int(test_num))]
    # 利用集合的无重复性来判断是否有人生日一致
    if len(bds) != len(set(bds)):
        count += 1
print("Simulation finished...")
end = time.perf_counter()

# Finish
timecost = round(end - start, 2)
print("Rounds: ", rounds)
print("count: ", count)
print("Rate: ", count/rounds)
print("time: ", timecost)
```

由于每次测试用时几到几十小时，总共只进行了3次模拟
1. 32位uuid总量大约6e40，生成千万的uuid，总共模拟1000次，碰撞概率都为0？
2. 32位uuid总量大约6e40，生成一亿的uuid，总共模拟500次，碰撞概率？
3. 32位uuid总量大约6e40，生成十亿的uuid，总共模拟200次，有没有发生碰撞？
<!-- 4. 32位uuid总量大约6e40，生成百亿的uuid，总共模拟100次，有没有发生碰撞？
1. 32位uuid总量大约6e40，生成千亿的uuid，总共模拟100次，有没有发生碰撞？ -->

### 性能测试：
#### 计算机配置
* CPU：3.4GHz，4核8线程

#### 测试方式
* 利用100个协程同时生成

#### 测试结果

|UUID版本|生成速度|
|-----|-----|
|版本1|6934835 UUID/s|
|版本2（通过Group ID生成）|6227007 UUID/s|
|版本2（通过Person ID生成）|6217744 UUID/s|
|版本2（通过自定义的命名空间+自定义的ID生成）|6281429 UUID/s|
|版本3|17116305 UUID/s|
|版本4|5279727 UUID/s|
|版本5|12414673 UUID/s|

### 其他
```
var (
    NameSpaceDNS  = Must(Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
    NameSpaceURL  = Must(Parse("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
    NameSpaceOID  = Must(Parse("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
    NameSpaceX500 = Must(Parse("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
    Nil           UUID // empty UUID, all zeros
)
```
Well known namespace IDs and UUIDs
