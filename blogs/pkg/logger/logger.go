package logger

import (
	"blogs/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	log   *zap.Logger
	sugar *zap.SugaredLogger
)

func Init(cfg *config.LogConfig) error {
	// 日志编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 指定时间字段的键名,输出示例："time": "2024-01-15T10:30:45.123Z"
		LevelKey:       "level",                        // 指定日志级别字段的键名,输出示例："level": "info" 或 "level": "error"
		NameKey:        "logger",                       // 指定日志器名称字段的键名
		CallerKey:      "caller",                       // 指定调用位置字段的键名,输出示例："caller": "handlers/user.go:45"
		FunctionKey:    zapcore.OmitKey,                // 完全省略函数名
		MessageKey:     "msg",                          // 指定日志消息字段的键名,输出示例："msg": "user created successfully"
		StacktraceKey:  "stacktrace",                   // 指定堆栈跟踪字段的键名,当日志级别为 panic 或 fatal 时，会自动包含堆栈信息
		LineEnding:     zapcore.DefaultLineEnding,      // 指定每条日志的换行符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 指定日志级别的格式化方式为普通的"info", "error"
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 指定时间字段的格式化方式为"2024-01-15T10:30:45.123Z"
		EncodeDuration: zapcore.SecondsDurationEncoder, // 指定时间间隔字段的格式化方式为浮点数
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 指定调用位置字段的格式化方式为相对路径
	}

	// 日志级别
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel // 打印所有级别日志（debug/info/warn/error）
	case "info":
		level = zapcore.InfoLevel // 打印 info/warn/error
	case "warn":
		level = zapcore.WarnLevel // 打印 warn/error
	case "error":
		level = zapcore.ErrorLevel // 打印error
	}

	// 文件输出
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// 创建多个输出（文件 + 控制台）
	var writers []zapcore.WriteSyncer
	writers = append(writers, zapcore.AddSync(fileWriter))
	writers = append(writers, zapcore.AddSync(os.Stdout))

	// 核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),   // 日志格式JSON
		zapcore.NewMultiWriteSyncer(writers...), // 输出目标，文件+控制台
		level,                                   // 日志级别
	)

	// 创建 logger
	// zap.AddCaller()：打印日志调用的文件和行号（如 logger/logger.go:50）
	// zap.AddCallerSkip(1)：跳过当前日志封装层的调用栈，显示真实业务代码的调用位置
	// zap.AddStacktrace(zapcore.ErrorLevel)：仅在打印 error 级别日志时，附带栈追踪信息，方便排查错误
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = log.Sugar()

	return nil
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}

// Sync 同步日志缓冲区
func Sync() error {
	return log.Sync()
}

// GetLogger 获取原始 logger
func GetLogger() *zap.Logger {
	return log
}

// GetSugaredLogger 获取 sugared logger
func GetSugaredLogger() *zap.SugaredLogger {
	return sugar
}
