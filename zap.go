package zap_log

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	ZapLog *zap.Logger
	Sync   func() error
}

// Logger 配置zap日志,将zap日志库引入
func Logger() *ZapLogger {
	//配置zap日志库的编码器
	encoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
}

// NewZapLogger return a zap logger.
func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	//日志切割
	// writeSyncer := getLogWriter()
	//设置日志级别
	level.SetLevel(zap.InfoLevel)

	//开发模式下打印到标准输出
	// --根据配置文件判断输出到控制台还是日志文件--

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),                         // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台和文件
		level, // 日志级别
	)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{ZapLog: zapLogger, Sync: zapLogger.Sync}
}

// Log 实现log接口
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.ZapLog.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.ZapLog.Debug("", data...)
	case log.LevelInfo:
		l.ZapLog.Info("", data...)
	case log.LevelWarn:
		l.ZapLog.Warn("", data...)
	case log.LevelError:
		l.ZapLog.Error("", data...)
	case log.LevelFatal:
		l.ZapLog.Fatal("", data...)
	}
	return nil
}
