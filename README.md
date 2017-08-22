编译
cd ../go-web-demo/src/person.mgtv.com/main/demo; go install

启动
cd ../go-web-demo/bin; ./demo

功能介绍
    1.项目结构
    go-web-demo
    --bin
      --conf   // 配置文件
      --template // 响应输出的模板
    --pkg
    --src 
      --person.mgtv.com
        --controller // 控制层
        --dao  
        --framework
          --config // 配置文件管理
          --database // 数据库包
          --httpclient // http包
          --logs // 日志包
          --mvc // 路由规则处理
          --redis //Redis包
          --resultcode // 错误码常量
        --main // main包 启动入口
        --model // 模块层
        --service // 业务逻辑层
        --thirdparty // 第三方业务依赖
      --vendor // 库依赖

1、路由
   路由比较简单、自己配置 

2、日志
    日志配置灵活，按天切割，1000ms或10000条日志，刷到磁盘.
    其中access.xml输出请求数据：时间|ip|URI|queryString|等...
    其中system.xml输出系统日志，如有异常自动分割到error.log日志

3、mysql
    支持事务提交
    
4、redis
    支持redisCluster
  
5、httpclient
    支持配置连接超时、读写超时、每个host的Idel数量、以及idel的时长...        

项目启动就会初始化数据库、redis、httpclient、日志等

简单测试了下性能：
1.从mysql获取一条记录 100000次 并发量500 qps=8313.81
http://192.168.9.29:8188/demo/get_feed?feedId=001b8b0b22294872996a7535493e49f2

2.从redis获取一个string 100000次 并发量500 qps=12726.70
http://192.168.9.29:8188/demo/get_key

3.事务提交
http://192.168.9.29:8188/demo/multi_commit

4.http判断是否关注
http://192.168.9.29:8188/demo/is_followed?uid=618225aa61d11662225896da1e398550&artistId=b1e4d1fee05caf22440fee8f4eb8c837
