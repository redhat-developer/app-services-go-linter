package analyzer

import (
	"github.com/alexal/go-i18n-linter/pkg/localize"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"log"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goi18nlinter",
	Doc:      "goi18nlinter",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var messages map[string]string

func run(pass *analysis.Pass) (interface{}, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.SelectorExpr)(nil),
		(*ast.BasicLit)(nil),
	}

	if messages == nil {
		messages = make(map[string]string)
		localizer, err := localize.New(nil)
		if err != nil {
			log.Panicln(err)
		}
		for _, file := range localizer.GetTranslations() {
			for _, msg := range file.Messages {
				messages[msg.ID] = msg.One
			}
		}
	}

	ins.Preorder(nodeFilter, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.CallExpr:
			s, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			if s.Sel.Name == "MustLocalize" || s.Sel.Name == "MustLocalizeError" {
				if len(n.Args) > 0 {
					args, ok := n.Args[0].(*ast.BasicLit)
					if !ok {
						return
					}
					str := strings.Trim(args.Value, "\"")
					if messages[str] == "" {
						pass.Reportf(args.Pos(), "Translation string '%s' doesn't exist", str)
					}
				}
			}
		}
	})

	return nil, nil
}
