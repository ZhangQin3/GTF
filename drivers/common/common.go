package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	GoPkgDir       = PkgDir()
	GoPath         = os.Getenv("gopath")
	ScriptsPkgDir  = PkgDir() + "gtf/scripts/"
	ScriptsSrcDir  = GoPath + "src/gtf/scripts/"
	TSPkgDir       = PkgDir() + "gtf/testsuites/"
	TSSrcDir       = GoPath + "src/gtf/testsuites/"
	GtfPkgDir      = PkgDir() + "gtf/"
	GoBinDir       = GoPath + "bin/"
	ProcessorLevel = os.Getenv("PROCESSOR_LEVEL")
)

func init() {

	if _, err := os.Stat(GoBinDir); os.IsNotExist(err) {
		os.MkdirAll(GoBinDir, 0777)
	}

	if _, err := os.Stat(TSPkgDir); os.IsNotExist(err) {
		os.MkdirAll(TSPkgDir, 0777)
	}

	if _, err := os.Stat(ScriptsPkgDir); os.IsNotExist(err) {
		os.MkdirAll(ScriptsPkgDir, 0777)
	}
}

func GetFileDate(fileDir string, fileName string) time.Time {
	fileInfo, err := os.Stat(fileDir + fileName)
	if err != nil {
		panic(err)
	}
	return fileInfo.ModTime()
}

// If the package has been compiled
func IsFileExist(dir, fileName string) bool {
	if _, err := os.Stat(dir + fileName); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ExecOSCmd(cmdStr string, args ...interface{}) {
	if len(args) != 0 {
		cmdStr = fmt.Sprintf(cmdStr, args...)
	}
	fmt.Printf("[DEGUG]: %s\n", cmdStr)
	out, err := exec.Command("cmd.exe", "/c", cmdStr).CombinedOutput()
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", out)
		panic(err)
	}
	if len(out) != 0 {
		fmt.Printf("%s\n", out)
	}
}

func PkgDir() string {
	var osType = "windows"
	if os.Getenv("OSTYPE") == "linux" {
		osType = "linux"
	}

	if ProcessorLevel == "6" {
		return fmt.Sprintf("%spkg/%s_amd64/", GoPath, osType)
	} else {
		return fmt.Sprintf("%spkg/%s_386/", GoPath, osType)
	}
}

func BinDir() string {
	gopath := os.Getenv("gopath")
	return gopath + `bin\`
}

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

// the StdGoPkg means a package which follows the requirements of go cmd tool
func CompileStdGoPkg(pkgName string) {
	ExecOSCmd(`go install ` + pkgName)
}

// Compile all go files in a subdir of the dir drivers to a pkg
// and put the pkg to the pkgloc
func CompileGoFilesInDir(dirName, pkgLoc string) {
	files, _ := filepath.Glob(fmt.Sprintf("../%s/*.go", dirName))
	input := strings.Join(files, " ")
	ExecOSCmd("go tool %sg -o %s%s.a -I %s -pack %s", ProcessorLevel, pkgLoc, dirName, GoPkgDir, input)
}

// CompileGoFile compiles packages with only single go file.
// The subdir MUST be the dir name in the $GOPATH\src dir, such as tests, ...
func CompileSingleGoFile(fileName, fileDir, pkgLoc string) {
	var doComepile = true
	var filePrefix = strings.TrimSuffix(fileName, ".go")
	var pkgFileName = filePrefix + ".a"

	if IsFileExist(GoPkgDir, pkgFileName) {
		pkgModTime := GetFileDate(GoPkgDir, pkgFileName)
		goModTime := GetFileDate(fileDir, fileName)
		if pkgModTime.After(goModTime) {
			doComepile = false
		}
	}
	if doComepile {
		ExecOSCmd("go tool %sg -o %s%s -I %s -pack %s%s", ProcessorLevel, pkgLoc, pkgFileName, GoPkgDir, fileDir, fileName)
	}
}

func CompileGtfCompiler() {
	ExecOSCmd(`go tool 6g -o %scompiler.a -I %s -pack ../compiler/compiler.go`, GoPkgDir, GoPkgDir)
	ExecOSCmd(`go tool 6l -o %scompiler.exe -L %s %scompiler.a`, GoPkgDir, GoPkgDir, GoPkgDir)
}
