package main

import (
	"os"

	"./goyaml"
)

func main() {
	// parameters
	// arg1 source package path
	packageSource := os.Args[1]
	// arg2 output directory
	ymlOutput := os.Args[2]
	// arg3 package prefix, parent folder strings of package folder, for example: arm is the package prefix of package arm.advisor
	packagePrefix := os.Args[3]
	// arg4 source git repo
	sourceRepo := os.Args[4]
	goyaml.GoYAMLGeneration(packageSource, ymlOutput, packagePrefix, sourceRepo)
}
