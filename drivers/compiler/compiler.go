package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gtf/drivers/common"
	tsuite "gtf/testsuites/tsuite"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type scriptScheme struct {
	Script      string `json:",omitempty"`
	Repetitions int    `json:",omitempty"`
	Other       string `json:",omitempty"`
}

func main() {
	compileGtfPkg()
	// compile all packages in the directory gtf/library
	common.CompileStdGoPkgInDir("gtf/library")
	CompileTestScripts()
	CompileExecuteGoFile("execute.go")

	// Remove All temp files
	os.RemoveAll(`temp`)
}

func compileGtfPkg() {
	common.CompileStdGoPkg("gtf/drivers/log")
	common.CompileStdGoPkg("gtf/drivers/assert")
	common.CompileStdGoPkg("gtf/drivers/csvdata")
	common.CompileStdGoPkg("gtf/drivers/uuid")
	common.CompileMultiFilesPkg("gtf", common.GtfDriversPkgDir)
}

var (
	imports     = "import ("
	pkgs        = "gtf.TestSuiteSchema.TestScripts = map[string]interface{}{\n"
	repetitions = "gtf.TestSuiteSchema.Repetitions = map[string]int{\n"
)

func CompileTestScripts() {
	scheme := decodeTestSuiteScheme()
	if len(scheme) == 0 {
		panic("There is not any test script in the testSuiteScheme.")
	}

	for _, obj := range scheme {
		fileName := obj.Script

		// Test file doex NOT exist.
		if !common.DoesFileExist(common.ScriptsSrcDir, fileName) {
			fmt.Println("[WARNNING]: Test file " + fileName + " does NOT exist!")
			continue
		}
		fmt.Println(fileName)
		// just compile go file, excluding the .csv file
		if strings.HasSuffix(fileName, ".go") {
			common.CompileSingleFilePkg(fileName, common.ScriptsSrcDir, common.ScriptsPkgDir)
		}
		appendExecuteInfo(fileName, obj)
	}
	GenerateExecuteGoFile()
}

func appendExecuteInfo(fileName string, s scriptScheme) {
	baseName := strings.TrimSuffix(fileName, ".go")
	imports = imports + fmt.Sprintf("%s `gtf/scripts/%s`\n", baseName, baseName)
	if strings.HasSuffix("fileName", ".csv") {
		fmt.Println("---------------1234")
		pkgs = pkgs + fmt.Sprintf("`%s`: `csv`,\n", baseName)
	} else {
		pkgs = pkgs + fmt.Sprintf("`%s`: new(%s.Test),\n", baseName, baseName)
	}

	if s.Repetitions == 0 {
		// Execute each test file at least one time.
		s.Repetitions = 1
	}
	repetitions = repetitions + fmt.Sprintf("`%s`: %d,\n", baseName, s.Repetitions)
}

func encloseExecuteInfo() {
	imports = imports + "\n)"
	pkgs = pkgs + "\t}"
	repetitions = repetitions + "\t}"
}

func CompileExecuteGoFile(fileName string) {
	var filePrefix = strings.TrimSuffix(fileName, ".go")
	var execFileName = filePrefix + ".exe"
	var pkgFileName = ` temp/` + filePrefix + ".a"

	if common.DoesFileExist(common.GoBinDir, execFileName) {
		pkgModTime := common.GetFileModTime(common.GoBinDir, execFileName)
		goModTime := common.GetFileModTime(`../execute/`, fileName)
		if pkgModTime.After(goModTime) {
			return
		}
	}

	common.ExecOSCmd("go tool compile -o %s -I %s -pack ../execute/%s", pkgFileName, common.GoPkgDir, fileName)
	common.ExecOSCmd("go tool link -o %s%s -L %s%s", common.GoBinDir, execFileName, common.GoPkgDir, pkgFileName)
}

func GenerateExecuteGoFile() {
	b, err := ioutil.ReadFile(`../execute/prototype/execute_prototype.go`)
	if err != nil {
		panic(err)
	}

	encloseExecuteInfo()

	impt := regexp.MustCompile(`//Import Here`)
	b = impt.ReplaceAll(b, []byte(imports))

	pkg := regexp.MustCompile(`//testPkgs`)
	b = pkg.ReplaceAll(b, []byte(pkgs))

	repet := regexp.MustCompile(`//repetitions`)
	b = repet.ReplaceAll(b, []byte(repetitions))

	ioutil.WriteFile(`../execute/execute.go`, b, 0644)
	common.ExecOSCmd(`gofmt -w ../execute/execute.go`)
}

func decodeTestSuiteScheme() []scriptScheme {
	var tSuite = new(tsuite.TSuite)
	var sch []scriptScheme
	tSuite.SuiteScheme()

	regexpCommentLine := regexp.MustCompile(`^ *//`)
	regexpNonNullLine := regexp.MustCompile(`\S+`)
	scanner := bufio.NewScanner(strings.NewReader(tSuite.Scheme))
	for scanner.Scan() {
		bytes := scanner.Bytes()
		if regexpNonNullLine.Match(bytes) && !regexpCommentLine.Match(bytes) {
			var s scriptScheme
			err := json.Unmarshal(bytes, &s)
			if err != nil {
				panic(err)
			}
			sch = append(sch, s)
		}
	}
	return sch
}

func init() {
	if _, err := os.Stat(`temp`); os.IsNotExist(err) {
		os.Mkdir(`temp`, 0777)
	}
}
