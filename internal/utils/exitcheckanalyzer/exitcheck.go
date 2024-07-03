package exitcheckanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ExitCheck finds use osExit with exit code and return error
var ExitCheck = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check os.Exit exit code for main func and main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Ident:
				if x.Name == "Exit" {
					pass.Reportf(x.Pos(), "can't use osExit with exit code")
				}
			}

			return true
		})
	}
	return nil, nil
}
