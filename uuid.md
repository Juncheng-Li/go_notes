# UUID
## UUID 的格式：
```
123e4567-e89b-12d3-a456-426655440000
xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
```
* 一共36个字符
* 四位数字 M表示 UUID 版本，数字 N的一至三个最高有效位表示 UUID 变体。在例子中，M 是 1 而且 N 是 a(10xx)，这意味着此 UUID 是 "变体1"、"版本1" UUID；即基于时间的 DCE/RFC 4122 UUID。

**UUID record layout**

|Name|Length (bytes)|Length (hex digits)|Contents|
|---|---|---|---|
|time_low|4|8|整数：低位 32 bits 时间|
|time_mid|2|4|整数：中间位 16 bits 时间|
|time_hi_and_version|2|4|最高有效位中的 4 bits“版本”，后面是高 12 bits 的时间|
|clock_seq_hi_and_res clock_seq_low|2|4|最高有效位为 1-3 bits“变体”，后跟13-15 bits 时钟序列|
|node|6|12|48 bits 节点 ID|

## UUID的版本：
* 版本1：UUID 是根据时间和节点ID（通常是生成UUID的计算机MAC地址）生成
* 版本2（DCE安全版本）：UUID 是根据标识符（通常是组或者用户ID）、时间和节点ID生成
* 版本3（MD5）和版本5（SHA1）：确定性UUID通过散列（hashing）名字空间（namespace）标识符和名称生成
* 版本4：UUID使用随机性或伪随机性生成 - NewRandom() (UUID, error)

## 碰撞可能
* 版本1，版本2：当多次生成相同的 UUID 并将其分配给不同的指示对象时，就会发生冲突。对于使用来自网卡的唯一MAC地址的标准 版本1和2 的 UUID ，只有当实施与标准不同时，无论是无意还是故意，都可能发生冲突
* 版本3，版本4，版本5：碰撞概率可以由生日问题获得 比如：P(有两个人生日一样) = (1 - (365!))/(365^n * (365 - n)!) n为：n个人中间两个人生日相同的概率。根据wiki上的计算 在 103 万亿 个 版本4  UUID 中找到重复的概率是十亿分之一 

### 碰撞可能计算

### 碰撞可能模拟
由于通过生日悖论的方法进行计算进行了碰撞可能的模拟

```python
```
## 性能测试

## 其他
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
