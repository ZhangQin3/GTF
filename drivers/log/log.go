package log

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
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
	if err != nil {
		t := time.Now().UnixNano()
		log.panicTime = t
		plainText(fmt.Sprintf("panic_here_%d", t))
		panic(err)
	}
}

func DoCatch(f interface{}, params ...interface{}) {
	method := reflect.ValueOf(f)

	if method.Kind() != reflect.Func {
		panic("The first param of the gtf.Execute must be a testcase method!")
	}

	len := len(params)
	vParams := make([]reflect.Value, len)
	if len != 0 {
		for i := 0; i < len; i++ {
			vParams[i] = reflect.ValueOf((params)[i])
		}
	}

	defer func() {
		if err := recover(); err != nil {
			Error(err)
			var buf []byte = make([]byte, 1500)
			runtime.Stack(buf, true)
			Warning("===========================-------------------------")
			log.Output("PANIC", LFileAndConsole, fmt.Sprintf("%s\n", buf))
		}
	}()

	ret := method.Call(vParams)
	Warning("--------", ret[0])
}
