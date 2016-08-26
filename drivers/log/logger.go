package log

import (
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
	metaIndex   map[string]string // used to handle repeated step index.
	panicTime   int64
}

type TestcaseHdr struct {
	TcID   string
	Time   string
	Text   string
	Anchor int64
}

type TestScriptHdr struct {
	Time string
	Text string
}

type TestScriptTlr struct {
	Time string
	Text string
}

type TestResultSummary struct {
	TcID        string
	Description string
	Result      bool
	FailedSteps string
	TcAnchor    int64
	Duration    string // minute
}

type stepInfo struct {
	Index string
	Time  string
	Text  string
}

// NewLogger creates a new Logger. the Logger can writes log to the console and
// a specific file.
func NewLogger(fileName, tmplFielName string) *Logger {
	t, err := template.ParseFiles(tmplFielName)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	log = &Logger{stdOut: os.Stdout, file: file, template: t, fileName: fileName}
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
	case TestcaseHdr, TestScriptHdr, TestScriptTlr, toggleText, toggleImage:
		data = v
	default:
		panic("Log.Output: Can't accept the info type.")
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

func (l *Logger) FileName() string {
	return l.fileName
}

func (l *Logger) GetTemplate() *template.Template {
	return l.template
}

func (l *Logger) GenerateStep(index, description string) {
	step := &tcStep{Index: index, Description: description}
	l.Steps = append(l.Steps, step)
	l.currentStep = step
}

func (l *Logger) CleanupSteps() {
	l.Steps = l.Steps[0:0]
	l.metaIndex = nil
}

func (l *Logger) PanicTime() int64 {
	return l.panicTime
}

func (l *Logger) ZeroPanicTime() {
	l.panicTime = 0
}

func (l *Logger) LabelError() {
	l.currentStep.IsFailed = true
}
