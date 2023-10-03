// Package custom реализует собственный анализатор кода.
package custom

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.exit usage in main.go file",
	Run:  run,
}

const name string = "main"

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != name {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			x, ok := node.(*ast.FuncDecl)
			if ok {
				if x.Name.Name == name {
					ast.Inspect(x, func(n ast.Node) bool {
						f, ok := n.(*ast.CallExpr)
						if ok {
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
	return true, nil
}
