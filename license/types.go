package license

type String string

type Module struct {
	Name     string
	Licenses []File
}

type File struct {
	Path     string
	Contents string
}
