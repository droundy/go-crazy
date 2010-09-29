package main

import (
	"go/ast"
	"go/token"
	"github.com/droundy/go-crazy/transform"
)

func Inline(fast *ast.File, name string) *ast.File {
	extractor := ExtractFunctionDeclaration{name, nil}
	out := transform.Walk(&extractor, fast).(*ast.File)
	inliner := InlineFunction(extractor)
	return transform.Walk(&inliner, out).(*ast.File)
}

type InlineFunction ExtractFunctionDeclaration

func (v *InlineFunction) Visit(node interface{}) interface{} {
	if n,ok := node.(*ast.CallExpr); ok {
		if id,ok := n.Fun.(*ast.Ident); ok && id.Name == v.Name {
			return &ast.ParenExpr {
				n.Pos(),
				&ast.CallExpr{
					&ast.FuncLit{
						v.ItsDecl.Type,
						v.ItsDecl.Body,
					},
					n.Lparen,
					n.Args,
					n.Ellipsis,
					n.Rparen,
				},
				n.Rparen,
			}
		}
	}
	return nil
}

type ExtractFunctionDeclaration struct {
	Name string
	ItsDecl *ast.FuncDecl
}

func (v *ExtractFunctionDeclaration) Visit(node interface{}) interface{} {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Name.Name == v.Name {
			v.ItsDecl = n
			var nopos token.Position
			return &ast.GenDecl{
				n.Doc,
				n.Pos(),
				token.CONST,
				nopos,
				[]ast.Spec{&ast.ValueSpec{
						nil,
						[]*ast.Ident{ast.NewIdent("i_inlined_"+v.Name)},
						nil,
						[]ast.Expr{&ast.BasicLit{nopos, token.INT, []byte("0")}},
						nil,
					},
				},
				nopos,
			}
		}
	}
	return nil
}
