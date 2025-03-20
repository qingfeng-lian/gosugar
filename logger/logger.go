package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	//projectRedis "demo-project/src/datasource/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var loggerWrapper *zap.Logger
var isOpenDebug bool

type ZapLoggerConf struct {
	FileName   string `json:"file_name" yaml:"file_name"`     // 日志文件和路径 例如：log/application.log
	MaxSize    int    `json:"max_size" yaml:"max_size"`       // 单个日志文件的最大大小（以MB为单位）
	MaxAge     int    `json:"max_age" yaml:"max_age"`         // 保留旧文件的最大数量
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` // 旧文件保留的最大天数
	Compress   bool   `json:"compress" yaml:"compress"`       // 是否压缩旧文件
	Level      string `json:"level" yaml:"level"`             // 日志级别
}

type ZapLogger struct {
	conf ZapLoggerConf
}

// NewZapLogger 日志初始化
func NewZapLogger(conf ZapLoggerConf) error {
	var err error
	if conf.FileName == "" ||
		conf.MaxSize == 0 ||
		conf.MaxAge == 0 ||
		conf.MaxBackups == 0 {
		// 参数为空抛出异常
		return errors.New("logger config error")
	}
	// 编码器配置
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.LevelKey = "level"
	encoderCfg.NameKey = "logger"
	encoderCfg.CallerKey = "caller"
	//encoderCfg.CallerKey = ""
	encoderCfg.MessageKey = "msg"
	encoderCfg.StacktraceKey = "stacktrace"
	encoderCfg.LineEnding = zapcore.DefaultLineEnding
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	// encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeDuration = zapcore.SecondsDurationEncoder
	//encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	// 设置时区，例如：使用"Asia/Shanghai"，根据需要替换
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return errors.New("failed to load location")
	}
	// 自定义时间编码器，应用时区
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05.000"))
	}

	// Lumberjack 日志分割器
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.FileName,
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackups,
		Compress:   conf.Compress,
	})

	// 初始化日志级别
	atomicLevel := zap.NewAtomicLevel()
	_ = atomicLevel.UnmarshalText([]byte(conf.Level))

	// 同步写入
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer),
		atomicLevel.Level(),
	)

	loggerWrapper = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))

	return err
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		traceId := ctx.Value("traceId")
		if traceId != nil {
			fields = append(fields, zap.String("traceId", fmt.Sprint(traceId)))
		}
	}
	loggerWrapper.Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		traceId := ctx.Value("traceId")
		if traceId != nil {
			fields = append(fields, zap.String("traceId", fmt.Sprint(traceId)))
		}
	}
	loggerWrapper.Error(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		traceId := ctx.Value("traceId")
		if traceId != nil {
			fields = append(fields, zap.String("traceId", fmt.Sprint(traceId)))
		}
	}
	if isOpenDebug {
		Info(ctx, msg, fields...)
	} else {
		loggerWrapper.Debug(msg, fields...)
	}
}

func Infof(ctx context.Context, msg string, fields ...interface{}) {
	if ctx != nil {
		traceId := ctx.Value("traceId")
		if traceId != nil {
			fields = append(fields, zap.String("traceId", fmt.Sprint(traceId)))
		}
	}
	loggerWrapper.Sugar().Infof(msg, fields...)
}
