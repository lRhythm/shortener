/*
Package main - реализация multichecker.

# Описание

Содержит вызов набора анализаторов, в т.ч. пользовательского.

# Анализаторы

1. golang.org/x/tools/go/analysis/passes: printf, shift, shadow, structtag.

2. honnef.co/go/tools: staticcheck (проверки SA*), stylecheck (проверки ST*).

3. github.com/gostaticanalysis: nilerr, sqlrows.

nilerr проверяет конструкции вида:

	func f() error {
	  err := do()
	  if err != nil {
		return nil // miss
	  }
	}

и

	func f() error {
	  err := do()
	  if err != nil {
		return err // miss
	  }
	}

sqlrows проверяет распространённую ошибку при использовании *sql.Rows - необходимость вызова rows.Close() в defer.

4. exitcheck.

Пользовательский анализатор, запрещающий использовать прямой вызов os.Exit в main:main

# Использование

Запуск из корневой директории проекта:

  go run ./cmd/staticlint/main.go ./...
*/

package main

import (
	"go/ast"

	"github.com/gostaticanalysis/nilerr"
	"github.com/gostaticanalysis/sqlrows/passes/sqlrows"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var checks []*analysis.Analyzer

	// Стандартные статические анализаторы golang.org/x/tools/go/analysis/passes.
	checks = append(checks, printf.Analyzer, shift.Analyzer, shadow.Analyzer, structtag.Analyzer)

	// Все статические анализаторы класса SA пакета staticcheck.io.
	for _, a := range staticcheck.Analyzers {
		checks = append(checks, a.Analyzer)
	}

	// Прочие анализаторы пакета staticcheck.io.
	for _, a := range stylecheck.Analyzers {
		checks = append(checks, a.Analyzer)
	}

	// Публичные анализаторы.
	checks = append(checks, sqlrows.Analyzer, nilerr.Analyzer)

	// Пользовательский анализатор, запрещающий использовать прямой вызов os.Exit в main:main.
	checks = append(checks, ExitAnalyzer)

	multichecker.Main(checks...)
}

// ExitAnalyzer - пользовательский анализатор, запрещающий использовать прямой вызов os.Exit в main:main.
var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "prohibition of call os.Exit() in main:main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	selectorExprFunc := func(x *ast.SelectorExpr) {
		switch t := x.X.(type) {
		case *ast.Ident:
			if t.Name == "os" && x.Sel.Name == "Exit" {
				pass.Reportf(x.Sel.NamePos, "call os.Exit")
			}
		}
	}
	callExprFunc := func(x *ast.CallExpr) {
		switch t := x.Fun.(type) {
		case *ast.SelectorExpr:
			selectorExprFunc(t)
		}
	}
	exprStmtFunc := func(x *ast.ExprStmt) {
		switch t := x.X.(type) {
		case *ast.CallExpr:
			callExprFunc(t)
		}
	}
	blockStmtFunc := func(x *ast.BlockStmt) {
		for _, v := range x.List {
			switch t := v.(type) {
			case *ast.ExprStmt:
				exprStmtFunc(t)
			}
		}
	}
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				switch t := node.(type) {
				case *ast.FuncDecl:
					if t.Name.Name == "main" {
						blockStmtFunc(t.Body)
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
