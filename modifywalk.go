// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
)

// A Visitor's Visit method is invoked for each node encountered by
// Walk.  If the result modifiednode is not nil, walk replaces node
// with modifiednode in its output.
type Visitor interface {
	Visit(node interface{}) (modifiednode interface{})
}


func walkIdent(v Visitor, x *ast.Ident) *ast.Ident {
	if x != nil {
		return Walk(v, x).(*ast.Ident)
	}
	return nil
}


func walkCommentGroup(v Visitor, g *ast.CommentGroup) *ast.CommentGroup {
	if g != nil {
		return Walk(v, g).(*ast.CommentGroup)
	}
	return nil
}


func walkBlockStmt(v Visitor, b *ast.BlockStmt) *ast.BlockStmt {
	if b != nil {
		return Walk(v, b).(*ast.BlockStmt)
	}
	return nil
}


// Walk traverses an AST in depth-first order: If node != nil, it
// invokes v.Visit(node). If the modifiednode returned by
// v.Visit(node) is not nil, Walk returns a node with this node
// modified.
//
// Walk may be called with any of the named ast node types. It also
// accepts arguments of type []*ast.Field, []*ast.Ident, []ast.Expr, []ast.Stmt and []ast.Decl;
// the respective children are the slice elements.
//
func Walk(v Visitor, node interface{}) (modifiednode interface{}) {
	if node == nil {
		return
	}
	if newnode := v.Visit(node); newnode != nil {
		return newnode
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding declaration in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *ast.Comment:
		// nothing else to do

	case *ast.CommentGroup:
		for i, c := range n.List {
			n.List[i] = Walk(v, c).(*ast.Comment)
		}

	case *ast.Field:
		n.Doc = walkCommentGroup(v, n.Doc)
		n.Names = Walk(v, n.Names).([]*ast.Ident)
		n.Type = Walk(v, n.Type).(ast.Expr)
		n.Tag = Walk(v, n.Tag).(*ast.BasicLit)
		n.Comment = walkCommentGroup(v, n.Comment)

	case *ast.FieldList:
		for i, f := range n.List {
			n.List[i] = Walk(v, f).(*ast.Field)
		}

	// ast.Expressions
	case *ast.BadExpr, *ast.Ident, *ast.Ellipsis, *ast.BasicLit:
		// nothing to do

	case *ast.FuncLit:
		if n != nil {
			n.Type = Walk(v, n.Type).(*ast.FuncType)
		}
		n.Body = walkBlockStmt(v, n.Body)

	case *ast.CompositeLit:
		n.Type = Walk(v, n.Type).(ast.Expr)
		n.Elts = Walk(v, n.Elts).([]ast.Expr)

	case *ast.ParenExpr:
		n.X = Walk(v, n.X).(ast.Expr)

	case *ast.SelectorExpr:
		n.X = Walk(v, n.X).(ast.Expr)
		n.Sel = walkIdent(v, n.Sel)

	case *ast.IndexExpr:
		n.X = Walk(v, n.X).(ast.Expr)
		n.Index = Walk(v, n.Index).(ast.Expr)

	case *ast.SliceExpr:
		n.X = Walk(v, n.X).(ast.Expr)
		n.Index = Walk(v, n.Index).(ast.Expr)
		n.End = Walk(v, n.End).(ast.Expr)

	case *ast.TypeAssertExpr:
		n.X = Walk(v, n.X).(ast.Expr)
		n.Type = Walk(v, n.Type).(ast.Expr)

	case *ast.CallExpr:
		n.Fun = Walk(v, n.Fun).(ast.Expr)
		n.Args = Walk(v, n.Args).([]ast.Expr)

	case *ast.StarExpr:
		n.X = Walk(v, n.X).(ast.Expr)

	case *ast.UnaryExpr:
		n.X = Walk(v, n.X).(ast.Expr)

	case *ast.BinaryExpr:
		n.X = Walk(v, n.X).(ast.Expr)
		n.Y = Walk(v, n.Y).(ast.Expr)

	case *ast.KeyValueExpr:
		n.Key = Walk(v, n.Key).(ast.Expr)
		n.Value = Walk(v, n.Value).(ast.Expr)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			n.Len = Walk(v, n.Len).(ast.Expr)
		}
		n.Elt = Walk(v, n.Elt).(ast.Expr)

	case *ast.StructType:
		n.Fields = Walk(v, n.Fields).(*ast.FieldList)

	case *ast.FuncType:
		n.Params = Walk(v, n.Params).(*ast.FieldList)
		if n.Results != nil {
			n.Results = Walk(v, n.Results).(*ast.FieldList)
		}

	case *ast.InterfaceType:
		n.Methods = Walk(v, n.Methods).(*ast.FieldList)

	case *ast.MapType:
		n.Key = Walk(v, n.Key).(ast.Expr)
		n.Value = Walk(v, n.Value).(ast.Expr)

	case *ast.ChanType:
		n.Value = Walk(v, n.Value).(ast.Expr)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		n.Decl = Walk(v, n.Decl).(ast.Decl)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		n.Label = walkIdent(v, n.Label)
		n.Stmt = Walk(v, n.Stmt).(ast.Stmt)

	case *ast.ExprStmt:
		n.X = Walk(v, n.X).(ast.Expr)

	case *ast.IncDecStmt:
		n.X = Walk(v, n.X).(ast.Expr)

	case *ast.AssignStmt:
		n.Lhs = Walk(v, n.Lhs).([]ast.Expr)
		n.Rhs = Walk(v, n.Rhs).([]ast.Expr)

	case *ast.GoStmt:
		if n.Call != nil {
			n.Call = Walk(v, n.Call).(*ast.CallExpr)
		}

	case *ast.DeferStmt:
		if n.Call != nil {
			n.Call = Walk(v, n.Call).(*ast.CallExpr)
		}

	case *ast.ReturnStmt:
		n.Results = Walk(v, n.Results).([]ast.Expr)

	case *ast.BranchStmt:
		n.Label = walkIdent(v, n.Label)

	case *ast.BlockStmt:
		n.List = Walk(v, n.List).([]ast.Stmt)

	case *ast.IfStmt:
		if n.Init != nil {
			n.Init = Walk(v, n.Init).(ast.Stmt)
		}
		n.Cond = Walk(v, n.Cond).(ast.Expr)
		n.Body = walkBlockStmt(v, n.Body)
		if n.Else != nil {
			n.Else = Walk(v, n.Else).(ast.Stmt)
		}

	case *ast.CaseClause:
		n.Values = Walk(v, n.Values).([]ast.Expr)
		n.Body = Walk(v, n.Body).([]ast.Stmt)

	case *ast.SwitchStmt:
		n.Init = Walk(v, n.Init).(ast.Stmt)
		n.Tag = Walk(v, n.Tag).(ast.Expr)
		n.Body = walkBlockStmt(v, n.Body)

	case *ast.TypeCaseClause:
		n.Types = Walk(v, n.Types).([]ast.Expr)
		n.Body = Walk(v, n.Body).([]ast.Stmt)

	case *ast.TypeSwitchStmt:
		n.Init = Walk(v, n.Init).(ast.Stmt)
		n.Assign = Walk(v, n.Assign).(ast.Stmt)
		n.Body = walkBlockStmt(v, n.Body)

	case *ast.CommClause:
		n.Lhs = Walk(v, n.Lhs).(ast.Expr)
		n.Rhs = Walk(v, n.Rhs).(ast.Expr)
		n.Body = Walk(v, n.Body).([]ast.Stmt)

	case *ast.SelectStmt:
		n.Body = walkBlockStmt(v, n.Body)

	case *ast.ForStmt:
		n.Init = Walk(v, n.Init).(ast.Stmt)
		n.Cond = Walk(v, n.Cond).(ast.Expr)
		n.Post = Walk(v, n.Post).(ast.Stmt)
		n.Body = walkBlockStmt(v, n.Body)

	case *ast.RangeStmt:
		n.Key = Walk(v, n.Key).(ast.Expr)
		n.Value = Walk(v, n.Value).(ast.Expr)
		n.X = Walk(v, n.X).(ast.Expr)
		n.Body = walkBlockStmt(v, n.Body)

	// ast.Declarations
	case *ast.ImportSpec:
		n.Doc = walkCommentGroup(v, n.Doc)
		n.Name = walkIdent(v, n.Name)
		n.Path = Walk(v, n.Path).(*ast.BasicLit)
		n.Comment = walkCommentGroup(v, n.Comment)

	case *ast.ValueSpec:
		n.Doc = walkCommentGroup(v, n.Doc)
		n.Names = Walk(v, n.Names).([]*ast.Ident)
		n.Type = Walk(v, n.Type).(ast.Expr)
		n.Values = Walk(v, n.Values).([]ast.Expr)
		n.Comment = walkCommentGroup(v, n.Comment)

	case *ast.TypeSpec:
		n.Doc = walkCommentGroup(v, n.Doc)
		n.Name = walkIdent(v, n.Name)
		n.Type = Walk(v, n.Type).(ast.Expr)
		n.Comment = walkCommentGroup(v, n.Comment)

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		n.Doc = walkCommentGroup(v, n.Doc)
		for i, s := range n.Specs {
			n.Specs[i] = Walk(v, s).(ast.Spec)
		}

	case *ast.FuncDecl:
		n.Doc = walkCommentGroup(v, n.Doc)
		if n.Recv != nil {
			n.Recv = Walk(v, n.Recv).(*ast.FieldList)
		}
		n.Name = walkIdent(v, n.Name)
		if n.Type != nil {
			n.Type = Walk(v, n.Type).(*ast.FuncType)
		}
		n.Body = walkBlockStmt(v, n.Body)

	// Files and packages
	case *ast.File:
		n.Doc = walkCommentGroup(v, n.Doc)
		n.Name = walkIdent(v, n.Name)
		n.Decls = Walk(v, n.Decls).([]ast.Decl)
		for i, g := range n.Comments {
			n.Comments[i] = Walk(v, g).(*ast.CommentGroup)
		}

	case *ast.Package:
		for i, f := range n.Files {
			n.Files[i] = Walk(v, f).(*ast.File)
		}

	case []*ast.Ident:
		for i, x := range n {
			n[i] = Walk(v, x).(*ast.Ident)
		}

	case []ast.Expr:
		for i, x := range n {
			n[i] = Walk(v, x).(ast.Expr)
		}

	case []ast.Stmt:
		for i, x := range n {
			n[i] = Walk(v, x).(ast.Stmt)
		}

	case []ast.Decl:
		for i, x := range n {
			n[i] = Walk(v, x).(ast.Decl)
		}

	default:
		fmt.Printf("ast.Walk: unexpected type %T", n)
		panic("ast.Walk")
	}

	return node
}
