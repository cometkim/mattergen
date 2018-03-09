package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func main() {
	fset := token.NewFileSet()

	// TODO: Read all model from mattermost-server
	f, err := parser.ParseFile(fset, "sample/user.go", nil, 0)
	if err != nil {
		panic(err)
	}

	decls := parseStructDecls(f)
	for _, decl := range decls {
		fmt.Println(decl.Name)
		for _, field := range parseStructJsonFields(decl.Pos) {
			fmt.Printf(" - %s %s\n", field.Name, field.Type)

			// TODO: Map Field.Type to ts, flow, and graphql
		}
	}

	// TODO: Generate files using template
}

type StructDecl struct {
	Name string
	Pos  *ast.StructType
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
							Pos:  t,
						})
					}
				}
			}
		}
	}
	return decls
}

type StructField struct {
	Type string
	Name string
	Pos  *ast.Field
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
				Pos:  f,
			})
		}
	}
	return fields
}
