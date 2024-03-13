package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strings"
)

type File struct {
	pkg         *Package
	file        *ast.File
	typeName    string
	values      []Value
	trimPrefix  string
	lineComment bool
}

func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		return true
	}
	typ := ""
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec)
		if vspec.Type == nil && len(vspec.Values) > 0 {
			typ = ""
			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			typ = id.Name
		}
		if vspec.Type != nil {
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != f.typeName {
			continue
		}
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("常量%s没有值", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			value := obj.(*types.Const).Val()

			v := Value{
				originalName: name.Name,
				Val:          value,
				signed:       info&types.IsUnsigned == 0,
				str:          value.String(),
			}
			if c := vspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
				v.Msg = strings.TrimSpace(c.Text())
			} else {
				v.Msg = strings.TrimPrefix(v.originalName, f.trimPrefix)
			}
			f.values = append(f.values, v)
		}
	}
	return false
}
