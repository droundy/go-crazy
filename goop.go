package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/droundy/goopt"
	"go/scanner"
	"go/token"
	"go/parser"
	"go/ast"
	"go/printer"
)

const (
	amBoring = iota
	amAfterDot
	amAfterFunc
	amInReceiver
	amAfterReceiver
	amAfterMethodDot
)

func main() {
	goopt.Parse(func() []string { return []string{} })
	if len(goopt.Args) != 1 {
		fmt.Println("We need the name of a go file to process!")
		os.Exit(1)
	}
	filename := goopt.Args[0]

	bytes,err := ioutil.ReadFile(filename)
	// Make a copy that we're going to convert into parsable go (but not
	// compilable go):
	parseable := make([]byte, len(bytes))
	for i,b := range bytes {
		parseable[i] = b
	}
	if err != nil {
		panic(err)
	}
	var scan scanner.Scanner
	scan.Init(filename, bytes, nil, 0)

	// Variables for scanning
	tok := token.COMMA // anything but EOF or PACKAGE
	var lit []byte
	var pos token.Position

	// Keeping track of where we are:
	where := amBoring
	index := 0
	for tok != token.EOF {
 		pos, tok, lit = scan.Scan()
		//fmt.Println("tok is", tok, "with literal", string(lit), "at", pos)
		switch tok {
		case token.PERIOD:
			//fmt.Printf("%s", bytes[index:pos.Offset])
			index = pos.Offset
			if where == amAfterReceiver {
				where = amAfterMethodDot
			} else {
				where = amAfterDot
			}
			continue
		case token.ADD, token.SUB, token.ADD_ASSIGN:
			switch where {
			case amAfterDot:
				parseable[index] = ' '
				index += 1
				//fmt.Print(" ")
				where = amBoring
			case amAfterMethodDot:
				for _,c := range tok2meth(tok) {
					parseable[index] = uint8(c)
					index++
				}
				continue
			}
		case token.FUNC:
			//fmt.Println("XXXX Am in func")
			where = amAfterFunc
		case token.LPAREN:
			if where == amAfterFunc {
				//fmt.Println("XXXX Am in receiver")
				where = amInReceiver
			} else {
				where = amBoring
			}
		case token.RPAREN:
			if where == amInReceiver {
				//fmt.Println("XXXX Am after receiver")
				where = amAfterReceiver
			} else {
				where = amBoring
			}
		default:
			if where != amInReceiver {
				where = amBoring
			}
		}
		//fmt.Printf("%s%s", bytes[index:pos.Offset], lit)
		index = pos.Offset + len(lit)
 	}
	
	//fmt.Println("Processed and parseable", filename,"...")
	//fmt.Print(string(parseable))

	fileast,err := parser.ParseFile(filename, parseable, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	
	
	//fmt.Println("AST is...")
	newast := Walk(MyVisitor(bytes), fileast)
	printer.Fprint(os.Stdout, newast)
}

type MyVisitor []byte
func (v MyVisitor) Visit(node interface{}) interface{} {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if v[n.OpPos.Offset-1] == '.' {
			//fmt.Println("I found a binary expression!")
			newX := Walk(v, n.X).(ast.Expr)
			newY := Walk(v, n.Y).(ast.Expr)
			return &ast.CallExpr{
				&ast.SelectorExpr{newX, ast.NewIdent(tok2meth(n.Op))},
				n.Pos(),
				[]ast.Expr{newY},
				n.Pos(),
			}
		}
	}
	return nil
}

func tok2meth(t token.Token) string {
	switch t {
	case token.ADD:
		return "P_"
	case token.SUB:
		return "M_"
	case token.ADD_ASSIGN:
		return "PE_"
	}
	panic("Bad token in tok2meth")
}
