package analyzer

import (
	"github.com/alexal/go-i18n-linter/pkg/localize"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"log"
	"strings"
	"sync"
)

var (
	Analyzer = &analysis.Analyzer{
		Name:             "goi18nlinter",
		Doc:              "goi18nlinter",
		Run:              new(instantiate).run,
		Requires:         []*analysis.Analyzer{inspect.Analyzer},
		RunDespiteErrors: true,
	}

	path              string
	mustLocalize      string
	mustLocalizeError string
)

type instantiate struct {
	once     sync.Once
	messages map[string]string
}

func (i *instantiate) run(pass *analysis.Pass) (interface{}, error) {
	task := func() {
		i.messages = make(map[string]string)
		localizer, err := localize.New(nil, path)
		if err != nil {
			log.Panicln(err)
		}

		for _, file := range localizer.GetTranslations() {
			for _, msg := range file.Messages {
				i.messages[msg.ID] = msg.One
			}
		}
	}
	runOnce(&i.once, task)

	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.SelectorExpr)(nil),
		(*ast.BasicLit)(nil),
	}

	ins.Preorder(nodeFilter, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.CallExpr:
			s, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			if s.Sel.Name == mustLocalize || s.Sel.Name == mustLocalizeError {
				if len(n.Args) > 0 {
					args, ok := n.Args[0].(*ast.BasicLit)
					if !ok {
						return
					}
					str := strings.Trim(args.Value, "\"")
					if i.messages[str] == "" {
						pass.Reportf(args.Pos(), "Translation string '%s' doesn't exist", str)
					}
				}
			}
		}
	})

	return nil, nil
}

func runOnce(once *sync.Once, onceBody func()) {
	once.Do(onceBody)
}

func init() {
	Analyzer.Flags.StringVar(&path, "path", "", "Path to the directory with localization files")
	Analyzer.Flags.StringVar(&mustLocalize, "mustLocalize", "MustLocalize", "")
	Analyzer.Flags.StringVar(&mustLocalizeError, "mustLocalizeError", "MustLocalizeError", "")
}
