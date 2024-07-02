package main

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// Config — имя файла конфигурации.
const Config = `config.json`

// ConfigData описывает структуру файла конфигурации.
type ConfigData struct {
	Staticcheck []string
}

func main() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}
	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}

func osExitLint() bool {
	src := `package main
import "fmt"

func calc(a int) {
   b := (a+8)*2
   return b
}

func main() {
    fmt.Println(calc(3), calc(5))
    fmt.Println("Hello, world!")
	os.Exit(0)
}`

	// дерево разбора AST ассоциируется с набором исходных файлов FileSet
	fset := token.NewFileSet()
	// парсер может работать с файлом
	// или исходным кодом, переданным в виде строки
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		// проверяем, какой конкретный тип лежит в узле
		switch x := n.(type) {
		case *ast.Ident:
			if x.Name == "osExit" {
				return true
			}
		}
		return false
	})

	return false
}
