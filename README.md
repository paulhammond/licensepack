# licensepack

licensepack generates Go source files containing licensing information for all
dependencies of your Go project. This should help distributors of your software
follow the terms of most open source licenses.

## Usage

licensepack is intended to be run as a `go-generate` script. Start by adding
licensepack to your `go.mod` file:

```
go get -d github.com/paulhammond/licensepack
```

Then add a `go-generate` line to your code, along with a variable definition. By
default licensepack uses a variable called licenses:

```
//go:generate go run github.com/paulhammond/licensepack ./path/to/main
var licenses string
```

Then run `go generate`:

```
go generate
```

A file called `licenses.go` will be created with contents looking something like
this:

```
// Code generated by licensepack; DO NOT EDIT.
package main

func init() {
	licenses = "## github.com/your/package\n" +
		"\n" +
		"LICENSE:\n" +
		"\n" +
		"Copyright (c) 2009 Your Name. All rights reserved.\n" +
		…
}
```

Then you can add some code to output the contents of the `licenses` variable,
perhaps if a command line flag is provided:

```
func main(){
	…
	var credits = flag.Bool("credits", false, "show open source credits")

	flag.Parse()
	if *credits {
		fmt.Println(licenses)
		os.Exit(0)
	}
	…
}
```

A complete example is available in [`examples/hello`](examples/hello).

By default licensepack generates code wrapped in an `init()` function so that
you can avoid committing the generated file without complaints from the Go
tooling. By recreating the file every time your project is compiled you ensure
the license information is always kept up to date, at the cost of adding
`go generate` to your build scripts. This is customizable, see below for
details.

### Adjusting the generated code

Simple changes, including package name, variable name and the output filename
can be adjusted with command line flags:

```
//go:generate go run github.com/paulhammond/licensepack -pkg cmd -var Credits -file path/to/credits.go ./cmd
```

### Built in templates

licensepack uses `text/template` to generate code, allowing you to customize the
output to meet your needs. A few built in templates are provided, including one
that creates a struct instead of a string:

```
//go:generate go run github.com/paulhammond/licensepack ./cmd -tmpl struct
```

A complete example of using this template is provided in `examples/struct`.

If you believe that `init()` functions are always bad or would prefer to commit
the generated file, then templates that generate a variable, const or function
declaration are also provided with the `-var`, `-const`, and `-func` suffixes.
For example:

```
//go:generate go run github.com/paulhammond/licensepack -tmpl struct-func .
//go:generate go run github.com/paulhammond/licensepack -tmpl string-const .
```

If you’d prefer to generate plain text output, perhaps because you’d like to
embed licensing information using a different mechanism, there is a plain text
template available. To use this you’ll need to disable gofmt. For example:

```
//go:generate go run github.com/paulhammond/licensepack -file licenses.txt -nofmt -tmpl text ./cmd
```

### Custom templates

If you’d like to customize the template you can provide the path to a template
file. For example:

```
//go:generate go run github.com/paulhammond/licensepack -tmpl ./licenses.tmpl ./cmd
```

A complete example is provided in `examples/custom`. The built in templates can
be used as a starting point and hopefully demonstrate usage of the available
template helper functions. By default the generated code is formatted before
being written to a file so you can be lax with your whitespace.

### Keeping licensepack in go.mod

If you follow the instructions above, then run `go mod tidy`, you will find that
licensepack has been removed from your `go.mod` file. To prevent this from
happening you have two options:

1. If you have an existing `tools.go` file, as [recommended by the Go team][tools],
you can add licensepack to it.

2. If that seems like too much work you can import
`github.com/paulhammond/licensepack/license` in your code and declare your
variable as a `license.String`, which is an alias for `string`.

[tools]: https://github.com/golang/go/issues/25922

## Alternatives

[go-licence-detector](https://github.com/elastic/go-licence-detector) also
generates license notices. It has many options not available in licensepack, but
invocation is more complex and, at the time of writing, includes indirect
dependencies that are not actually linked into the compiled binary.

## Credits/License

Run `licensecheck -credits` for license information.
