package main

import (
    "fmt"
    "os"
    "./goyaml"
)

func main() {
    if len(os.Args) != 4 {
        fmt.Println("Usage: ./godoc2yml <package_source_path> <package_name> <yaml_output_path>")
        os.Exit(-1)
    }

    // arg1 directory which contains the package,
    packageSource := os.Args[1]
    // arg2 package name
    packageName := os.Args[2]
    // arg3 output directory, wher package.yml saved to
    ymlOutput := os.Args[3]
    goyaml.GoYAMLGeneration(packageSource, packageName, ymlOutput)
}

