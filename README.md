编译<br>
cd ../go-web-demo/src/person.mgtv.com/main/demo; go install<br>

启动<br>
cd ../go-web-demo/bin; ./demo<br>

功能介绍<br>
1.项目结构<br>
　　　go-web-demo<br>
　　　--bin<br>
　　　--conf   // 配置文件<br>
　　　--template // 响应输出的模板<br>
　　　--pkg<br>
　　　--src <br>
　　　　--person.mgtv.com<br>
　　　　　--controller // 控制层<br>
　　　　　--dao  <br>
　　　　　--framework<br>
　　　　　--config // 配置文件管理<br>
　　　　　--database // 数据库包<br>
　　　　　--httpclient // http包<br>
　　　　　--logs // 日志包<br>
　　　　　--mvc // 路由规则处理<br>
　　　　　--redis //Redis包<br>
　　　　　--resultcode // 错误码常量<br>
　　　　　--main // main包 启动入口<br>
　　　　　--model // 模块层<br>
　　　　　--service // 业务逻辑层<br>
　　　　　--thirdparty // 第三方业务依赖<br>
　　　--vendor // 库依赖<br>

1、路由<br>
   路由比较简单、自己配置 <br>

2、日志<br>
    日志配置灵活，按天切割，1000ms或10000条日志，刷到磁盘.<br>
    其中access.xml输出请求数据：时间|ip|URI|queryString|等...<br>
    其中system.xml输出系统日志，如有异常自动分割到error.log日志<br>

3、mysql<br>
    支持事务提交<br>
    
4、redis<br>
    支持redisCluster<br>
  
5、httpclient<br>
    支持配置连接超时、读写超时、每个host的Idel数量、以及idel的时长...    <br>    

项目启动就会初始化数据库、redis、httpclient、日志等<br>

简单测试了下性能：<br>
1.从mysql获取一条记录 100000次 并发量500 qps=8313.81<br>
http://192.168.9.29:8188/demo/get_feed?feedId=001b8b0b22294872996a7535493e49f2<br>

2.从redis获取一个string 100000次 并发量500 qps=12726.70<br>
http://192.168.9.29:8188/demo/get_key<br>

3.事务提交<br>
http://192.168.9.29:8188/demo/multi_commit<br>

4.http判断是否关注<br>
http://192.168.9.29:8188/demo/is_followed?uid=618225aa61d11662225896da1e398550&artistId=b1e4d1fee05caf22440fee8f4eb8c837<br>
