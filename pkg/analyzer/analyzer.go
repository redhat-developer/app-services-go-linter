package analyzer

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "goi18nlinter",
	Doc:  "goi18nlinter",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.CallExpr:
			s, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if s.Sel.Name == "MustLocalize" || s.Sel.Name == "MustLocalizeError" {
				fmt.Println(s.Sel.Name)
				if len(n.Args) > 0 {
					args, ok := n.Args[0].(*ast.BasicLit)
					if !ok {
						return true
					}
					fmt.Println(args.Value)
				}
			}
		}
		return true
	}

	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}

	return nil, nil
}
