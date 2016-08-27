package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	GoPkgDir         = PkgDir()
	GoPath           = os.Getenv(`gopath`)
	ScriptsPkgDir    = PkgDir() + `gtf\scripts\`
	ScriptsSrcDir    = GoPath + `src\gtf\scripts\`
	DataFilesDir     = GoPath + `src\gtf\datafiles\`
	AWFilesDir       = GoPath + `src\gtf\awfiles\`
	TsPkgDir         = PkgDir() + `gtf\testsuites\`
	TsSrcDir         = GoPath + `src\gtf\testsuites\`
	GtfSrcDir        = GoPath + `src\`
	GtfPkgDir        = PkgDir() + `gtf\`
	GoBinDir         = GoPath + `bin\`
	GtfDriversPkgDir = PkgDir() + `gtf\drivers\`
	ProcessorLevel   = os.Getenv(`PROCESSOR_LEVEL`)
	DriversDir       = driversDir()
)

func init() {
	os.MkdirAll(GoBinDir, 0777)
	os.MkdirAll(TsPkgDir, 0777)
	os.MkdirAll(ScriptsPkgDir, 0777)
}

// the StdGoPkg means a package which follows the requirements of go cmd tool
func CompileStdGoPkg(pkgName string) {
	ExecOSCmd(`go install ` + pkgName)
}

// the StdGoPkg means a package which follows the requirements of go cmd tool
func CompileStdGoPkgInDir(dir string) {
	full_dir := GtfSrcDir + dir
	filepath.Walk(full_dir, func(fp string, fi os.FileInfo, err error) error {
		if fi.IsDir() && fp != full_dir {
			ExecOSCmd(`go install ` + dir + "/" + fi.Name())
		}
		return nil
	})
}

// Compile all go files in a subdir of the dir drivers to a pkg and put the pkg to the pkgLocation
func CompileMultiFilesPkg(dirName, pkgLocation string) {
	files, _ := filepath.Glob(fmt.Sprintf(`..\%s\*.go`, dirName))
	input := strings.Join(files, " ")
	ExecOSCmd("go tool compile -o %s%s.a -I %s -pack %s", pkgLocation, dirName, GoPkgDir, input)
}

// CompileGoFile compiles packages with only single go file.
func CompileSingleFilePkg(fileName, fileDir, pkgLocation string) {
	var filePrefix = strings.TrimSuffix(fileName, ".go")
	var pkgFileName = filePrefix + ".a"

	if DoesFileExist(GoPkgDir, pkgFileName) {
		pkgModTime := GetFileModTime(GoPkgDir, pkgFileName)
		goModTime := GetFileModTime(fileDir, fileName)
		if pkgModTime.After(goModTime) {
			return
		}
	}
	ExecOSCmd("go tool compile -o %s%s -I %s -pack %s%s", pkgLocation, pkgFileName, GoPkgDir, fileDir, fileName)
}

func CompileGtfCompiler() {
	ExecOSCmd(`go tool compile -o %scompiler.a -I %s -pack ..\compiler\compiler.go`, GoPkgDir, GoPkgDir)
	ExecOSCmd(`go tool link -o %scompiler.exe -L %s %scompiler.a`, GoPkgDir, GoPkgDir, GoPkgDir)
}

func PkgDir() string {
	if ProcessorLevel == "6" {
		return fmt.Sprintf(`%spkg\windows_amd64\`, GoPath)
	} else {
		return fmt.Sprintf(`%spkg\windows_386\`, GoPath)
	}
}

func BinDir() string {
	// gopath := os.Getenv("gopath")
	return GoPath + `bin\`
}

func driversDir() string {
	_, commonFile, _, _ := runtime.Caller(0)
	dir, _ := path.Split(path.Dir(commonFile))
	return dir
}

func GetFileModTime(fileDir string, fileName string) time.Time {
	fileInfo, err := os.Stat(fileDir + fileName)
	if err != nil {
		panic(err)
	}
	return fileInfo.ModTime()
}

func DoesFileExist(dir, fileName string) bool {
	if _, err := os.Stat(dir + fileName); err == nil {
		return true
	} else {
		return false
	}
}

func CopyFile(dst, src string) (w int64, err error) {
	d, err := os.Create(dst)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer d.Close()

	s, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer s.Close()

	return io.Copy(d, s)
}

func ExecOSCmd(cmd string, args ...interface{}) {
	if args != nil {
		cmd = fmt.Sprintf(cmd, args...)
	}

	fmt.Printf("[DEGUG]: %s\n", cmd)

	out, err := exec.Command("cmd.exe", "/c", cmd).CombinedOutput()
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", out)
		panic(err)
	}

	fmt.Printf("%s\n", out)
}
