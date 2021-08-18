package main

import (
	"bytes"
	"embed"
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

// Options
var (
	tmpl     = flag.String("tmpl", "init", "code template")
	pkg      = flag.String("pkg", "main", "package for generated code")
	file     = flag.String("file", "licenses.go", "filename for generated code")
	variable = flag.String("var", "Licenses", "variable name")
	credits  = flag.Bool("credits", false, "show open source credits")
	help     = flag.Bool("help", false, "show help")
)

// Non-options
const comment = "Code generated by licensepack; DO NOT EDIT."

//go:generate go run github.com/paulhammond/licensepack .
var Licenses license.String

//go:embed tmpl/*.tmpl
var tmplFS embed.FS

func mustFS(fs fs.FS, err error) fs.FS {
	if err != nil {
		panic(err)
	}
	return fs
}

func WrapQuote(indent, s string) string {
	quoted := fmt.Sprintf("%q", s)
	wrapped := strings.ReplaceAll(quoted, `\n`, `\n" +`+"\n"+indent+`"`)
	cleaned := strings.TrimSuffix(wrapped, ` +`+"\n"+indent+`""`)
	return cleaned
}

func templates() *template.Template {
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
	template.Must(t.ParseFS(mustFS(fs.Sub(tmplFS, "tmpl")), "*"))
	return t
}

func main() {
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
		fmt.Println(Licenses)
		os.Exit(0)
	}

	if *help || len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	fmt.Printf("licensepack: generating %s\n", *file)

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
		return (modules[i].Name == mainPath) || modules[i].Name < modules[j].Name
	})

	var src bytes.Buffer
	_, err = src.Write([]byte("// " + comment + "\n"))
	if err != nil {
		log.Fatal(err)
	}

	t := templates()
	err = t.ExecuteTemplate(&src, *tmpl+".tmpl", struct {
		Pkg     string
		Var     string
		Modules []license.Module
	}{
		Pkg:     *pkg,
		Var:     *variable,
		Modules: modules,
	})
	if err != nil {
		log.Fatal(err)
	}

	formatted, err := format.Source(src.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(*file, formatted, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
