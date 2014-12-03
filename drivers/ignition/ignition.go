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
	common.CompileGoFilesInDir("suite", common.GtfPkgDir)
	common.CompileSingleGoFile("tsuite.go", common.TSSrcDir, common.TSPkgDir)
	common.CompileGtfCompiler()
}

func CompileGtf() {
	common.ExecOSCmd("%scompiler.exe", common.GoPkgDir)
}
