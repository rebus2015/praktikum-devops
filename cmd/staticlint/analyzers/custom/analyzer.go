// Package custom реализует собственный анализатор кода
package custom

import (
	"flag"
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name:             "exitcheck",
	Doc:              "check for os.exit usage in main.go file",
	URL:              "",
	Flags:            flag.FlagSet{},
	Run:              run,
	RunDespiteErrors: false,
	Requires:         []*analysis.Analyzer{},
	ResultType:       nil,
	FactTypes:        []analysis.Fact{},
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					ast.Inspect(x, func(n ast.Node) bool {
						switch f := n.(type) {
						case *ast.CallExpr:
							fun, ok := f.Fun.(*ast.SelectorExpr)
							if ok && fmt.Sprintf("%s.%s", fun.X, fun.Sel.Name) == "os.Exit" {
								pass.Reportf(f.Pos(), "os.Exit usage detected in function main, file main")
							}
						}
						return true
					})
					return false
				}
			}
			return true
		},
		)
	}
	return nil, nil
}
