package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/paulhammond/licensepack/license"
)

//go:generate go run github.com/paulhammond/licensepack -tmpl struct .
var licenses []license.Module

func main() {
	var credits = flag.Bool("credits", false, "show open source credits")

	flag.Parse()

	if *credits {

		data := map[string]map[string]string{}

		for _, module := range licenses {
			if module.Name == "github.com/paulhammond/licensepack/examples/struct" {
				continue
			}
			files := map[string]string{}
			for _, file := range module.Licenses {
				files[file.Path] = file.Contents
			}
			data[module.Name] = files
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(&data)

		os.Exit(0)
	}

	fmt.Println("hello, world")
}
