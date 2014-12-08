package log

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type stepInfo struct {
	StepIndex string
	Time      string
	Text      string
}

type tcStep struct {
	StepIndex       string
	StepDescription string
	IsFailed        bool
}

// LogStep(1.1, "Setup xxxx"), 1.1 is "main index"."subIndex"
func Step(StepIndex interface{}, f string, v ...interface{}) {
	var idx string
	switch si := StepIndex.(type) {
	case int:
		idx = strconv.Itoa(si)
	case float64:
		// Only support x.0-x.9 due to the prec is 1 here.
		// TODO: Enhance the issue.
		idx = strconv.FormatFloat(si, 'f', 1, 64)
	case string:
		idx = si
		if idx == "." {
			if log.currentStep.StepIndex == "PRE-FIRST-STEP" {
				idx = "1"
			}
			preId, err := strconv.Atoi(log.currentStep.StepIndex)
			if err == nil {
				idx = strconv.Itoa(preId + 1)
			}
		}
	default:
		panic("Type does not support.")
	}
	// Calculate the step index that has "." index or index same with its previous one.
	calculateStepIndex(&idx, `(\w+)\.?(\d*)`)
	// Handle the step index that is same with any one before it or has the same main index.
	handleRepeatStepIndex(&idx)

	StepDescription := fmt.Sprintf(f, v...)
	stp := &tcStep{StepIndex: idx, StepDescription: StepDescription}
	log.Steps = append(log.Steps, stp)
	log.currentStep = stp

	log.Output("STEP",
		LFileAndConsole,
		stepInfo{
			idx,
			time.Now().Format("2006-01-02 15:04:05"),
			StepDescription,
		})
}

// LogStep(1.1, "%s", "Setup xxxx"), NOT exported for now
func stepf(StepIndex interface{}, format string, v ...interface{}) {
	var idx string
	switch si := StepIndex.(type) {
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
	StepDescription := fmt.Sprintf(format, v...)
	log.GenerateStep(idx, StepDescription)
	log.Output("STEP",
		LFileAndConsole,
		stepInfo{
			idx,
			time.Now().Format("2006-01-02 15:04:05"),
			StepDescription,
		})
}

func (l *Logger) GenerateStep(idx, dscr string) {
	step := &tcStep{StepIndex: idx, StepDescription: dscr}
	l.Steps = append(l.Steps, step)
	l.currentStep = step
}

func calculateStepIndex(idx *string, reg string) {
	regStepIndex, _ := regexp.Compile(reg)
	match := regStepIndex.FindAllStringSubmatch(log.currentStep.StepIndex, 1)
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

func handleRepeatStepIndex(idx *string) {
	regexpStepIndex, _ := regexp.Compile(`(\w+)\.?(\d*)`)
	if log.metaIndex == nil {
		log.metaIndex = make(map[string]string)
	}

	if !strings.Contains(*idx, ".") {
		if i, ok := log.metaIndex[*idx]; ok {
			match := regexpStepIndex.FindAllStringSubmatch(i, 1)
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
			log.metaIndex[*idx] = i
			*idx = i
		} else {
			log.metaIndex[*idx] = *idx
		}
	} else {
		match := regexpStepIndex.FindAllStringSubmatch(*idx, 1)
		if match != nil {
			log.metaIndex[match[0][1]] = *idx
		}
	}
}
