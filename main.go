package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "sample/user.go", nil, 0)
	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			if gen.Tok == token.TYPE {
				for _, spec := range gen.Specs {
					if t, ok := spec.(*ast.TypeSpec); ok {
						// ast.Print(fset, t)
						if s, ok := t.Type.(*ast.StructType); ok {
							fmt.Println(t.Name)
							for _, f := range s.Fields.List {
								fmt.Print("  - ")
								if ft, ok := f.Type.(*ast.Ident); ok {
									fmt.Print(ft.Name)
								} else if ft, ok := f.Type.(*ast.StarExpr); ok {
									fmt.Print(ft.X.(*ast.Ident).Name)
								}
								fmt.Print(" ", f.Tag.Value)
								fmt.Println()
							}
							fmt.Println()
						}
					}
				}
			}
		}
	}
}
