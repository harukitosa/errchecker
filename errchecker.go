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

func search(body []ast.Stmt, idx int) bool {
	if idx == -1 {
		return false
	}
	flag := true
	// fmt.Println("--------------------------------")
	for _, stmt := range body {
		// fmt.Printf("%+v\n", stmt)
		switch let := stmt.(type) {
		case *ast.ReturnStmt:
			if len(let.Results) <= idx {
				return false
			}
			if !isNil(let.Results[idx]) {
				return false
			}
			// fmt.Println("return stmt true")
			return true
		case *ast.IfStmt:
			flag = search(let.Body.List, idx)
			if let.Else != nil {
				block, ok := let.Else.(*ast.BlockStmt)
				if !ok {
					continue
				}
				flag = search(block.List, idx)
			}
			if !flag {
				return flag
			}
		case *ast.ForStmt:
			flag = search(let.Body.List, idx)
			if !flag {
				return flag
			}
		case *ast.SwitchStmt:
			flag = search(let.Body.List, idx)
			if !flag {
				return flag
			}
		}
	}
	// fmt.Printf("flag:%t\n", flag)
	// fmt.Println("--------------------------------")
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

func errAllNilAnon(lit *ast.FuncLit, pass *analysis.Pass) bool {
	idx := returnErrIndex(lit.Type, pass)
	if idx == -1 {
		return false
	}
	return search(lit.Body.List, idx)
}

// Check if all places that return error return nil
func errAllNil(decl *ast.FuncDecl, pass *analysis.Pass) bool {
	idx := returnErrIndex(decl.Type, pass)
	if idx == -1 {
		return false
	}
	flag := search(decl.Body.List, idx)
	return flag
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
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			if errAllNil(decl, pass) {
				pass.Reportf(decl.Pos(), "It returns nil in all the places where it should return error %d", decl.Pos())
				return
			}
		case *ast.FuncLit:
			if errAllNilAnon(decl, pass) {
				pass.Reportf(decl.Pos(), "It returns nil in all the places where it should return error %d", decl.Pos())
				return
			}
		}
	})
	return nil, nil
}
