package log

var gLogger Logger = NewLoggerZap("debug", "")


type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Close()
}

func Export(logger Logger) {
	gLogger = logger
}

func Debug(v ...interface{}) {
	gLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	gLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	gLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	gLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	gLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	gLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	gLogger.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	gLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	gLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	gLogger.Fatalf(format, v...)
}

func Close() {
	gLogger.Close()
}
