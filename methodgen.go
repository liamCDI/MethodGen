//Written by Liam Kelly to try to save a lot of time porting a ton a C-code

//The following utility is meant to make generating methods for structures easier
//by leveraging the go tempaltes and AST libraries. The utility does this by creating
//a simplified version of the struct's field declarations and passing them to a given template.
//
//A simple template to create an `Equals` function that check if all field are equal could look like the following
//
//  package {{.Pkg}}
//
//  func (a *{{.Name}})Equals(b *{{.Name}})bool{{"{"}}
//      {{range .Fields}}if a.{{.Name}}!=b.{{.Name}}{
//          return false
//      }
//      {{end}}
//      return true
//  {{"}"}}
//
//This template would be stored in a seperate file. In this example the file will be called `equal.go.tmpl`.
//
//To render the function you would first install this utility in your PATH, then place a
//`go:generate` directive in your code calling this function. For a struct called `SomeStruct` and a template
//`equal.go.tmpl` the directive would look like:
//
//  //go:generate methodgen -tmpl=$GOPATH/src/tmpls/equal.go.tmpl -struct=SomeStruct
//
//Finally run the `go generate` to generate the function in its own file with the name <struct name>_<template name>.go
//
//The above example assumes the generate directive is in the same file as `SomeStruct`, if it is not then the file
//must be specified with the `-in` argument
package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//Field is meant to represent the field declartions in a struct
type Field struct {
	Name  string //Name of field
	Type  string //Type of Field. For an Array or slice, this is the element type
	Tag   string //Tag of field
	Slice bool   //Slice or Array Flag
	Len   int    //Only used if Slice flag to true. If the value is not 0, then element is considered an Array with that length
}

//Struct is meant to
type Struct struct {
	Name   string  //Name of struct
	Pkg    string  //Name of Package
	Fields []Field //Fields of struct
}

func main() {
	fin := flag.String("in", "", "File to be read. If omited, program will try to read enviroment variables set by `go generate`")
	strc := flag.String("struct", "", "struct to add method to")
	tmpl := flag.String("tmpl", "", "template file to render results")
	pkg := flag.String("pkg", "", "package name. If omited, program will try to read enviroment variables set by `go generate`")

	flag.Parse()

	epkg := os.Getenv("GOPACKAGE")
	if epkg == "" {
		epkg = *pkg
	}

	in := os.Getenv("GOFILE")

	if *fin == "" && in == "" {
		log.Fatal("need an input file specified")
	}
	if in == "" {
		in = *fin
	}

	if *strc == "" || *tmpl == "" {
		log.Fatal("need struct name, file where it is defined, and  template file to generate function")
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, in, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	out := &Struct{*strc, epkg, []Field{}}

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == *strc {
				fields := x.Type.(*ast.StructType).Fields.List
				var tag string
				for _, field := range fields {
					switch field.Type.(type) {
					case *ast.Ident:
						stype := field.Type.(*ast.Ident).Name
						tag = ""
						if field.Tag != nil {
							tag = field.Tag.Value
						}
						out.Fields = append(out.Fields, Field{field.Names[0].Name, stype, tag, false, 0})
					case *ast.ArrayType:
						tag = ""
						if field.Tag != nil {
							tag = field.Tag.Value
						}
						aryln, _ := strconv.Atoi(field.Type.(*ast.ArrayType).Len.(*ast.BasicLit).Value)
						out.Fields = append(out.Fields, Field{field.Names[0].Name, field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name, tag, true, aryln})
					}
				}
			}
		}
		return true
	})

	rndr, err := template.ParseFiles(*tmpl)
	if err != nil {
		log.Fatal(err)
	}

	base := filepath.Base(*tmpl)
	idx := strings.IndexByte(base, '.')
	if idx > 0 {
		base = base[0:idx]
	}

	ofn := path.Dir(in) + "/" + out.Name + "_" + base + ".go"
	ofd, err := os.Create(ofn)
	if err != nil {
		log.Fatal(err)
	}
	defer ofd.Close()

	err = rndr.Execute(ofd, out)
	if err != nil {
		log.Fatal(err)
	}

}
