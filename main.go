package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	var targetId string

	flag.StringVar(&targetId, "target", "", "Target want to use auto-generated type definitions (Required)")
	flag.Parse()

	if targetId == "" {
		exitWithError("--target must be specified. currently supports typescript or flow")
	}

	target, err := NewTarget(targetId)
	if err != nil {
		exitWithError(err)
	}

	modelpath, err := checkPackageDir()
	if err != nil {
		exitWithError(err)
	}

	fset := token.NewFileSet()

	pkg, err := parser.ParseDir(fset, modelpath, func(file os.FileInfo) bool {
		isNotTest := !strings.HasSuffix(file.Name(), "_test.go")
		return isNotTest
	}, 0)
	if err != nil {
		exitWithError(err)
	}

	tmplBase := targetId + ".tmpl"
	tmpl := template.Must(template.New(tmplBase).ParseFiles(path.Join("templates/", tmplBase)))

	for filePath, fileNode := range pkg["model"].Files {
		fmt.Printf("Inspecting %s\n", filePath)

		decls := inspectStructDecls(fileNode)
		var newDecls []StructDecl

		for _, decl := range decls {
			fields := inspectStructFields(decl.Ref)
			var newFields []Field

			for _, field := range fields {
				if field.Tag == "" {
					continue
				}

				tag, err := parseJsonTag(field.Tag)
				if err != nil {
					continue
				}

				if tag.Name == "-" {
					continue
				}

				t := target.convertType(field.Type)
				if t == "" {
					continue
				}

				newFields = append(newFields, Field{
					Type: t,
					Name: tag.Name,
				})
			}

			if len(newFields) != 0 {
				decl.Fields = newFields
				newDecls = append(newDecls, decl)
			}
		}

		if len(newDecls) != 0 {
			for _, decl := range newDecls {
				fmt.Printf("Decl: %s\n", decl.Name)
				for _, field := range decl.Fields {
					fmt.Printf("\tField { Type:%s, Name:%s }\n", field.Type, field.Name)
				}
			}

			fileBase := filepath.Base(filePath)
			fileExt := filepath.Ext(fileBase)
			filename := fileBase[:len(fileBase)-len(fileExt)]

			outpath := path.Join("output", filename+target.Ext)

			f, err := os.OpenFile(
				outpath,
				os.O_CREATE|os.O_WRONLY,
				0664,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer f.Close()

			err = tmpl.Execute(f, newDecls)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%s is generated\n", outpath)

		} else {
			fmt.Println("No available declarations, skipped.")
		}
	}
}

func checkPackageDir() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}

	pkgpath := path.Join(gopath, "src", "github.com", "mattermost", "mattermost-server")
	if _, err := os.Stat(pkgpath); os.IsNotExist(err) {
		return "", fmt.Errorf("Can't find mattermost/mattermost-server in your GOPATH")
	}

	return path.Join(pkgpath, "model"), nil
}

func exitWithError(message ...interface{}) {
	fmt.Println(message)
	os.Exit(1)
}
