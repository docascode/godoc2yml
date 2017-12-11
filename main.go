package main

import (
"./goyaml"
)

func main() {
	// parameters
	// arg1 directory which contains the package,
	packageSource := "D:/Bunny/Work/GitRepo/golang-sampleprojects/go-github"
	// arg2 package name
	packageName := "github"
	// arg3 output directory, wher package.yml saved to
	ymlOutput := "D:/Bunny/Work"
	goyaml.GoYAMLGeneration(packageSource, packageName, ymlOutput)
}

