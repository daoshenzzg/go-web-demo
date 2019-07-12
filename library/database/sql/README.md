#### database/sql

##### 项目简介
MySQL数据库驱动，进行封装加入了链路追踪和统计。

如果需要SQL级别的超时管理 可以在业务代码里面使用context.WithDeadline实现 推荐超时配置放到application.toml里面 方便热加载

##### 依赖包
1. [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql)

```$xslt
[mysql]
    addr = "127.0.0.1:3306"
    dsn = "root:123456@tcp(127.0.0.1:3306)/bilibili_answer?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8,utf8mb4"
    active = 5
    idle = 2
    idleTimeout ="4h"
    queryTimeout = "1000ms"
    execTimeout = "1000ms"
    tranTimeout = "2000ms"
```