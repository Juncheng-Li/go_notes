# HTTP 知识点
## 请求报文
请求行
get, push, put, delete（请求方法）		url		协议版本
请求头
服务器需要知道的附加信息
（空一行）

请求体（就是参数）

请求方法
Get	请求资源
Post	提交资源
Head	获取相应头
Put	替换资源
Delete	删除资源
Options	允许客户端查看服务器性能
Trace	回显服务器收到的请求，用于测试或诊断

常见的请求头
Host	主机ip地址或域名
User-Agent		客户端相关信息，如操作系统，浏览器等信息
Accept			指定客户端接受信息类型，如：image/jpg，text/html，applicaion
Accept-Charset	客户端接受的字符集，如gb2312、iso-8859-1
accept-encoding	可接受的内容编码，如gzip
accept-language	接受的语言，如accept-language：zh-cn
Authorization		客户端给服务器授权认证的信息
Cookie			携带的cookie信息
Referer		当前文档的url，即从哪个文件过来的	- 盗链，测流量，竞价排名等等
content-type		请求体内容类型，如content-type：application/x-www-form-url
content-length	数据长度
cache-control		缓存机制，如cache-control：no-cache
Pragma		防止页面被缓存，和cache-conrol：no-cache作用一样

## 响应报文
响应行
http版本	状态码		状态描述（ok）
响应头
比较重要的有：set-cookie	content-type：响应的类型和字符集 如 content-type：text		location：指明重新定向的地址
（空行）
响应体


参数传值法：get
表单传值法：post

### get和post的区别
* 我们可以看到GET请求消息体为空，POST请求带有消息体。
* GET提交的数据会放在URL之后，以?分割URL和传输数据，参数之间以&相连，如EditPosts.aspx?name=test1&id=123456。POST方法是把提交的数据放在HTTP包的body中。
* GET提交的数据大小有限制（因为浏览器对URL的长度有限制），而POST方法提交的数据没有限制。
* GET方式提交数据，会带来安全问题，比如一个登录页面，通过GET方式提交数据时，用户名和密码将出现在URL上，如果页面可以被缓存或者其他人可以访问这台机器，就可以从历史记录获得该用户的账号和密码