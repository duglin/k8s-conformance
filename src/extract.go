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

// Suite ...
type Suite struct {
	TOC   sort.StringSlice
	Tests Tests
}

// Test ...
type Test struct {
	FileName    string
	Line        int
	Name        string
	Description string
}

type Tests []Test

func (a Tests) Len() int      { return len(a) }
func (a Tests) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Tests) Less(i, j int) bool {
	return strings.Compare(a[i].Name, a[j].Name) < 0
}

func AddTests(suite *Suite, fileName string) *Suite {
	fset := token.NewFileSet()

	fileAst, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, decl := range fileAst.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name != nil {
				add := Test{
					FileName: fileName,
					Name:     fn.Name.Name,
				}

				if fn.Doc != nil && fn.Doc.List != nil {
					add.Line = fset.Position(fn.Doc.List[0].Slash).Line
				} else {
					add.Line = fset.Position(fn.Type.Func).Line
				}

				if i := strings.IndexAny(fn.Name.Name, "0123456789"); i > 0 {
					toc := fn.Name.Name[0:i]

					i = suite.TOC.Search(toc)
					if i == suite.TOC.Len() || suite.TOC[i] != toc {
						suite.TOC = append(suite.TOC[0:i],
							append([]string{toc}, suite.TOC[i:]...)...)
					}
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
				suite.Tests = append(suite.Tests, add)
			}
		}
	}

	return suite
}

func main() {
	suite := Suite{}

	var srcFn = flag.String("src", "src/funcs.go", "File to write src file")
	var docFn = flag.String("doc", "tests.md", "File to write doc file")
	flag.Parse()

	for _, fn := range flag.Args() {
		AddTests(&suite, fn)
	}

	sort.Sort(suite.Tests)

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

	docFile.WriteString("## Table of Contents\n\n")

	for i, t := range suite.TOC {
		docFile.WriteString(fmt.Sprintf("%d. [%s](#%s)\n", 1+i, t, strings.ToLower(t)))
	}

	docFile.WriteString("\n")

	srcFile.WriteString("package main\n\n")
	srcFile.WriteString("import \"../utils\"\n")
	srcFile.WriteString("import \"../tests\"\n\n")
	srcFile.WriteString("var TestMap = map[string]func(*utils.Test){\n")

	prevTOC := -1

	for _, r := range suite.Tests {
		if prevTOC == -1 || !strings.HasPrefix(r.Name, suite.TOC[prevTOC]) {
			prevTOC = prevTOC + 1
			docFile.WriteString(fmt.Sprintf("## %s\n\n", suite.TOC[prevTOC]))
		}

		fmt.Fprintf(docFile, "### [%s](%s#L%d)\n\n%s\n\n", r.Name, r.FileName,
			r.Line, r.Description)

		str := fmt.Sprintf("\t\"%s\": tests.%s,\n", r.Name, r.Name)
		srcFile.WriteString(str)
	}

	srcFile.WriteString("}\n\nvar TestNames = []string{\n")
	for _, r := range suite.Tests {
		srcFile.WriteString("\t\"" + r.Name + "\",\n")
	}
	srcFile.WriteString("}\n")
}
