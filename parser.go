package main

import (
	"fmt"
	"go/ast"
	"regexp"
)

type StructDecl struct {
	Name   string
	Fields []Field
	Ref    *ast.StructType
}

type StructField struct {
	Name string
	Type string
	Tag  string
	Ref  *ast.Field
}

type JsonTag struct {
	Name      string
	Omitempty bool
}

type Field struct {
	Type string
	Name string
}

func (decl *StructDecl) appendField(field Field) {
	decl.Fields = append(decl.Fields, field)
}

func inspectStructDecls(node *ast.File) []StructDecl {
	var decls []StructDecl

	ast.Inspect(node, func(node ast.Node) bool {
		spec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structDecl, ok := spec.Type.(*ast.StructType)
		if !ok {
			return false
		}

		decls = append(decls, StructDecl{
			Name:   spec.Name.Name,
			Fields: nil,
			Ref:    structDecl,
		})

		return false
	})

	return decls
}

func inspectStructFields(node *ast.StructType) []StructField {
	var fields []StructField

	ast.Inspect(node, func(node ast.Node) bool {
		field, ok := node.(*ast.Field)
		if !ok {
			return true
		}

		var t string
		switch fieldType := field.Type.(type) {
		case *ast.Ident:
			t = fieldType.Name
		case *ast.StarExpr:
			switch refType := fieldType.X.(type) {
			case *ast.Ident:
				t = refType.Name
			default:
				return false
			}
		case *ast.ArrayType:
			// How supports Array of identifier?
			return false
		default:
			return false
		}

		var name string
		if field.Names != nil && field.Names[0] != nil {
			name = field.Names[0].Name
		} else {
			// Skip incomplete name
			return false
		}

		var tag string
		if field.Tag != nil {
			tag = field.Tag.Value
		} else {
			// Skip untagged field
			return false
		}

		fields = append(fields, StructField{
			Name: name,
			Type: t,
			Tag:  tag,
			Ref:  field,
		})

		return false
	})

	return fields
}

func parseJsonTag(tag string) (*JsonTag, error) {
	jsonTagRegExp := regexp.MustCompile(`json:"(\w+|\-)(,omitempty)?"`)
	matches := jsonTagRegExp.FindStringSubmatch(tag)

	if len(matches) == 0 {
		return nil, fmt.Errorf("%s is not JSON tag", tag)
	}

	return &JsonTag{
		Name:      matches[1],
		Omitempty: len(matches) == 3,
	}, nil
}
