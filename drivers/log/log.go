package log

import (
	"fmt"
	"gtf/drivers/uuid"
	// "runtime"
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

func plainText(v ...interface{}) {
	log.Output("PLAINTEXT", LFileAndConsole, fmt.Sprintln(v...))
}

func DoPanic(err error) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		Error(err)
	// 		var buf []byte = make([]byte, 1500)
	// 		runtime.Stack(buf, true)
	// 		Warning("===========================-------------------------")
	// 		log.Output("PANIC", LFileAndConsole, fmt.Sprintf("%s\n", buf))
	// 	}
	// }()
	Warning("===========================-------------------------123")
	uuid := uuid.Rand()
	log.panicLocation = &uuid
	plainText(fmt.Sprintf("%s", uuid))
	panic(err)
}
