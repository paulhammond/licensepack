package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"

	"github.com/paulhammond/licensepack/license"
)

const comment = "Code generated by licensepack; DO NOT EDIT."

// there are simpler examples of licensepack usage in the examples directory
//
//go:generate go run github.com/paulhammond/licensepack .
var licenses string

//go:embed tmpl/*.tmpl
var tmplFS embed.FS

func WrapQuote(indent, s string) string {
	quoted := fmt.Sprintf("%q", s)
	wrapped := strings.ReplaceAll(quoted, `\n`, `\n" +`+"\n"+indent+`"`)
	cleaned := strings.TrimSuffix(wrapped, ` +`+"\n"+indent+`""`)
	return cleaned
}

func parseTemplate(name string) (*template.Template, error) {
	t := template.New("")
	t.Funcs(template.FuncMap{
		"eval": func(name string, arg interface{}) (string, error) {
			var buf bytes.Buffer
			err := t.ExecuteTemplate(&buf, name, arg)
			return buf.String(), err
		},
		"quote": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
		"trim":      strings.TrimSpace,
		"wrapquote": WrapQuote,
	})

	fs2, err := fs.Sub(tmplFS, "tmpl")
	if err != nil {
		panic(err)
	}

	if regexp.MustCompile("^[a-z\\-]+$").MatchString(name) {
		name = name + ".tmpl"
		_, err = t.ParseFS(fs2, name)
	} else {
		_, err = t.ParseFiles(name)
	}
	if err != nil {
		return nil, err
	}
	t = t.Lookup(filepath.Base(name))
	if t == nil {
		return nil, errors.New("could not find template")
	}
	return t, nil
}

func main() {
	var (
		tmpl     = flag.String("tmpl", "string", "code template")
		pkg      = flag.String("pkg", "main", "package for generated code")
		file     = flag.String("file", "licenses.go", "filename for generated code (- for stdout)")
		variable = flag.String("var", "licenses", "variable name")
		nofmt    = flag.Bool("nofmt", false, "do not run output through gofmt")
		credits  = flag.Bool("credits", false, "show open source credits")
		help     = flag.Bool("help", false, "show help")
	)

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: licensepack [options] [pkg]")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *credits {
		fmt.Println("This software builds on many open source projects. We're grateful to the")
		fmt.Println("developers of the following projects for their hard work.")
		fmt.Println("")
		fmt.Println(licenses)
		os.Exit(0)
	}

	if *help || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedDeps | packages.NeedImports | packages.NeedName | packages.NeedModule}
	pkgs, err := packages.Load(cfg, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading packages: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	paths := map[string]string{
		"golang.org/pkg": build.Default.GOROOT,
	}
	mainPath := ""
	packages.Visit(pkgs, func(p *packages.Package) bool {
		// here we assume that the only packages that don't have a module are in
		// the standard library. Maybe that's not true?
		if p.Module != nil {
			paths[p.Module.Path] = p.Module.Dir

			if p.Module.Main {
				mainPath = p.Module.Path
			}
		}
		return true
	}, nil)

	modules := make([]license.Module, 0, len(paths))
	for k, v := range paths {
		isMain := (k == mainPath)
		files := []license.File{}
		dirEntries, err := os.ReadDir(v)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range dirEntries {
			re := regexp.MustCompile(`(?i)(licen(c|s)e|notice|copying)`)
			if re.MatchString(entry.Name()) {
				path := filepath.Join(v, entry.Name())
				if ".go" == filepath.Ext(path) {
					continue
				}

				stat, err := os.Stat(path)
				if err != nil {
					log.Fatal(err)
				}
				if stat.IsDir() {
					continue
				}

				contents, err := os.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}

				if strings.Contains(string(contents), comment) {
					continue
				}

				trimmed := strings.TrimSpace(string(contents))

				files = append(files, license.File{Path: entry.Name(), Contents: trimmed})
			}
		}
		if len(files) == 0 && !isMain {
			fmt.Printf("Error: could not find license for %s", k)
			os.Exit(1)
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Path < files[j].Path
		})

		modules = append(modules, license.Module{
			Name:     k,
			Licenses: files,
		})
	}

	sort.Slice(modules, func(i, j int) bool {
		if modules[j].Name == mainPath {
			return false
		}
		return (modules[i].Name == mainPath) || modules[i].Name < modules[j].Name
	})

	var src bytes.Buffer

	t, err := parseTemplate(*tmpl)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(&src, struct {
		Pkg     string
		Var     string
		Comment string
		Modules []license.Module
	}{
		Pkg:     *pkg,
		Var:     *variable,
		Modules: modules,
		Comment: comment,
	})
	if err != nil {
		log.Fatal(err)
	}

	output := src.Bytes()
	if !*nofmt {
		output, err = format.Source(src.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}

	if *file == "-" {
		_, err = fmt.Print(string(output))
	} else {
		err = os.WriteFile(*file, output, 0666)
	}

	if err != nil {
		log.Fatal(err)
	}
}
