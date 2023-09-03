package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name              string `mapstructure:"name"`
	Mode              string `mapstructure:"mode"`
	Version           string `mapstructure:"version"`
	Port              string `mapstructure:"port"`
	StartTime         string `mapstructure:"start_time"`
	MachineID         int64  `mapstructure:"machine_id"`
	ReportFilePath    string `mapstructure:"report_file_path"`
	*LogConfig        `mapstructure:"log"`
	*MySQLConfig      `mapstructure:"mysql"`
	*RedisConfig      `mapstructure:"redis"`
	*RedisQueueConfig `mapstructure:"redisQueue"`
	*AsynqConfig      `mapstructure:"asynq"`
}

type LogConfig struct {
	Level      string `mapstructure:"name"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PollSize int    `mapstructure:"pool_size"`
}

type RedisQueueConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PollSize int    `mapstructure:"pool_size"`
}

type AsynqConfig struct {
	Concurrency int `mapstructure:"concurrency"`
}

func Init() (err error) {
	viper.SetConfigName("config") // 指定配置文件名称
	viper.SetConfigType("yaml")   // 指定配置文件类型
	viper.AddConfigPath("./conf") // 指定
	err = viper.ReadInConfig()    // 读取配置信息
	if err != nil {
		zap.L().Info(fmt.Sprintf("viper read error:%v", err))
		return
	}

	// 把读取到的配置信息反序列化到 conf 变量中
	if err = viper.Unmarshal(Conf); err != nil {
		zap.L().Info(fmt.Sprintf("viper Unmarshal error:%v", err))
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info(fmt.Sprintf("viper setting update success:%v", in.String()))
		if err = viper.Unmarshal(Conf); err != nil {
			zap.L().Info(fmt.Sprintf("viper Unmarshal2 error:%v", err))
		}
	})

	return
}
