# 有关servlet的笔记
## servlet的生命周期
1. servlet创建的三种方式1 - servlet接口
```java
public class ServletDemo1 implements Servlet {

    //public ServletDemo1(){}

     //生命周期方法:当Servlet第一次被创建对象时执行该方法,该方法在整个生命周期中只执行一次
     //init：创建时候的函数
    public void init(ServletConfig arg0) throws ServletException {
                System.out.println("=======init=========");
        }

    //生命周期方法:对客户端响应的方法,该方法会被执行多次，每次请求该servlet都会执行该方法
    //service：service运行时候的函数
    public void service(ServletRequest arg0, ServletResponse arg1)
            throws ServletException, IOException {
        System.out.println("hehe");

    }

    //生命周期方法:当Servlet被销毁时执行该方法
    //destroy：service结束后的方法
    public void destroy() {
        System.out.println("******destroy**********");
    }

    //当停止tomcat时也就销毁的servlet。
    //getServeletConfig：tomcat也结束过后的方法
    public ServletConfig getServletConfig() {

        return null;
    }

    //getServletInfo：tomcat也结束过后的方法
    public String getServletInfo() {

        return null;
    }
}
```
2. servlet创建的方式2 - 继承GenericServlet类（极少使用）
   ```java
   public class ServletDemo2 extends GenericServlet {

    @Override
    public void service(ServletRequest arg0, ServletResponse arg1)
            throws ServletException, IOException {
        System.out.println("heihei");
        }
    }
   ```
3. servlet创建方式3 - 继承HttpServlet的方法（最常使用？）
```java
public class ServletDemo3 extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp)
            throws ServletException, IOException {
        System.out.println("haha");
    }

    @Override
    protected void doPost(HttpServletRequest req, HttpServletResponse resp)
            throws ServletException, IOException {
        System.out.println("ee");
        doGet(req,resp);
    }

}
```
* 更多关于HttpServlet
```java
abstract class HttpServlet extends GenericServlet{
 
   //HttpServlet中的service()
   protected void service(HttpServletRequest httpServletRequest,
                       HttpServletResponse httpServletResponse){
        //该方法通过httpServletRequest.getMethod()判断请求类型调用doGet() doPost()
   }
 
   //必须实现父类的service()方法
   public void service(ServletRequest servletRequest,ServletResponse servletResponse){
      HttpServletRequest request;
      HttpServletResponse response;
      try{
         request=(HttpServletRequest)servletRequest;
         response=(HttpServletResponse)servletResponse;
      }catch(ClassCastException){
         throw new ServletException("non-http request or response");
      }
      //调用service()方法
      this.service(request,response);
   }
}   
```