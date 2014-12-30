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
func Step(index interface{}, description string, v ...interface{}) {
	var idx string
	switch i := index.(type) {
	case int:
		idx = strconv.Itoa(i)
	case float64:
		// Only support x.0-x.9 due to the prec is 1 here.
		idx = strconv.FormatFloat(i, 'f', 1, 64)
	case string:
		idx = i
		if idx == "." {
			if log.currentStep.Index == "PreTest" {
				idx = "1"
			}
			preId, err := strconv.Atoi(log.currentStep.Index)
			if err == nil {
				idx = strconv.Itoa(preId + 1)
			}
		}
	default:
		panic("Type does not support.")
	}

	// Calculate the step index that has "." index or index same with its previous one.
	idx = calculateStepIndex(idx)

	// Handle the step index that is same with any one before it or has the same main index.
	idx = handleRepeatStepIndex(idx)

	dscr := fmt.Sprintf(description, v...)
	step := &tcStep{Index: idx, Description: dscr}
	log.Steps = append(log.Steps, step)
	log.currentStep = step

	log.Output("STEP",
		LFileAndConsole,
		stepInfo{
			idx,
			time.Now().Format("2006-01-02 15:04:05"),
			dscr,
		})
}

var regexpStepIndex = regexp.MustCompile(`(\w+)\.?(\d*)`)

func calculateStepIndex(index string) string {
	match := regexpStepIndex.FindAllStringSubmatch(log.currentStep.Index, 1)
	if match != nil && match[0][1] == index || index == "." {
		if match[0][2] == "" && match[0][2] == "" {
			index = match[0][1] + `.1`
		} else if match[0][2] == "" {
			index = index + `.1`
		} else {
			s, err := strconv.Atoi(match[0][2])
			if err != nil {
				panic(err)
			}
			index = match[0][1] + `.` + strconv.Itoa(s+1)
		}
	} else if index == "." {
		if match[0][2] == "" {
			index = match[0][1] + `.1`
		} else {
			s, err := strconv.Atoi(match[0][2])
			if err != nil {
				panic(err)
			}
			index = match[0][1] + `.` + strconv.Itoa(s+1)
		}
	}

	return index
}

func handleRepeatStepIndex(index string) string {
	if log.metaIndex == nil {
		log.metaIndex = make(map[string]string)
	}

	if !strings.Contains(index, ".") {
		if i, ok := log.metaIndex[index]; ok {
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
			log.metaIndex[index] = i
			index = i
		} else {
			log.metaIndex[index] = index
		}
	} else {
		match := regexpStepIndex.FindAllStringSubmatch(index, 1)
		if match != nil {
			log.metaIndex[match[0][1]] = index
		}
	}
	return index
}
