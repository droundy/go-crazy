package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/droundy/goopt"
	"github.com/droundy/goop/parser"
	"go/ast"
	"go/printer"
	"github.com/droundy/goop/transform"
)

func main() {
	goopt.Parse(func() []string { return []string{} })
	if len(goopt.Args) != 1 {
		fmt.Println("We need the name of a go file to process!")
		os.Exit(1)
	}
	filename := goopt.Args[0]

	bytes,err := ioutil.ReadFile(filename)

	fileast,err := parser.ParseFile(filename, bytes, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	
	
	//fmt.Println("AST is...")
	newast := transform.Walk(MyVisitor(bytes), fileast)
	printer.Fprint(os.Stdout, newast)
}

type MyVisitor []byte
func (v MyVisitor) Visit(node interface{}) interface{} {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if v[n.OpPos.Offset] == '.' {
			//fmt.Println("I found a binary expression!")
			newX := transform.Walk(v, n.X).(ast.Expr)
			newY := transform.Walk(v, n.Y).(ast.Expr)
			return &ast.CallExpr{
				&ast.SelectorExpr{newX, ast.NewIdent(parser.MungeOperator(n.Op))},
				n.Pos(),
				[]ast.Expr{newY},
				n.Pos(),
			}
		}
	}
	return nil
}
