package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

// Description ....
type Description struct {
	Name        string
	Description string
}

type Descriptions []Description

func (a Descriptions) Len() int           { return len(a) }
func (a Descriptions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Descriptions) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) < 0 }

func extractDescriptions(fileName string) Descriptions {
	result := []Description{}
	fset := token.NewFileSet()

	fileAst, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, decl := range fileAst.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name != nil {
				add := Description{
					Name: fn.Name.Name,
				}
				if fn.Doc.Text != nil {
					add.Description = fn.Doc.Text()
					/*
						lines := strings.Split(fn.Doc.Text(), "\n")
						for _, line := range lines {
							add.Description += line
						}
					*/
				}
				result = append(result, add)
			}
		}
	}

	return result
}

func main() {
	res := Descriptions{}

	for _, fn := range os.Args[1:] {
		res = append(res, extractDescriptions(fn)[:]...)
	}

	sort.Sort(res)
	for _, r := range res {
		fmt.Printf("## %s\n\n%s\n\n", r.Name, r.Description)
	}
}
