package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/paulhammond/licensepack/license"
)

//go:generate go run github.com/paulhammond/licensepack -tmpl struct .
var Licenses []license.Module

func main() {
	var credits = flag.Bool("credits", false, "show open source credits")

	flag.Parse()

	if *credits {
		for _, module := range Licenses {
			if module.Name == "github.com/paulhammond/licensepack/examples/struct" {
				continue
			}
			color.New(color.FgCyan, color.Bold).Println(module.Name)
			fmt.Println("")
			for _, file := range module.Licenses {
				color.New(color.FgWhite, color.Bold).Println(file.Path)
				fmt.Println("")
				color.White(file.Contents)
				fmt.Println("")
			}
		}

		os.Exit(0)
	}

	fmt.Println("hello, world")
}
