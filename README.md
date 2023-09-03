# rpc-go-php
rpc 导入数据和删除; go语言作为服务端，php作为客户端，简单的rpc调用（主要是go语言的学习，切勿在生产中使用）

**特别注意：配置文件中的数据库相关信息要修改**

> go服务端
> 
`vi conf/config.yaml`

    report_file_path: "/var/www/report"    // import 导入数据时，生成报告的地址

    **特别注意：/var/www/tmp 用户组权限是 mysql:mysql **
    tmp_file_path: "/var/www/tmp"          // delete 操作中需要回收时，先查询数据写入txt保留的地址

    monitor_remote_addr: "192.168.0.110"   // 修改为自己的服务器地址


    mysql:
        host: "127.0.0.1"
        port: 3306
        user: "root"            // 修改为自己的登录名
        password: "123456"      // 修改为自己的登录密码
        dbname: "test"          // 修改为自己的数据库名
        max_open_conns: 20
        max_idle_conns: 14

    redis:
        host: "127.0.0.1"
        port: 6379
        password: ""            // 修改为自己的密码
        db: 0
        poolSize: 10
    
    redisQueue:
        addr: "127.0.0.1:6379"
        password: ""           // 修改为自己的密码
        db: 1
        poolSize: 10

运行项目前请先根据实际情况修改配置文件、文件夹权限

> rpc 服务启动
>
cd go项目所在文件夹

启动：`go run main.go`

开发环境可使用air，具体如何使用请自行查询

> php 调用

根据 $sync 这个参数决定 php 调用方是否同步阻塞

true =》同步请求等rpc返回

false =》直接结束

调用：php import.php


> 结果演示

![image](https://github.com/Alke-meng/rpc-go-help/blob/main/images/1.jpg)

![image](https://github.com/Alke-meng/rpc-go-help/blob/main/images/2.jpg)

![image](https://github.com/Alke-meng/rpc-go-help/blob/main/images/3.jpg)

![image](https://github.com/Alke-meng/rpc-go-help/blob/main/images/4.jpg)

![image](https://github.com/Alke-meng/rpc-go-help/blob/main/images/5.jpg)


> 简单的性能监控：
> http://192.168.0.110:9089/debug/statsview#
>
> 192.168.0.110 是自己服务器地址 9089 根据 conf/config.yaml 中 monitor_remote_addr 配置而来
> 
> 详细内容可查看 https://github.com/go-echarts/statsview
>


> 异步任务 asynq：
> 启动： ./asynqmon -port 9087  --redis-url=redis://:@localhost:6379/1
> http://192.168.0.110:9087/
>
> 192.168.0.110 是自己服务器地址 9087 启动地址
>
> 详细内容可查看 https://github.com/hibiken/asynq
> ，asynqmon 二进制文件可自行下载使用
>

### 该项目只是学习使用golang,rpc简单利用,实际生成使用请参考以下方式

> golang rpc 使用建议选择专业的rpc框架

> php 与 golang 的通信可参考
> https://github.com/roadrunner-server/roadrunner
> 
> 或者采用swoole,hyperf 框架中 goTask
> https://github.com/hyperf/gotask
