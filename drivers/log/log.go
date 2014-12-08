package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"text/template"
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu          sync.Mutex // ensures atomic writes; protects the following fields.
	text        string     // for accumulating text to write
	stdOut      io.Writer  // destination for output
	file        *os.File   // file destination for the output, if any.
	fileName    string
	template    *template.Template // the template used to writ log to file
	Steps       []*tcStep
	currentStep *tcStep
	metaIndex   map[string]string
}

type TcHeaderInfo struct {
	TcID string
	Time string
	Text string
}

type TsHeaderInfo struct {
	Time string
	Text string
}

type TcResultSummary struct {
	TcID        string
	Description string
	Result      bool
	FailedSteps string
}

type flags int

const (
	LOnlyConsol flags = 1 << iota
	LOnlyFile
	Ldebug
	LFileAndConsole flags = LOnlyConsol | LOnlyFile
)

var log *Logger

// NewLogger creates a new Logger. the Logger can writes log to the console and
// a specific file.
func NewLogger(fileName, tmplFielName string) *Logger {
	t, err := template.ParseFiles(tmplFielName)
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	log = &Logger{stdOut: os.Stdout, file: logFile, template: t, fileName: fileName}
	return log
}

// level: DEBUG, DEFAULT, INFO, WARNING, ERROR, SETUP, D_PASS, D_FAIL,
func (l *Logger) Output(level string, flags flags, info interface{}) {
	if l == nil {
		panic("The Logger l is not initialized.")
	}

	var data interface{}
	switch v := info.(type) {
	case nil:
		l.text = ""
		data = v
	case string:
		l.text = v
		data = v
	case stepInfo:
		l.text = v.Text
		data = v
	case TcHeaderInfo, TsHeaderInfo:
		data = v
	default:
		panic("Can't accept the info type.")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if flags&LOnlyConsol != 0 {
		// write to consol
		if len(l.text) > 0 && l.text[len(l.text)-1] != '\n' {
			l.text = l.text + "\n"
		}
		if _, err := l.stdOut.Write([]byte("[" + level + "] " + l.text)); err != nil {
			panic(err)
		}
	}

	if flags&LOnlyFile != 0 {
		// write to file
		if l.file == nil {
			panic("The Log Output file is nil.")
		}
		if err := l.template.ExecuteTemplate(l.file, level, data); err != nil {
			panic(err)
		}
	}
}

func (l *Logger) ReopenFile() {
	if l == nil {
		panic("The Logger log is not initialized.")
	}
	logFile, err := os.OpenFile(l.fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	l.file = logFile
}

func (l *Logger) CloseFile() {
	if l == nil {
		panic("The Logger log is not initialized.")
	}
	l.file.Close()
}

func (l *Logger) GetFileName() string {
	return l.fileName
}

func (l *Logger) GetTemplate() *template.Template {
	return l.template
}

func Log(v ...interface{}) {
	log.Output("INFO", LFileAndConsole, fmt.Sprintln(v...))
}

func Logf(format string, v ...interface{}) {
	log.Output("INFO", LFileAndConsole, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	if len(log.Steps) == 0 {
		log.GenerateStep("PRE-TEST", "PRE-TEST")
	}
	log.currentStep.IsFailed = true
	log.Output("ERROR", LFileAndConsole, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	if len(log.Steps) == 0 {
		log.GenerateStep("PRE-TEST", "PRE-TEST")
	}
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
