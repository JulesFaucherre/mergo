package logger

import "log"

var (
	Verbosity = 0
)

func SetVerbosity(v int) {
	Verbosity = v
}

func Fatal(format string, v ...interface{}) {
	if Verbosity >= 0 {
		log.Fatalf(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if Verbosity >= 0 {
		log.Printf(format, v...)
	}
}

func Warning(format string, v ...interface{}) {
	if Verbosity >= 1 {
		log.Printf(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if Verbosity >= 2 {
		log.Printf(format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if Verbosity >= 3 {
		log.Printf(format, v...)
	}
}

func Silly(format string, v ...interface{}) {
	if Verbosity >= 4 {
		log.Printf(format, v...)
	}
}
