package log

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu          sync.Mutex         // ensures atomic writes; protects the following fields.
	tmpl        *template.Template // the template used to writ log to file
	text        string             // for accumulating text to write
	stdOut      io.Writer          // destination for output
	logFile     *os.File           // file destination for the output, if any.
	logFileName string
	Steps       []*tcStep
	curStp      *tcStep
	metaIdx     map[string]string
}

type stepInfo struct {
	StepIdx string
	Time    string
	Text    string
}

type TcHdrInfo struct {
	TcID string
	Time string
	Text string
}

type TsHdrInfo struct {
	Time string
	Text string
}

type TcResultToTbl struct {
	TcID      string
	TcDescr   string
	TcResult  bool
	FaildStps string
}

const (
	LOnlyConsol = 1 << iota
	LOnlyFile
	LFileConsole = LOnlyConsol | LOnlyFile
	Ldebug       = 1 << iota
)

type tcStep struct {
	StpIdx    string
	StpDscr   string
	StpFailed bool
}

var log *Logger

// NewLogger creates a new Logger. the Logger can writes log to the console and
// a specific file.
func NewLogger(logFileName, tmplFielName string) *Logger {
	t, err := template.ParseFiles(tmplFielName)
	if err != nil {
		panic(err)
	}

	logFd, err := os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	log = &Logger{stdOut: os.Stdout, logFile: logFd, tmpl: t, logFileName: logFileName}
	return log
}

// level: DEBUG, DEFAULT, INFO, WARNING, ERROR, SETUP, D_PASS, D_FAIL,
func (l *Logger) Output(level string, writeTo int, info interface{}) {
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
	case TcHdrInfo, TsHdrInfo:
		data = v
	default:
		panic("Can't accept the info type.")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if writeTo&LOnlyConsol != 0 {
		// write to consol
		if len(l.text) > 0 && l.text[len(l.text)-1] != '\n' {
			l.text = l.text + "\n"
		}
		if _, err := l.stdOut.Write([]byte("[" + level + "] " + l.text)); err != nil {
			panic(err)
		}
	}

	if writeTo&LOnlyFile != 0 {
		// write to file
		if l.logFile == nil {
			panic("The Output file [logFile] is nil.")
		}
		if err := l.tmpl.ExecuteTemplate(l.logFile, level, data); err != nil {
			panic(err)
		}
	}
}

// >>>>>>>>>>>>>>>>>> Exported Functions >>>>>>>>>>>>>>>>>>>
// Not exported
func println(level string, v ...interface{}) {
	level = strings.ToUpper(level)
	log.Output(level, LFileConsole, fmt.Sprintln(v...))
}

func Print(level string, v ...interface{}) {
	level = strings.ToUpper(level)
	log.Output(level, LFileConsole, fmt.Sprintln(v...))
}

func Printf(level string, format string, v ...interface{}) {
	level = strings.ToUpper(level)
	log.Output(level, LFileConsole, fmt.Sprintf(format, v...))
}

func Log(v ...interface{}) {
	log.Output("INFO", LFileConsole, fmt.Sprintln(v...))
}

