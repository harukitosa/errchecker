package errchecker

import (
	"go/ast"
	"log"

	"github.com/gostaticanalysis/analysisutil"
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

// errorChecker returns the index of error in the return value
// todo: 複数のエラーがあった場合の対処
func errorCheker(n *ast.FuncDecl, pass *analysis.Pass) (int, error) {
	// 受け取った関数定義の引数リスト
	fieldList := n.Type.Results.List
	var isErrorExist bool
	index := -1
	for idx, t := range fieldList {
		switch ty := t.Type.(type) {
		case *ast.Ident:
			// Question:ここを厳密に型比較で行う場合はどうしたらいいのか？
			s := pass.TypesInfo.Types[ty]
			if analysisutil.ImplementsError(s.Type) {
				isErrorExist = true
				index = idx
			}
		}
	}
	if isErrorExist {
		return index, nil
	}
	// errorを返り値として持たなかったら
	return index, nil
}

func search(body []ast.Stmt, idx int) (bool, error) {
	isReturnNil := true
	var err error
	for _, stmt := range body {
		switch let := stmt.(type) {
		case *ast.ReturnStmt:
			switch lit := let.Results[idx].(type) {
			// nilの場合は*ast.Indentにふり分けられる
			case *ast.Ident:
				if lit.Name != "nil" {
					return false, nil
				}
			default:
				return false, nil
			}
		case *ast.IfStmt:
			isReturnNil, err = search(let.Body.List, idx)
			if err != nil {
				return isReturnNil, err
			}
		case *ast.ForStmt:
			isReturnNil, err = search(let.Body.List, idx)
			if err != nil {
				return isReturnNil, err
			}
		}
	}
	return isReturnNil, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	a, ok := pass.ResultOf[inspect.Analyzer]
	if !ok {
		log.Println("*inspector.Inspector assertion error")
	}
	inspect, ok := a.(*inspector.Inspector)
	if !ok {
		log.Println("*inspector.Inspector assertion error")
	}
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	inspect.Preorder(nodeFilter, func(decl ast.Node) {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			// idxはエラーが現れる場所の数値
			idx, err := errorCheker(decl, pass)
			if err != nil {
				log.Println(err)
			}
			flag, err := search(decl.Body.List, idx)
			if err != nil {
				log.Println(err)
			}
			if flag {
				pass.Reportf(decl.Pos(), "It returns nil in all the places where it should return error %d", decl.Pos())
			}
		}
	})
	return nil, nil
}
