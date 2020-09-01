package errchecker

import (
	"go/ast"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "errchecker identify the functions that return error and whose return value of error is all nil"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "errchecker",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	inspect.Preorder(nodeFilter, func(decl ast.Node) {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			idx, err := errorCheker(decl)
			if err != nil {
				log.Println(err)
			}
			// 追加
			// 本文decl.Body.Listからreturn文がある一行を取得している
			for _, stmt := range decl.Body.List {
				// return文を取得している
				ret, _ := stmt.(*ast.ReturnStmt)
				if ret == nil {
					continue
				}
				// Resultsが0の時第一引数
				isReturnNil := true
				switch lit := ret.Results[idx].(type) {
				case *ast.BasicLit:
					log.Println(lit.Kind)
				case *ast.Ident:
					log.Printf("error:%s", lit.Name)
					if lit.Name != "nil" {
						isReturnNil = false
					}
				}
				if isReturnNil {
					pass.Reportf(stmt.Pos(), "It returns nil in all the places where it should return error:%d", stmt.Pos())
				}
			}
		}
	})
	return nil, nil
}

// errorChecker returns the index of error in the return value
func errorCheker(n *ast.FuncDecl) (int, error) {
	// 受け取った関数定義の引数リスト
	fieldList := n.Type.Results.List
	var isErrorExist bool
	index := -1
	for idx, t := range fieldList {
		switch ty := t.Type.(type) {
		case *ast.Ident:
			if ty.Name == "error" {
				isErrorExist = true
				index = idx
			}
		}
		// 返り値の型が表示される
		// log.Println(t.Type)
	}
	if isErrorExist {
		return index, nil
	}
	// errorを返り値として持たなかったら
	return index, nil
}
