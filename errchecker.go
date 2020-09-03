package errchecker

import (
	"errors"
	"go/ast"

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

func isNil(exp ast.Expr) bool {
	switch lit := exp.(type) {
	case *ast.Ident:
		if lit.Name != "nil" {
			return false
		}
	default:
		return false
	}
	return true
}

// returnIdxNil checks if all nils are returned at the specified index
func returnIdxNil(body []ast.Stmt, idx int) bool {
	if idx == -1 {
		return false
	}
	flag := true
	for _, stmt := range body {
		switch let := stmt.(type) {
		case *ast.ReturnStmt:
			if len(let.Results) <= idx {
				return false
			}
			if !isNil(let.Results[idx]) {
				return false
			}
			return true
		case *ast.IfStmt:
			flag = returnIdxNil(let.Body.List, idx)
			if let.Else != nil {
				block, ok := let.Else.(*ast.BlockStmt)
				if !ok {
					continue
				}
				flag = returnIdxNil(block.List, idx)
			}
			if !flag {
				return flag
			}
		case *ast.ForStmt:
			flag = returnIdxNil(let.Body.List, idx)
			if !flag {
				return flag
			}
		case *ast.SwitchStmt:
			flag = returnIdxNil(let.Body.List, idx)
			if !flag {
				return flag
			}
		}
	}
	return flag
}

// errorChecker returns the index of error in the return value.
// If there is no error in the return value, it returns -1
func returnErrIndex(n *ast.FuncType, pass *analysis.Pass) int {
	index := -1
	results := n.Results
	if results == nil {
		return -1
	}
	fieldList := results.List
	if fieldList == nil {
		return -1
	}
	for idx, t := range fieldList {
		switch ty := t.Type.(type) {
		case *ast.Ident:
			s := pass.TypesInfo.Types[ty]
			if analysisutil.ImplementsError(s.Type) {
				index = idx
			}
		}
	}
	return index
}

// Check if all places that return error return nil
func errAllNil(node ast.Node, pass *analysis.Pass) bool {
	switch n := node.(type) {
	case *ast.FuncLit:
		idx := returnErrIndex(n.Type, pass)
		if idx == -1 {
			return false
		}
		flag := returnIdxNil(n.Body.List, idx)
		return flag
	case *ast.FuncDecl:
		idx := returnErrIndex(n.Type, pass)
		if idx == -1 {
			return false
		}
		flag := returnIdxNil(n.Body.List, idx)
		return flag
	}
	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	a, ok := pass.ResultOf[inspect.Analyzer]
	if !ok {
		return nil, errors.New("*inspector.Inspector assertion error")
	}
	inspect, ok := a.(*inspector.Inspector)
	if !ok {
		return nil, errors.New("*inspector.Inspector assertion error")
	}
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}
	inspect.Preorder(nodeFilter, func(decl ast.Node) {
		if errAllNil(decl, pass) {
			pass.Reportf(decl.Pos(), "It returns nil in all the places where it should return error %d", decl.Pos())
			return
		}

	})
	return nil, nil
}
