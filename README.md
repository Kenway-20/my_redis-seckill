# MyRedis and Distributed Seckill System

- Technology Stack:
  - Go, MySQL, Redis, RabbitMQ
- Key Modules:
  - MyRedis Middleware: Independently developed middleware with Redis functionality, called MyRedis
  - Distributed Authorization: Implements distributed consistency using consistent hashing algorithm, and performs authorization interception with tokens to block illegal traffic
  - Distributed Quantity Service: Implements an interface service for quantity control to solve the overselling problem
  - RabbitMQ Message Queue: Web servers write messages to the message queue, and a separate service consumes them asynchronously
- Current Progress:
  - Completed development of Redis
  - Implemented high-concurrency distributed validation
  - Token-based authorization with the ability to horizontally scale
  - Resolved data consistency issues

## How to Use
### MyRedis
- Configure in my_redis/redis.conf:
  - Bind IP
  - Listening port
  - Enable AOF persistence and specify the AOF file name
  - Self address and addresses of other distributed nodes
- Run my_redis/main.go
- After establishing a TCP connection, use the following commands to operate Redis:
~~~
//*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
 //*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n
 //*2\r\n$6\r\nselect\r\n$1\r\n1\r\n
~~~

### Distributed Seckill System
- Configure in seckill_mysense/configs/config.json:
  - backend/api_port: Binding port
  - backend/api_secret: Project secret key, a randomly generated string
  - mysql/pass_word: MySQL password
  - mysql/host: MySQL server IP address
  - db_name: Database name
- Run seckill_mysense/main.go
- Access ip:port/check to participate in the seckill

## 🔗link
My Blog: [Kenway'Blog](http://kenway-20.com/)

# MyRedis和分布式秒杀系统
- 技术栈：
  - Go、MySQL、Redis、RabbitMQ
- 重点模块：
  - MyRedis中间件：独立开发的具备Redis功能的中间件MyRedis
  - 分布式权限验证：采用一致性hash算法实现分布式，使用token进行权限认证拦截部分非法流量
  - 分布式数量服务：实现数量控制的接口服务，用于解决超卖问题
  - 消息队列RabbitMQ：由后端web服务器写入消息队列，另起服务进行异步消费
- 目前进度：
  - 完成Redis的开发
  - 实现高并发分布式验证
  - token权限验证和结构可横向扩展
  - 解决数据一致性问题

## 如何使用
### MyRedis
- 在my_redis/redis.conf中配置
  - 绑定ip
  - 监听端口
  - 是否开启aof持久化及aof文件名称
  - 自身地址和其他分布式节点地址
- 运行my_redis/main.go文件
- tcp连接后，使用以下命令操作redis
~~~
//*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
//*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n
//*2\r\n$6\r\nselect\r\n$1\r\n1\r\n
~~~

### 分布式秒杀系统
- 在seckill_mysense/configs/config.json中配置
  - backend/api_port：绑定端口
  - backend/api_secret：项目密钥，一个随机生成的字符串
  - mysql/pass_word：mysql密码
  - mysql/host：mysql运行ip地址
  - db_name：数据库名称
- 运行seckill_mysense/main.go文件
- 通过访问ip:port/check进行抢购


## 🔗 链接
我的博客：[Kenway'Blog](http://kenway-20.com/)
