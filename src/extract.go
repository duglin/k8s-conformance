package main

import (
	"flag"
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

	var srcFn = flag.String("src", "src/funcs.go", "File to write src file")
	var docFn = flag.String("doc", "tests.md", "File to write doc file")
	flag.Parse()

	for _, fn := range flag.Args() {
		res = append(res, extractDescriptions(fn)[:]...)
	}

	sort.Sort(res)

	srcFile, err := os.Create(*srcFn)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	defer srcFile.Close()

	docFile, err := os.Create(*docFn)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	defer docFile.Close()

	srcFile.WriteString("package main\n\n")
	srcFile.WriteString("import \"../utils\"\n")
	srcFile.WriteString("import \"../tests\"\n\n")
	srcFile.WriteString("var TestMap = map[string]func(*utils.Test){\n")

	for _, r := range res {
		fmt.Fprintf(docFile, "## %s\n\n%s\n\n", r.Name, r.Description)

		str := fmt.Sprintf("\t\"%s\": tests.%s,\n", r.Name, r.Name)
		srcFile.WriteString(str)
	}
	srcFile.WriteString("}\n\nvar TestNames = []string{\n")
	for _, r := range res {
		srcFile.WriteString("\t\"" + r.Name + "\",\n")
	}
	srcFile.WriteString("}\n")
}
