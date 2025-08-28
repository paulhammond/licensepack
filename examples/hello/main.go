package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
)

//go:generate go run github.com/paulhammond/licensepack .
//go:embed credits.txt
var credits string

func main() {
	showCredits := flag.Bool("credits", false, "show open source credits")

	flag.Parse()
	if *showCredits {
		fmt.Println(credits)
		os.Exit(0)
	}

	fmt.Println("hello, world")
}
