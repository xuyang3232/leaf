package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// levels
const (
	debugLevel   = 0
	infoLevel = 1
	warnLevel = 2
	errorLevel   = 3
	panicLevel   = 4
	fatalLevel   = 5
)

const (
	printDebugLevel   = "[debug ] "
	printInfoLevel 	  = "[info ] "
	printWarnLevel 	  = "[warn ] "
	printErrorLevel   = "[error ] "
	printPanicLevel   = "[panic ] "
	printFatalLevel   = "[fatal ] "
)

type LoggerLeaf struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File
}

func NewLoggerLeaf(strLevel string, pathname string, flag int) (*LoggerLeaf, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "info":
		level = infoLevel
	case "warn":
		level = warnLevel
	case "error":
		level = errorLevel
	case "panic":
		level = panicLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	if pathname != "" {
		now := time.Now()

		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		file, err := os.Create(path.Join(pathname, filename))
		if err != nil {
			return nil, err
		}

		baseLogger = log.New(file, "", flag)
		baseFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(LoggerLeaf)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile

	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *LoggerLeaf) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *LoggerLeaf) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

	if level == fatalLevel {
		os.Exit(1)
	}else if level == panicLevel{
		panic(fmt.Sprintf(format,a...))
	}
}

func (logger *LoggerLeaf) doPrint(level int, printLevel string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	logger.baseLogger.Output(3, fmt.Sprintf(printLevel,a...))

	if level == fatalLevel {
		os.Exit(1)
	}else if level == panicLevel{
		panic(fmt.Sprintf(printLevel,a...))
	}
}

func (logger *LoggerLeaf) Debugf(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *LoggerLeaf) Infof(format string, a ...interface{}) {
	logger.doPrintf(infoLevel, printInfoLevel, format, a...)
}

func (logger *LoggerLeaf) Warnf(format string, a ...interface{}) {
	logger.doPrintf(warnLevel, printWarnLevel, format, a...)
}

func (logger *LoggerLeaf) Errorf(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *LoggerLeaf) Panicf(format string, a ...interface{}) {
	logger.doPrintf(panicLevel, printPanicLevel, format, a...)
}

func (logger *LoggerLeaf) Fatalf(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func (logger *LoggerLeaf) Debug(a ...interface{}) {
	logger.doPrint(debugLevel, printDebugLevel, a...)
}

func (logger *LoggerLeaf) Info( a ...interface{}) {
	logger.doPrint(infoLevel, printInfoLevel, a...)
}

func (logger *LoggerLeaf) Warn( a ...interface{}) {
	logger.doPrint(warnLevel, printWarnLevel, a...)
}

func (logger *LoggerLeaf) Error(a ...interface{}) {
	logger.doPrint(errorLevel, printErrorLevel, a...)
}

func (logger *LoggerLeaf) Panic(a ...interface{}) {
	logger.doPrint(panicLevel, printPanicLevel, a...)
}

func (logger *LoggerLeaf) Fatal(a ...interface{}) {
	logger.doPrint(fatalLevel, printFatalLevel, a...)
}


