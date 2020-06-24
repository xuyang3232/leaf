package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)


type LoggerZap struct {
	sugarLogger *zap.SugaredLogger
}

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.DebugLevel
}

func NewLoggerZap(strLevel string, pathname string) *LoggerZap{
	if pathname != ""{
		return newLoggerZapFile(strLevel,pathname)
	}
	return newLoggerZapStdOut(strLevel, pathname)
}

func newLoggerZapFile(strLevel string, pathname string)*LoggerZap{
	encoder := getEncoder()
	logLevel := getLoggerLevel(strLevel)
	defaultLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	logName := path.Join(pathname, "default.log")
	errLogName := path.Join(pathname, "error.log")
	infoWriter := getLogWriter(logName)
	errorWriter := getLogWriter(errLogName)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), defaultLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	log := zap.New(core, zap.AddCaller(),zap.AddCallerSkip(2))
	return &LoggerZap{sugarLogger:log.Sugar()}
}

func newLoggerZapStdOut(strLevel string, pathname string)*LoggerZap{
	encoder := getEncoder()
	logLevel := getLoggerLevel(strLevel)

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(logLevel))
	log := zap.New(core, zap.AddCaller(),zap.AddCallerSkip(2))
	return &LoggerZap{sugarLogger:log.Sugar()}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 10,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func (logger *LoggerZap) Debugf(format string, a ...interface{}) {
	logger.sugarLogger.Debugf(format,a...)
}

func (logger *LoggerZap) Infof(format string, a ...interface{}) {
	logger.sugarLogger.Infof(format,a...)
}

func (logger *LoggerZap) Warnf(format string, a ...interface{}) {
	logger.sugarLogger.Warnf(format,a...)
}

func (logger *LoggerZap) Errorf(format string, a ...interface{}) {
	logger.sugarLogger.Errorf(format,a...)
}

func (logger *LoggerZap) Panicf(format string,a ...interface{}) {
	logger.sugarLogger.Panicf(format,a...)
}

func (logger *LoggerZap) Fatalf(format string, a ...interface{}) {
	logger.sugarLogger.Fatalf(format,a...)
}

func (logger *LoggerZap) Debug(a ...interface{}) {
	logger.sugarLogger.Debug(a...)
}

func (logger *LoggerZap) Info( a ...interface{}) {
	logger.sugarLogger.Info(a...)
}

func (logger *LoggerZap) Warn( a ...interface{}) {
	logger.sugarLogger.Warn(a...)
}

func (logger *LoggerZap) Error(a ...interface{}) {
	logger.sugarLogger.Error(a...)
}

func (logger *LoggerZap) Panic(a ...interface{}) {
	logger.sugarLogger.Panic(a...)
}

func (logger *LoggerZap) Fatal(a ...interface{}) {
	logger.sugarLogger.Fatal(a...)
}

func (logger *LoggerZap) Close() {
	logger.sugarLogger.Sync()
}

