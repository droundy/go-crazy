package main

import (
	"fmt"
	"os"
	"exec"
	"github.com/droundy/goopt"
	"github.com/droundy/goop/parser"
	"go/printer"
	//"github.com/droundy/goop/transform"
)

func panicon(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func archnum() string {
	switch os.Getenv("GOARCH") {
	case "386": return "8"
	case "amd64": return "6"
		// what was the other one called?
	}
	return "5"
}

func main() {
	goopt.Parse(func() []string { return []string{} })
	if len(goopt.Args) != 1 {
		fmt.Println("We need the name of a go file to process!")
		os.Exit(1)
	}
	filename := goopt.Args[0]

	fileast,err := parser.ParseFile(filename, nil, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Parse error:\n", err)
		os.Exit(1)
	}

	// Let's create a file containing the parsed code...
	basename := filename[0:len(filename)-3]
	newfilename := basename+"-compiled.go"
	out,err := os.Open(newfilename, os.O_WRONLY + os.O_TRUNC + os.O_CREAT, 0644)
	panicon(err)
	panicon(printer.Fprint(out, fileast))
	out.Close()

	objname := basename+"-compiled."+archnum()
	if e := justrun(archnum()+"g", "-o", objname, newfilename); e != nil {
		fmt.Println("Error compiling", newfilename,"!")
		fmt.Println(e)
		os.Exit(1)
	}
	panicon(justrun(archnum()+"l", "-o", basename, objname))
}

func justrun(cmd string, args ...string) os.Error {
	abscmd,err := exec.LookPath(cmd)
	if err != nil { return os.NewError("Couldn't find "+cmd+": "+err.String()) }
	
	cmdargs := make([]string, len(args)+1)
	cmdargs[0] = cmd
	for i,a := range args {
		cmdargs[i+1] = a
	}
	pid, err := exec.Run(abscmd, cmdargs, nil, "",
		exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if err != nil { return err }
	wmsg,err := pid.Wait(0)
	if err != nil { return err }
	if wmsg.ExitStatus() != 0 {
		return os.NewError(cmd+" exited with status "+fmt.Sprint(wmsg.ExitStatus()))
	}
	return nil
}

/*
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
		} else if v[n.OpPos.Offset] == '*' && v[n.OpPos.Offset+1] == '.' {
			newX := transform.Walk(v, n.X).(ast.Expr)
			newY := transform.Walk(v, n.Y).(ast.Expr)
			return &ast.CallExpr{
				&ast.SelectorExpr{newY, ast.NewIdent("_mul_dot")},
				n.Pos(),
				[]ast.Expr{newX},
				n.Pos(),
			}
		}
	}
	return nil
}
*/
