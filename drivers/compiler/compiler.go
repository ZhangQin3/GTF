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

// All the dir operations below are relative to %GOPATH%/src/gtf/drivers/compile/ directory

type scriptScheme struct {
	Script      string `json:",omitempty"`
	Repetitions int    `json:",omitempty"`
	Other       string `json:",omitempty"`
}

func main() {
	compileGtfPkg()
	common.CompileStdGoPkg("gse")
	CompileTestScripts()
	CompileExecuteGoFile("execute.go")

	// Remove All temp files
	os.RemoveAll(`temp`)
}

func compileGtfPkg() {
	common.CompileStdGoPkg("gtf/drivers/common")
	common.CompileMultiFilesPkg("log", common.GtfPkgDir)
	common.CompileMultiFilesPkg("gtf", common.GoPkgDir)
}

var (
	imports     = "import ("
	pkgs        = "gtf.TestSuiteSchema.TestScripts = map[string]interface{}{\n"
	repetitions = "gtf.TestSuiteSchema.Repetitions = map[string]int{\n"
)

func CompileTestScripts() {
	testSuiteScheme := decodeTestSuiteScheme()
	if len(testSuiteScheme) == 0 {
		panic("There is not any test script in the testSuiteScheme.")
	}

	for _, obj := range testSuiteScheme {
		goFileName := obj.Script
		fmt.Println(goFileName)
		// Test file doex NOT exist.
		if !common.IsFileExist(common.ScriptsSrcDir, goFileName) {
			fmt.Println("[WARNNING]: Test file " + goFileName + " does NOT exist!")
			continue
		}
		common.CompileSingleFilePkg(goFileName, common.ScriptsSrcDir, common.ScriptsPkgDir)
		appendExecuteInfo(goFileName, obj.Repetitions)
	}
	GenerateExecuteGoFile()
}

func appendExecuteInfo(goFileName string, rep int) {
	goBaseName := strings.TrimSuffix(goFileName, ".go")
	imports = imports + fmt.Sprintf("%s \"gtf/scripts/%s\"\n", goBaseName, goBaseName)
	pkgs = pkgs + fmt.Sprintf("\"%s\": new(%s.Test),\n", goBaseName, goBaseName)

	if rep == 0 {
		rep = 1 // Execute each test file at least one time.
	}
	repetitions = repetitions + fmt.Sprintf("\"%s\": %d,\n", goBaseName, rep)
}

func encloseExecuteInfo() {
	imports = imports + "\n)"
	pkgs = pkgs + "\t}"
	repetitions = repetitions + "\t}"
}

func CompileExecuteGoFile(fileName string) {
	var doComepile bool = true
	var filePrefix = strings.TrimSuffix(fileName, ".go")
	var execFileName = filePrefix + ".exe"
	var pkgFileName = ` temp/` + filePrefix + ".a"

	if common.IsFileExist(common.GoBinDir, execFileName) {
		pkgModTime := common.GetFileDate(common.GoBinDir, execFileName)
		goModTime := common.GetFileDate(`../execute/`, fileName)
		if pkgModTime.After(goModTime) {
			doComepile = false
		}
	}

	if doComepile {
		proLevel := os.Getenv("PROCESSOR_LEVEL")
		common.ExecOSCmd("go tool %sg -o %s -I %s -pack ../execute/%s", proLevel, pkgFileName, common.GoPkgDir, fileName)
		common.ExecOSCmd("go tool %sl -o %s%s -L %s%s", proLevel, common.GoBinDir, execFileName, common.GoPkgDir, pkgFileName)
	}
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
