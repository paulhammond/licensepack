package main

import (
	"flag"
	"fmt"
	"os"
)

//go:generate go run github.com/paulhammond/licensepack .
var licenses string

func main() {
	var credits = flag.Bool("credits", false, "show open source credits")

	flag.Parse()
	if *credits {
		fmt.Println(licenses)
		os.Exit(0)
	}

	fmt.Println("hello, world")
}
