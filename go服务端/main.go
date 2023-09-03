package main

import (
	"ccgo/asynq"
	"ccgo/dao/asynqClient"
	"ccgo/dao/mysql"
	"ccgo/dao/redis"
	"ccgo/logger"
	"ccgo/routes"
	"ccgo/settings"
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-echarts/statsview"
	"github.com/go-echarts/statsview/viewer"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 1、加载配置
	if err := settings.Init(); err != nil {
		zap.L().Info(fmt.Sprintf("init setting err:%v", err))
		return
	}
	// 2、初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		zap.L().Info(fmt.Sprintf("logger init fail:%v", err))
		return
	}
	defer zap.L().Sync()

	defer mysql.Close()
	// 3、初始化MySQL
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		zap.L().Info(fmt.Sprintf("mysql init fail:%v", err))
		return
	}

	defer redis.Close()
	// 4、初始化Redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Info(fmt.Sprintf("redis init fail:%v", err))
		return
	}

	// 5、注册rpc
	routes.Setup()

	// 6、启动服务
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("port")))
	if err != nil {
		zap.L().Info(fmt.Sprintf("rpc init fail:%v", err))
		panic(err)
	}

	zap.L().Info(fmt.Sprintf("rpc init success,listen port:%v", viper.GetInt("port")))

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			_ = conn
			go jsonrpc.ServeConn(conn)
		}
	}()

	//7、任务投递
	defer asynqClient.Close()
	asynqClient.Init(settings.Conf.RedisQueueConfig)

	//8、启动异步任务
	go func() {
		if err := asynq.SetUp(settings.Conf.RedisQueueConfig); err != nil {
			zap.L().Info(fmt.Sprintf("asynq run server fail:%v", err))
			panic(err)
		}
	}()

	//9、启动监控
	viewer.SetConfiguration(
		viewer.WithInterval(viper.GetInt("monitor_interval")),
		viewer.WithAddr(fmt.Sprintf(":%d", viper.GetInt("monitor_local_port"))),
		viewer.WithLinkAddr(fmt.Sprintf("%s:%d", viper.GetString("monitor_remote_addr"), viper.GetInt("monitor_local_port"))),
	)
	go func() {
		mgr := statsview.New()
		mgr.Start()
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞

	<-quit // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	zap.L().Info("Server exiting")

}
