

> MethodGen
The following utility is meant to make generating methods for structures easier
by leveraging the go tempaltes and AST libraries. The utility does this by creating
a simplified version of the struct's field declarations and passing them to a given template.

A simple template to create an `Equals` function that check if all field are equal could look like the following


	package {{.Pkg}}
	
	func (a *{{.Name}})Equals(b *{{.Name}})bool{{"{"}}
	    {{range .Fields}}if a.{{.Name}}!=b.{{.Name}}{
	        return false
	    }
	    {{end}}
	    return true
	{{"}"}}

This template would be stored in a seperate file. In this example the file will be called `equal.go.tmpl`.

To render the function you would first install this utility in your PATH, then place a
`go:generate` directive in your code calling this function. For a struct called `SomeStruct` and a template
`equal.go.tmpl` the directive would look like:


	//go:generate methodgen -tmpl=$GOPATH/src/tmpls/equal.go.tmpl -struct=SomeStruct

Finally run the `go generate` to generate the function in its own file with the name <struct name>_<template name>.go

The above example assumes the generate directive is in the same file as `SomeStruct`, if it is not then the file
must be specified with the `-in` argument

### Advanced Usage
The struct tag `mg` can be used to store key value pairs in a Field's Tag Map. This allows for custom triggers in the tempalate
see the `PrintTag.tmpl` for examples.

### Notes
-currently does not support Slice or nested Structs (getting there)






- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