func Logf(format string, v ...interface{}) {
	log.Output("INFO", LFileConsole, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	if len(log.Steps) == 0 {
		log.GenStp("PRE-TEST", "PRE-TEST")
	}
	log.curStp.StpFailed = true
	log.Output("ERROR", LFileConsole, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	if len(log.Steps) == 0 {
		log.GenStp("PRE-TEST", "PRE-TEST")
	}
	log.curStp.StpFailed = true
	log.Output("ERROR", LFileConsole, fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) {
	log.Output("WARNING", LFileConsole, fmt.Sprintln(v...))
}

func Warningf(format string, v ...interface{}) {
	log.Output("WARNING", LFileConsole, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	log.Output("DEBUG", LFileConsole, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	log.Output("DEBUG", LFileConsole, fmt.Sprintf(format, v...))
}

// LogStep(1.1, "Setup xxxx"), 1.1 is "main index"."subIndex"
func Step(stepIdx interface{}, f string, v ...interface{}) {
	var idx string
	switch si := stepIdx.(type) {
	case int:
		idx = strconv.Itoa(si)
	case float64:
		// Only support x.0-x.9 due to the prec is 1 here.
		// TODO: Enhance the issue.
		idx = strconv.FormatFloat(si, 'f', 1, 64)
	case string:
		idx = si
		if idx == "." {
			if log.curStp.StpIdx == "PRE-FIRST-STEP" {
				idx = "1"
			}
			preId, err := strconv.Atoi(log.curStp.StpIdx)
			if err == nil {
				idx = strconv.Itoa(preId + 1)
			}
		}
	default:
		panic("Type does not support.")
	}
	// Calculate the step index that has "." index or index same with its previous one.
	calcStepIndex(&idx, `(\w+)\.?(\d*)`)
	// Handle the step index that is same with any one before it or has the same main index.
	handleRepeatIdx(&idx)

	stpDscr := fmt.Sprintf(f, v...)
	stp := &tcStep{StpIdx: idx, StpDscr: stpDscr}
	log.Steps = append(log.Steps, stp)
	log.curStp = stp

	log.Output("STEP",
		LFileConsole,
		stepInfo{
			idx,
			time.Now().Format("2006-01-02 15:04:05"),
			stpDscr,
		})
}

// LogStep(1.1, "%s", "Setup xxxx"), NOT exported for now
func stepf(stepIdx interface{}, format string, v ...interface{}) {
	var idx string
	switch si := stepIdx.(type) {
	case int:
		idx = strconv.Itoa(si)
	case float64:
		// only support x.0-x.9 due to the prec is 1 here
		idx = strconv.FormatFloat(si, 'f', 1, 64)
	case string:
		idx = si
	default:
		panic("Type does not support.")
	}
	stpDscr := fmt.Sprintf(format, v...)
	log.GenStp(idx, stpDscr)
	log.Output("STEP",
		LFileConsole,
		stepInfo{
			idx,
			time.Now().Format("2006-01-02 15:04:05"),
			stpDscr,
		})
}

// --------------------------------
func ReopenLogFile() {
	if log == nil {
		panic("The Logger log is not initialized.")
	}
	logFd, err := os.OpenFile(log.logFileName, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.logFile = logFd
}

func CloseLogFile() {
	if log == nil {
		panic("The Logger log is not initialized.")
	}
	log.logFile.Close()
}

func (l *Logger) GetLogTmpl() *template.Template {
	return l.tmpl
}

func (l *Logger) GetLogFileName() string {
	return l.logFileName
}

func (l *Logger) GenStp(idx, dscr string) {
	stp := &tcStep{StpIdx: idx, StpDscr: dscr}
	l.Steps = append(l.Steps, stp)
	l.curStp = stp
}

func calcStepIndex(idx *string, reg string) {
	regStpIdx, _ := regexp.Compile(reg)
	match := regStpIdx.FindAllStringSubmatch(log.curStp.StpIdx, 1)
	Warning(match)
	if match != nil && match[0][1] == *idx || *idx == "." {
		if match[0][2] == "" && match[0][2] == "" {
			*idx = match[0][1] + `.1`
		} else if match[0][2] == "" {
			*idx = *idx + `.1`
		} else {
			s, err := strconv.Atoi(match[0][2])
			if err != nil {
				panic(err)
			}
			*idx = match[0][1] + `.` + strconv.Itoa(s+1)
		}
	} else if *idx == "." {
		if match[0][2] == "" {
			*idx = match[0][1] + `.1`
		} else {
			s, err := strconv.Atoi(match[0][2])
			if err != nil {
				panic(err)
			}
			*idx = match[0][1] + `.` + strconv.Itoa(s+1)
		}
	}
}

func handleRepeatIdx(idx *string) {
	regStpIdx, _ := regexp.Compile(`(\w+)\.?(\d*)`)
	if log.metaIdx == nil {
		log.metaIdx = make(map[string]string)
	}

	if !strings.Contains(*idx, ".") {
		if i, ok := log.metaIdx[*idx]; ok {
			match := regStpIdx.FindAllStringSubmatch(i, 1)
			if match != nil {
				if match[0][2] == "" {
					i = i + `.1`
				} else {
					s, err := strconv.Atoi(match[0][2])
					if err != nil {
						panic(err)
					}
					i = match[0][1] + `.` + strconv.Itoa(s+1)
				}
			}
			log.metaIdx[*idx] = i
			*idx = i
		} else {
			log.metaIdx[*idx] = *idx
		}
	} else {
		match := regStpIdx.FindAllStringSubmatch(*idx, 1)
		if match != nil {
			log.metaIdx[match[0][1]] = *idx
		}
	}
}
