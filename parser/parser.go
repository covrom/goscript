package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func ParseFile(fn string) (*ast.File, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	pckMain := strings.NewReader("package main\nfunc main(){\n")
	pckMainEnd := strings.NewReader("}")
	astf, err := parser.ParseFile(fset, filepath.Base(fn), io.MultiReader(pckMain, f, pckMainEnd), 0)
	return astf, err
}

func RunAST(astf *ast.File, val interface{}) error {
	if len(astf.Decls) != 1 {
		return errors.New("only one function in global scope can be defined")
	}
	if reflect.TypeOf(val).Kind() != reflect.Ptr {
		return errors.New("not a pointer to struct")
	}
	if reflect.ValueOf(val).Elem().Kind() != reflect.Struct {
		return errors.New("not a pointer to struct")
	}
	fmain := astf.Decls[0].(*ast.FuncDecl)
	walkAst(fmain.Body.List)
	return nil
}

func walkAst(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		switch st := stmt.(type) {
		case *ast.AssignStmt:
			fmt.Printf("Присваивание %v\n", st)
		case *ast.ExprStmt:
			switch ex := st.X.(type) {
			case *ast.CallExpr:
				fmt.Printf("Вызов функции %v\n", ex)
			default:
				fmt.Printf("%#v\n", ex)
			}
		default:
			fmt.Printf("%#v\n", st)
		}
	}
}
