package main

import (
    "fmt"
    "os"
    "./goyaml"
)

func main() {
    if len(os.Args) != 5 {
        fmt.Println("Usage: ./godoc2yml <package_source_path> <yaml_output_path> <package_prefix> <source_repo>")
        os.Exit(-1)
    }

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
