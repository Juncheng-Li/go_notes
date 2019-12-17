# Spring boot学习笔记
## 准备工作
* 要记得三个组成部分
### properties 文件
```java
@Component
public class Book {

    @Value("${book.name}")
    private String name;
    @Value("${book.author}")
    private String author;

    // 省略getter和setter
}
```
**参数引用**
* 建议全部小写
```properties
book.name=SpringCloud
book.author=ZhaiYongchao
book.desc=${book.author}  is writing《${book.name}》
```
**随机数**
```properties
# 随机字符串
com.didispace.blog.value=${random.value}
# 随机int
com.didispace.blog.number=${random.int}
# 随机long
com.didispace.blog.bignumber=${random.long}
# 10以内的随机数
com.didispace.blog.test1=${random.int(10)}
# 10-20的随机数
com.didispace.blog.test2=${random.int[10,20]}
```
**List类型**
```properties
spring.my-example.url[0]=http://example.com
spring.my-example.url[1]=http://spring.io
# 或者
spring.my-example.url=http://example.com,http://spring.io
```

```YAML
spring:
  my-example:
    url:
      - http://example.com
      - http://spring.io

# 或者
spring:
  my-example:
    url: http://example.com, http://spring.io
```
**Map类型**
```properties
spring.my-example.foo=bar
spring.my-example.hello=world
```
```YAML
spring:
  my-example:
    foo: bar
    hello: world

#注意：如果Map类型的key包含非字母数字和-的字符，需要用[]括起来，比如：
spring:
  my-example:
    '[foo.baz]': bar
```

### 读取属性
* 方法一
```java
/** 
* 注意点
通过.分离各个元素
最后一个.将前缀与属性名称分开
必须是字母（a-z）和数字(0-9)
必须是小写字母
用连字符-来分隔单词
唯一允许的其他字符是[和]，用于List的索引
不能以数字开头
*/

//正确
this.environment.containsProperty("spring.jpa.database-platform")
//错误
this.environment.containsProperty("spring.jpa.databasePlatform")
```

* 方法二 - 读取：com.didispace.foo=bar
```java
//第一步
@Data
@ConfigurationProperties(prefix = "com.didispace")
public class FooProperties {

    private String foo;
    
}

//第二步
@SpringBootApplication
public class Application {

    public static void main(String[] args) {
        ApplicationContext context = SpringApplication.run(Application.class, args);

        Binder binder = Binder.get(context.getEnvironment());

        // 绑定简单配置
        FooProperties foo = binder.bind("com.didispace", Bindable.of(FooProperties.class)).get();
        System.out.println(foo.getFoo());
    }
}
```

* 读取List

读取这个List：
```properties
com.didispace.post[0]=Why Spring Boot
com.didispace.post[1]=Why Spring Cloud

com.didispace.posts[0].title=Why Spring Boot
com.didispace.posts[0].content=It is perfect!
com.didispace.posts[1].title=Why Spring Cloud
com.didispace.posts[1].content=It is perfect too!
```
```java
ApplicationContext context = SpringApplication.run(Application.class, args);

Binder binder = Binder.get(context.getEnvironment());

// 绑定List配置
List<String> post = binder.bind("com.didispace.post", Bindable.listOf(String.class)).get();
System.out.println(post);

List<PostInfo> posts = binder.bind("com.didispace.posts", Bindable.listOf(PostInfo.class)).get();
System.out.println(posts);
```

### 修改property的启动方式
```shell
java -jar xxx.jar --server.port=8888
```

### 配置文件的优先级
1. 命令行中传入的参数。
2. SPRING_APPLICATION_JSON中的属性。SPRING_APPLICATION_JSON是以JSON格式配置在系统环境变量中的内容。
3. java:comp/env中的JNDI属性。
4. Java的系统属性，可以通过System.getProperties()获得的内容。
5. 操作系统的环境变量
6. 通过random.*配置的随机属性
7. 位于当前应用jar包之外，针对不同{profile}环境的配置文件内容，例如：application-{profile}.properties或是YAML定义的配置文件
8. 位于当前应用jar包之内，针对不同{profile}环境的配置文件内容，例如：application-{profile}.properties或是YAML定义的配置文件
9. 位于当前应用jar包之外的application.properties和YAML配置内容
10. 位于当前应用jar包之内的application.properties和YAML配置内容
11. 在@Configuration注解修改的类中，通过@PropertySource注解定义的属性
12. 应用默认属性，使用SpringApplication.setDefaultProperties定义的内容

### 推荐的模版引擎
* Thymeleaf
* FreeMarker
* Velocity
* Groovy
* Mustache
* 不要使用JSP，支持JSP需要修改配置

### Thymeleaf模版
```properties
# Enable template caching.
spring.thymeleaf.cache=true 
# Check that the templates location exists.
spring.thymeleaf.check-template-location=true 
# Content-Type value.
spring.thymeleaf.content-type=text/html 
# Enable MVC Thymeleaf view resolution.
spring.thymeleaf.enabled=true 
# Template encoding.
spring.thymeleaf.encoding=UTF-8 
# Comma-separated list of view names that should be excluded from resolution.
spring.thymeleaf.excluded-view-names= 
# Template mode to be applied to templates. See also StandardTemplateModeHandlers.
spring.thymeleaf.mode=HTML5 
# Prefix that gets prepended to view names when building a URL.
spring.thymeleaf.prefix=classpath:/templates/ 
# Suffix that gets appended to view names when building a URL.
spring.thymeleaf.suffix=.html  spring.thymeleaf.template-resolver-order= # Order of the template resolver in the chain. spring.thymeleaf.view-names= # Comma-separated list of view names that can be resolved.
```

## 构建Restful API
* @Controller 修饰class，用来创建处理http请求的对象
* @RestController Spring4之后加入的注解，相较于@Controller不需要再配置@ResponseBody，默认返回json格式
* @RequestMapping 默认url映射

## Swagger2
* 先建立一个Swagger2.java - 这里面就是主页的简介
```java
@ApiOperation(value="更新用户详细信息", notes="根据url的id来指定更新对象，并根据传过来的user信息来更新用户详细信息")
@ApiImplicitParams({
    @ApiImplicitParam(name = "id", value = "用户ID", required = true, dataType = "Long"),
    @ApiImplicitParam(name = "user", value = "用户详细实体user", required = true, dataType = "User")
})
```

* http://localhost:8080/swagger-ui.html
* 