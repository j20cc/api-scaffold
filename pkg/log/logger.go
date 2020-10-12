package log

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var l *logger

type logger struct {
	engine *zap.Logger
	mode   string
}

// NewLogger init logger
func NewLogger() {
	if l != nil {
		return
	}

	l = &logger{
		mode: viper.GetString("app.mode"),
	}
	// zap hook
	hook := lumberjack.Logger{
		Filename:   viper.GetString("log.path"),     // 日志文件路径
		MaxSize:    viper.GetInt("log.max_size"),    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: viper.GetInt("log.max_backups"), // 日志文件最多保存多少个备份
		MaxAge:     viper.GetInt("log.max_age"),     // 文件最多保存多少天
		Compress:   true,                            // 是否压缩
	}
	// zap encoderConfig
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	_ = atomicLevel.UnmarshalText([]byte(viper.GetString("log.level")))
	var w zapcore.WriteSyncer
	if l.mode == "debug" {
		w = zapcore.AddSync(os.Stdout)
	} else {
		w = zapcore.AddSync(&hook)
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		atomicLevel,
	)
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", viper.GetString("app.name")))
	// 构造日志
	l.engine = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), development, filed)
}

// Debug log
func Debug(msg string, fields ...zap.Field) {
	l.engine.Debug(msg, fields...)
}

// Info log
func Info(msg string, fields ...zap.Field) {
	l.engine.Info(msg, fields...)
}

// Warn log
func Warn(msg string, fields ...zap.Field) {
	l.engine.Warn(msg, fields...)
}

// Error log
func Error(msg string, fields ...zap.Field) {
	l.engine.Error(msg, fields...)
}

// Fatal log
func Fatal(msg string, fields ...zap.Field) {
	l.engine.Fatal(msg, fields...)
}
