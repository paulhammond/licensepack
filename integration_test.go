package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"licensepack": main,
	})
}

func TestCLI(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting cwd: %v", err)
	}
	p := testscript.Params{
		Dir: "testdata",
		Setup: func(env *testscript.Env) error {
			env.Setenv("src", cwd)
			return nil
		},
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}
