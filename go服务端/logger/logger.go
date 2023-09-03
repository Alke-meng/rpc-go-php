package logger

import (
	"ccgo/controllers"
	"ccgo/settings"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxBackup,
		MaxBackups: maxAge,
	}
	return zapcore.AddSync(lumberjackLogger)
}

// 框架默认的日志
func CCgoResponseLogger(status controllers.ResCode, method, data, traceID string, cost time.Duration) {
	zap.L().Debug("response",
		zap.String("traceID", traceID),
		zap.String("method", method),
		zap.String("data", data),
		zap.Int64("status", int64(status)),
		zap.Duration("cost", cost),
	)
}

func CCgoRequestLogger(method, data string) {
	zap.L().Debug("request",
		zap.String("method", method),
		zap.String("data", data),
	)
}

func CCgoActionLogger(traceID, info string) {
	zap.L().Debug("action",
		zap.String("traceID", traceID),
		zap.String("info", info),
	)
}

func CCgoTaskLogger(traceID, info string) {
	zap.L().Debug("task",
		zap.String("traceID", traceID),
		zap.String("info", info),
	)
}

func CCgoTaskCostLogger(traceID, info string, cost time.Duration) {
	zap.L().Debug("task",
		zap.String("traceID", traceID),
		zap.String("info", info),
		zap.Duration("cost", cost),
	)
}

func Init(cfg *settings.LogConfig, mode string) (err error) {
	writeSyncer := getLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return
	}

	var core zapcore.Core
	if mode == "dev" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg)
	return err
}
