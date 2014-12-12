package log

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type tcStep struct {
	Index       string
	Description string
	IsFailed    bool
}

// LogStep(1.1, "Setup xxxx"), 1.1 is "main index"."subIndex"
func Step(stepIndex interface{}, stepDescription string, v ...interface{}) {
	var index string
	switch si := stepIndex.(type) {
	case int:
		index = strconv.Itoa(si)
	case float64:
		// Only support x.0-x.9 due to the prec is 1 here.
		// TODO: Enhance the issue.
		index = strconv.FormatFloat(si, 'f', 1, 64)
	case string:
		index = si
		if index == "." {
			if log.currentStep.Index == "PRE-FIRST-STEP" {
				index = "1"
			}
			preId, err := strconv.Atoi(log.currentStep.Index)
			if err == nil {
				index = strconv.Itoa(preId + 1)
			}
		}
	default:
		panic("Type does not support.")
	}
	// Calculate the step index that has "." index or index same with its previous one.
	calculateStepIndex(&index, `(\w+)\.?(\d*)`)
	// Handle the step index that is same with any one before it or has the same main index.
	handleRepeatStepIndex(&index)

	description := fmt.Sprintf(stepDescription, v...)
	stp := &tcStep{Index: index, Description: description}
	log.Steps = append(log.Steps, stp)
	log.currentStep = stp

	log.Output("STEP",
		LFileAndConsole,
		stepInfo{
			index,
			time.Now().Format("2006-01-02 15:04:05"),
			description,
		})
}

// LogStep(1.1, "%s", "Setup xxxx"), NOT exported for now
func stepf(stepIndex interface{}, stepDescription string, v ...interface{}) {
	var index string
	switch si := stepIndex.(type) {
	case int:
		index = strconv.Itoa(si)
	case float64:
		// only support x.0-x.9 due to the prec is 1 here
		index = strconv.FormatFloat(si, 'f', 1, 64)
	case string:
		index = si
	default:
		panic("Type does not support.")
	}
	StepDescription := fmt.Sprintf(stepDescription, v...)
	log.GenerateStep(index, StepDescription)
	log.Output("STEP",
		LFileAndConsole,
		stepInfo{
			index,
			time.Now().Format("2006-01-02 15:04:05"),
			StepDescription,
		})
}

func (l *Logger) GenerateStep(index, dscr string) {
	step := &tcStep{Index: index, Description: dscr}
	l.Steps = append(l.Steps, step)
	l.currentStep = step
}

func calculateStepIndex(index *string, reg string) {
	regStepIndex, _ := regexp.Compile(reg)
	match := regStepIndex.FindAllStringSubmatch(log.currentStep.Index, 1)
	Warning(match)
	if match != nil && match[0][1] == *index || *index == "." {
		if match[0][2] == "" && match[0][2] == "" {
			*index = match[0][1] + `.1`
		} else if match[0][2] == "" {
			*index = *index + `.1`
		} else {
			s, err := strconv.Atoi(match[0][2])
			if err != nil {
				panic(err)
			}
			*index = match[0][1] + `.` + strconv.Itoa(s+1)
		}
	} else if *index == "." {
		if match[0][2] == "" {
			*index = match[0][1] + `.1`
		} else {
			s, err := strconv.Atoi(match[0][2])
			if err != nil {
				panic(err)
			}
			*index = match[0][1] + `.` + strconv.Itoa(s+1)
		}
	}
}

func handleRepeatStepIndex(index *string) {
	regexpStepIndex, _ := regexp.Compile(`(\w+)\.?(\d*)`)
	if log.metaIndex == nil {
		log.metaIndex = make(map[string]string)
	}

	if !strings.Contains(*index, ".") {
		if i, ok := log.metaIndex[*index]; ok {
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
			log.metaIndex[*index] = i
			*index = i
		} else {
			log.metaIndex[*index] = *index
		}
	} else {
		match := regexpStepIndex.FindAllStringSubmatch(*index, 1)
		if match != nil {
			log.metaIndex[match[0][1]] = *index
		}
	}
}
