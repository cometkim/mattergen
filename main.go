package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
	"text/template"
)

type StructDecl struct {
	Name string
	Ref  *ast.StructType
}

type StructField struct {
	Type string
	Name string
	Ref  *ast.Field
}

type Spec struct {
	Type   string
	Name   string
	Fields []Spec
}

func main() {
	fset := token.NewFileSet()

	// TODO: Read all model from mattermost-server
	gopath := os.Getenv("GOPATH")
	pkgpath := "github.com/mattermost/mattermost-server"
	modelpath := path.Join(gopath, "src", pkgpath, "model")

	f, err := parser.ParseFile(fset, path.Join(modelpath, "user.go"), nil, 0)
	if err != nil {
		panic(err)
	}

	// TODO: Map Field.Type to ts, flow, graphql ...and whatever
	tsMap := map[string]string{
		"bool":      "boolean",
		"int":       "number",
		"int64":     "number",
		"string":    "string",
		"StringMap": "Map<string, string>",
	}

	var specs []Spec
	decls := parseStructDecls(f)
	for _, decl := range decls {
		var root Spec
		var children []Spec
		for _, field := range parseStructJsonFields(decl.Ref) {
			children = append(children, Spec{
				Name:   field.Name,
				Type:   tsMap[field.Type],
				Fields: nil,
			})
		}
		root.Name = decl.Name
		root.Type = "interface"
		root.Fields = children
		specs = append(specs, root)
	}

	t, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(os.Stdout, "typescript.tmpl", specs)
	if err != nil {
		panic(err)
	}
}

func parseStructDecls(f *ast.File) []StructDecl {
	var decls []StructDecl
	for _, fdecl := range f.Decls {
		if gen, ok := fdecl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				if s, ok := spec.(*ast.TypeSpec); ok {
					if t, ok := s.Type.(*ast.StructType); ok {
						decls = append(decls, StructDecl{
							Name: s.Name.Name,
							Ref:  t,
						})
					}
				}
			}
		}
	}
	return decls
}

func parseStructJsonFields(s *ast.StructType) []StructField {
	var fields []StructField
	for _, f := range s.Fields.List {
		if i := strings.Index(f.Tag.Value, "json"); i != -1 {
			n := strings.Split(f.Tag.Value, "json:\"")[1]
			n = strings.Split(n, "\"")[0]
			n = strings.Split(n, ",")[0]

			var t string
			if ft, ok := f.Type.(*ast.Ident); ok {
				t = ft.Name
			} else if ft, ok := f.Type.(*ast.StarExpr); ok {
				t = ft.X.(*ast.Ident).Name
			}

			fields = append(fields, StructField{
				Type: t,
				Name: n,
				Ref:  f,
			})
		}
	}
	return fields
}
