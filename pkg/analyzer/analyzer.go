package analyzer

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:             "goi18nlinter",
	Doc:              "goi18nlinter",
	Run:              run,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	RunDespiteErrors: true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.SelectorExpr)(nil),
		(*ast.BasicLit)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.CallExpr:
			s, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			if s.Sel.Name == "MustLocalize" || s.Sel.Name == "MustLocalizeError" {
				fmt.Println(s.Sel.Name)
				pass.Reportf(1, s.Sel.Name)
				if len(n.Args) > 0 {
					args, ok := n.Args[0].(*ast.BasicLit)
					if !ok {
						return
					}
					fmt.Println(args.Value)
				}
			}
		}
	})

	return nil, nil
}
