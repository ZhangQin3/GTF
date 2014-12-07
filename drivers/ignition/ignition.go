package main

import (
	"gtf/drivers/common"
)

func main() {
	compileGtfCompiler()
	CompileGtf()
}

func compileGtfCompiler() {

	common.CompileStdGoPkg("gtf/drivers/common")
	common.CompileMultiFilesPkg("suite", common.GtfPkgDir)
	common.CompileSingleFilePkg("tsuite.go", common.TSSrcDir, common.TSPkgDir)
	common.CompileGtfCompiler()
}

func CompileGtf() {
	common.ExecOSCmd("%scompiler.exe", common.GoPkgDir)
}
