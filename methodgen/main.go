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
//
//Advanced Usage
//
//The struct tag `mg` can be used to store key value pairs in a Field's Tag Map. This allows for custom triggers in the tempalate
//see the `PrintTag.tmpl` for examples.
//
//Notes
//
// -currently does not support Slice or nested Structs (getting there)
package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//Field is meant to represent the field declartions in a struct
type Field struct {
	Name  string            //Name of field
	Type  string            //Type of Field. For an Array or slice, this is the element type
	Tag   map[string]string //Tag of field
	Slice bool              //Slice or Array Flag
	Len   int               //Only used if Slice flag to true. If the value is not 0, then element is considered an Array with that length
}

//Struct is meant to
type Struct struct {
	Name   string  //Name of struct
	Pkg    string  //Name of Package
	Fields []Field //Fields of struct
}

//HasTagKey Checks if key in tag config string. Good to use to check if an import is needed
func (s *Struct) HasTagKey(k string) bool {
	for _, f := range s.Fields {
		for tkey := range f.Tag {
			if k == tkey {
				return true
			}
		}
	}
	return false
}

var deftag = "mg"

//ProcTag processes the field tag for our tag 'mg'
func ProcTag(tagfull string) map[string]string {
	out := make(map[string]string)

	tagdecstmp := strings.Trim(tagfull, "`")
	tagdecs := strings.Split(tagdecstmp, " ")
	for _, tagdec := range tagdecs {
		tmptagdec := strings.Split(tagdec, ":")
		if tmptagdec[0] == deftag {
			if len(tmptagdec) == 2 {
				config := strings.Split(strings.Trim(tmptagdec[1], "\""), ",")
				for _, conf := range config {
					if strings.Contains(conf, "=") {
						tmpconf := strings.Split(conf, "=")
						out[tmpconf[0]] = tmpconf[1]
					} else {
						out[conf] = ""
					}
				}
			} else {
				return nil
			}
		}
	}
	return out
}

func noescape(s string) template.HTML {
	return template.HTML(s)
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
				var tag map[string]string
				for _, field := range fields {
					switch field.Type.(type) {
					case *ast.Ident:
						stype := field.Type.(*ast.Ident).Name
						tag = nil
						if field.Tag != nil {
							tag = ProcTag(field.Tag.Value)
						}
						out.Fields = append(out.Fields, Field{field.Names[0].Name, stype, tag, false, 0})
					case *ast.ArrayType:
						tag = nil
						if field.Tag != nil {
							tag = ProcTag(field.Tag.Value)
						}
						aryln, _ := strconv.Atoi(field.Type.(*ast.ArrayType).Len.(*ast.BasicLit).Value)
						out.Fields = append(out.Fields, Field{field.Names[0].Name, field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name, tag, true, aryln})
					}
				}
			}
		}
		return true
	})

	funcMap := template.FuncMap{
		"compare":        strings.Compare,
		"contains":       strings.Contains,
		"containsAny":    strings.ContainsAny,
		"containsRune":   strings.ContainsRune,
		"count":          strings.Count,
		"equalFold":      strings.EqualFold,
		"fields":         strings.Fields,
		"fieldsFunc":     strings.FieldsFunc,
		"hasPrefix":      strings.HasPrefix,
		"hasSuffix":      strings.HasSuffix,
		"index":          strings.Index,
		"indexAny":       strings.IndexAny,
		"indexByte":      strings.IndexByte,
		"indexFunc":      strings.IndexFunc,
		"indexRune":      strings.IndexRune,
		"join":           strings.Join,
		"lastIndex":      strings.LastIndex,
		"lastIndexAny":   strings.LastIndexAny,
		"lastIndexByte":  strings.LastIndexByte,
		"lastIndexFunc":  strings.LastIndexFunc,
		"map":            strings.Map,
		"repeat":         strings.Repeat,
		"replace":        strings.Replace,
		"split":          strings.Split,
		"splitAfter":     strings.SplitAfter,
		"splitAfterN":    strings.SplitAfterN,
		"splitN":         strings.SplitN,
		"title":          strings.Title,
		"toLower":        strings.ToLower,
		"toLowerSpecial": strings.ToLowerSpecial,
		"toTitle":        strings.ToTitle,
		"toTitleSpecial": strings.ToTitleSpecial,
		"toUpper":        strings.ToUpper,
		"toUpperSpecial": strings.ToUpperSpecial,
		"trim":           strings.Trim,
		"trimFunc":       strings.TrimFunc,
		"trimLeft":       strings.TrimLeft,
		"trimLeftFunc":   strings.TrimLeftFunc,
		"trimPrefix":     strings.TrimPrefix,
		"trimRight":      strings.TrimRight,
		"trimRightFunc":  strings.TrimRightFunc,
		"trimSpace":      strings.TrimSpace,
		"trimSuffix":     strings.TrimSuffix,
		"noescape":       noescape,
	}

	rndr, err := template.New(filepath.Base(*tmpl)).Funcs(funcMap).ParseFiles(*tmpl)
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

	cmd := exec.Command("go", "fmt", ofn)
	_, err = cmd.Output()
	if err != nil {
		log.Fatal("Go format failed parsing output file\n" + err.Error())
	}

}
