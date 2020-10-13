package log

import (
	"net/http"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var l *logger

type logger struct {
	engine *zap.SugaredLogger
	mode   string
	level  zap.AtomicLevel
}

// NewLogger init logger
func NewLogger() {
	if l != nil {
		return
	}

	l = &logger{
		mode: viper.GetString("app.mode"),
	}
	// 设置日志级别
	atom := zap.NewAtomicLevel()
	if err := atom.UnmarshalText([]byte(viper.GetString("log.level"))); err != nil {
		panic(err)
	}
	// 默认输出终端，线上输出到json文件
	w := zapcore.AddSync(os.Stdout)
	e := getDevEncoder()
	if l.mode != "debug" {
		w = getProLogWriter()
		e = getProEncoder()
	}
	core := zapcore.NewCore(e, w, atom)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	l.engine = logger.Sugar()
	l.level = atom
	defer l.engine.Sync()
}

func getDevEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getProEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getProLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   viper.GetString("log.path"),     // 日志文件路径
		MaxSize:    viper.GetInt("log.max_size"),    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: viper.GetInt("log.max_backups"), // 日志文件最多保存多少个备份
		MaxAge:     viper.GetInt("log.max_age"),     // 文件最多保存多少天
		Compress:   true,                            // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// ServeHTTP can set logger level in fly
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.level.ServeHTTP(w, r)
}

// Debug log
func Debug(args ...interface{}) {
	l.engine.Debug(args)
}

// Info log
func Info(args ...interface{}) {
	l.engine.Info(args)
}

// Warn log
func Warn(args ...interface{}) {
	l.engine.Warn(args)
}

// Error log
func Error(args ...interface{}) {
	l.engine.Error(args)
}
