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
	"unicode/utf8"
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
	return walkAst(fmain.Body.List, val)
}

func walkAst(stmts []ast.Stmt, val interface{}) error {
	for _, stmt := range stmts {
		switch st := stmt.(type) {
		case *ast.AssignStmt:

			for ival, lhs := range st.Lhs {
				if lhsv, ok := lhs.(*ast.Ident); ok {

					rhs, ok := st.Rhs[ival].(ast.Expr)
					if !ok {
						return errors.New("only expressions on right side supported")
					}

					rhsv, err := evalExpr(rhs, val)
					if err != nil {
						return err
					}

					fmt.Printf("Присваивание %s = %v\n", lhsv.Name, rhsv)

					if err := setField(val, lhsv.Name, rhsv); err != nil {
						return err
					}

				} else {
					return errors.New("only idents on left side supported")
				}
			}
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
	return nil
}

func setField(val interface{}, fname string, setval interface{}) error {
	tv := reflect.TypeOf(val).Elem()
	for i := 0; i < tv.NumField(); i++ {
		ft := tv.Field(i)
		//only exported fields
		if ft.PkgPath == "" {
			fnm := ft.Name
			if t, ok := ft.Tag.Lookup("goscript"); ok {
				fnm = t
			}
			if strings.EqualFold(fnm, fname) {
				fv := reflect.Indirect(reflect.ValueOf(val)).FieldByName(ft.Name)
				if !fv.CanSet() {
					return errors.New("cant set the field: " + fnm)
				}
				fv.Set(reflect.ValueOf(setval))
				return nil
			}
		}
	}
	return errors.New("field not exists: " + fname)
}

func evalExpr(expr ast.Expr, val interface{}) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.STRING:
			return e.Value, nil
		case token.INT:
			var vv int
			_, err := fmt.Sscan(e.Value, &vv)
			if err != nil {
				return nil, err
			}
			return vv, nil
		case token.FLOAT:
			var vv float64
			_, err := fmt.Sscan(e.Value, &vv)
			if err != nil {
				return nil, err
			}
			return vv, nil
		case token.CHAR:
			vv, _ := utf8.DecodeRuneInString(e.Value)
			if vv == utf8.RuneError {
				return nil, errors.New("incorrect utf8 value")
			}
			return vv, nil

		}
	}
	return nil, errors.New("unrecognized expression")
}
