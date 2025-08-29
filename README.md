# licensepack

licensepack generates a text file containing licensing information for all
dependencies used by your Go project. This should help distributors of your
software follow the terms of most open source licenses.

## Usage

licensepack is intended to be run as a `go-generate` script. Start by adding
licensepack to your `go.mod` file:

```
go get -tool github.com/paulhammond/licensepack
```

Then add a [`go-generate`](https://go.dev/blog/generate) line to your code:
```
//go:generate go tool licensepack ./path/to/main
```

Then run `go generate`:

```
go generate ./...
```

A file called `credits.txt` will be created with contents looking something like
this:

```
## github.com/your/package

LICENSE:

Copyright (c) 2021 Your Name. All rights reserved.
…
```

You can then use the [embed package](https://pkg.go.dev/embed) to make that file
available in your code:

```
//go:embed credits.txt
var credits string
```

Then you can add some code to output the contents of that variable, perhaps if a
command line flag is provided:

```
func main(){
	…
	var showCredits = flag.Bool("credits", false, "show open source credits")

	flag.Parse()
	if *showCredits {
		fmt.Println(credits)
		os.Exit(0)
	}
	…
}
```

A complete example is available in [`examples/hello`](examples/hello).

### Built in templates

licensepack uses `text/template`, allowing you to customize the output to meet
your needs. A few built in templates are provided which provide differing levels
of repetitiveness:

```
//go:generate go tool licensepack -tmpl short path/to/main
//go:generate go tool licensepack -tmpl long path/to/main
```

### Custom templates

If you’d like to customize the template you can provide the path to a template
file. For example:

```
//go:generate go tool licensepack -tmpl ./licenses.tmpl .
```

A complete example is provided in `examples/custom`. The built in templates can
be used as a starting point and hopefully demonstrate usage of the available
template helper functions.

## Alternatives

[go-licenses](https://github.com/google/go-licenses) can generate similar
reports to licensepack (see the markdown example in the README). It has a lot
of additional functionality that licensepack lacks, including the ability to
detect a name for each license, but also has a lot of dependencies. If
licencepack doesn't have the functionality you need then this is likely the best
option for you.

[go-licence-detector](https://github.com/elastic/go-licence-detector) also
generates license notices. It has many options not available in licensepack, but
invocation is more complex and, at the time of writing, includes indirect
dependencies that are not actually linked into the compiled binary.

## Credits/License

Run `licensepack -credits` for license information.
