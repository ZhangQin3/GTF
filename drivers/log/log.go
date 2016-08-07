package log

import (
	"fmt"
)

type flags int

const (
	LOnlyConsol flags = 1 << iota
	LOnlyFile
	Ldebug
	LFileAndConsole flags = LOnlyConsol | LOnlyFile
)

var log *Logger

func Text(v ...interface{}) {
	log.Output("DEFAULT", LFileAndConsole, fmt.Sprintln(v...))
}

func Textf(format string, v ...interface{}) {
	log.Output("DEFAULT", LFileAndConsole, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	log.Output("INFO", LFileAndConsole, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	log.Output("INFO", LFileAndConsole, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	log.currentStep.IsFailed = true
	log.Output("ERROR", LFileAndConsole, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	log.currentStep.IsFailed = true
	log.Output("ERROR", LFileAndConsole, fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) {
	log.Output("WARNING", LFileAndConsole, fmt.Sprintln(v...))
}

func Warningf(format string, v ...interface{}) {
	log.Output("WARNING", LFileAndConsole, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	log.Output("DEBUG", LFileAndConsole, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	log.Output("DEBUG", LFileAndConsole, fmt.Sprintf(format, v...))
}

func DoPanic(err error) {
	Warning("panic here", err)
	panic(err)
}
