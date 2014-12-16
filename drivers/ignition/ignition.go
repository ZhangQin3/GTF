package main

import (
	"gtf/drivers/common"
)

func main() {
	compileGtfCompiler()
	CompileGtf()
	// RunExecute()
}

func compileGtfCompiler() {

	common.CompileStdGoPkg("gtf/drivers/common")
	common.CompileStdGoPkg("gtf/drivers/suite")
	// common.CompileMultiFilesPkg("suite", common.GtfPkgDir)
	common.CompileSingleFilePkg("tsuite.go", common.TsSrcDir, common.TsPkgDir)
	common.CompileGtfCompiler()
}

func CompileGtf() {
	common.ExecOSCmd("%scompiler.exe", common.GoPkgDir)
}

func RunExecute() {
	common.ExecOSCmd("%sexecute.exe", common.GoBinDir)
}
