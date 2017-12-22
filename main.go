package main

import (
    "os"
    "./goyaml"
)

func main() {
	// parameters
	// arg1 directory which contains the package,
	packageSource := os.Args[1]
	// arg2 package name
	packageName := os.Args[2]
	// arg3 output directory, wher package.yml saved to
	ymlOutput := os.Args[3]
	goyaml.GoYAMLGeneration(packageSource, packageName, ymlOutput)
}

