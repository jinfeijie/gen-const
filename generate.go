package main

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"golang.org/x/tools/go/packages"
	"log"
	"reflect"
	"strings"
)

type Generator struct {
	pkg        *Package
	trimPrefix string
}

func (g *Generator) format(headBuff, bodyBuff *Buffer) []byte {
	if strings.Contains(string(bodyBuff.Bytes()), "fmt.") {
		headBuff.WriteS("import \"fmt\"\n")
	}
	headBuff.WriteS("\n")
	headBuff.Write(bodyBuff.Bytes())
	src, err := format.Source(headBuff.Bytes())
	if err != nil {
		log.Printf("警告：内部错误：Go 生成无效： %s", err)
		log.Printf("警告：编译包分析错误")
		return headBuff.Bytes()
	}
	return src
}
func (g *Generator) generate(typeName string, buff *Buffer) {
	values := make([]Value, 0, 100)
	structTypeName := typeName + "Type"
	for _, file := range g.pkg.files {
		file.typeName = typeName
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}
	if len(values) == 0 {
		log.Fatalf("%s 没有为类型定义值", typeName)
	}
	t := reflect.TypeOf(values[0])
	v := reflect.ValueOf(values[0])
	buff.WriteF("type %s struct {\n", structTypeName)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() { // 导出字段
			f := t.Field(i)
			if val, ok := v.Field(i).Interface().(constant.Value); ok {
				buff.WriteF("\t%s\t%s\n", f.Name, strings.ToLower(val.Kind().String()))
			} else {
				buff.WriteF("\t%s\t%s\n", f.Name, f.Type.String())
			}
		}
	}
	buff.WriteS("}\n")
	buff.WriteS("\n")
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() { // 导出字段
			f := t.Field(i)
			if val, ok := v.Field(i).Interface().(constant.Value); ok {
				buff.WriteF("func (receiver *%s) Get%s() %s {\n", structTypeName, f.Name, strings.ToLower(val.Kind().String()))
			} else {
				buff.WriteF("func (receiver *%s) Get%s() %s {\n", structTypeName, f.Name, f.Type.String())
			}
			buff.WriteF("\treturn receiver.%s\n", f.Name)
			buff.WriteS("}\n")
			buff.WriteS("\n")
		}
	}
	var params []string
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() { // 导出字段
			f := t.Field(i)
			if val, ok := v.Field(i).Interface().(constant.Value); ok {
				params = append(params, strings.ToLower(f.Name)+" "+strings.ToLower(val.Kind().String()))
			} else {
				params = append(params, strings.ToLower(f.Name)+" "+f.Type.String())
			}
		}
	}
	buff.WriteF("func %sFunc(%s) *%s {\n", typeName, strings.Join(params, ", "), structTypeName)
	buff.WriteF("\treturn &%s{\n", structTypeName)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() { // 导出字段
			f := t.Field(i)
			buff.WriteF("\t\t%s: %s,\n", f.Name, strings.ToLower(f.Name))
		}
	}
	buff.WriteS("\t}\n")
	buff.WriteS("}\n")
	buff.WriteS("\n")
	buff.WriteF("func (receiver *%s) String() string {\n", structTypeName)
	buff.WriteF("\treturn \"%s (", structTypeName)
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() { // 导出字段
			f := t.Field(i)
			if val, ok := v.Field(i).Interface().(constant.Value); ok {
				if val.Kind() == constant.String {
					fields = append(fields, fmt.Sprintf("%s: \"+receiver.%s+\"", f.Name, f.Name))
				} else {
					fields = append(fields, fmt.Sprintf(`%s: "+fmt.Sprintf("%s", receiver.%s)+"`, f.Name, "%+v", f.Name))
				}
			} else {
				if f.Type.Kind() != reflect.String {
					fields = append(fields, fmt.Sprintf(`%s: "+fmt.Sprintf("%s", receiver.%s)+"`, f.Name, "%+v", f.Name))
				} else {
					fields = append(fields, fmt.Sprintf("%s: \"+receiver.%s+\"", f.Name, f.Name))
				}
			}
		}
	}
	buff.WriteS(strings.Join(fields, ", "))
	buff.WriteS(")\"\n")
	buff.WriteS("}\n")
	buff.WriteS("\n")
	buff.WriteS("var (\n")
	for _, v := range values {
		originalName := v.originalName
		if len(originalName) > 0 {
			v := []byte(originalName)
			x := string(v[0])
			t := []byte(strings.ToUpper(x))
			v[0] = t[0]
			originalName = string(v)
		}
		buff.WriteF("\t%s%s = %s(%+v, \"%s\")\n", originalName, typeName, typeName+"Func", v.Val, v.Msg)
	}
	buff.WriteS(")\n")
}

func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedTypesSizes,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("错误：%d 个软件包与 %v 匹配", len(pkgs), strings.Join(patterns, " "))
	}
	g.addPackage(pkgs[0])
}
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}
	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			typeName:    file.Name.String(),
			file:        file,
			pkg:         g.pkg,
			trimPrefix:  g.trimPrefix,
			lineComment: true,
		}
	}
}
