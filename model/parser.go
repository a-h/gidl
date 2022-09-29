package model

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	gopath "path"
	"path/filepath"
	"strings"

	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
)

func Parse(base, path string) (m *Model, err error) {
	fset := token.NewFileSet()
	dict := make(map[string][]*ast.Package)
	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			for _, v := range d {
				k := gopath.Join(base, path)
				dict[k] = append(dict[k], v)
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	m = New()
	imports := make(map[string]string)
	fields := make(map[string][]*Field)

	for pkg, p := range dict {
		for _, f := range p {
			var typ string
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					typ = x.Name.String()
					if !ast.IsExported(typ) {
						break
					}
					rawComments := x.Doc.Text()
					desc := doc.Synopsis(rawComments)
					typeID := fmt.Sprintf("%s.%s", pkg, typ)
					m.AddType(typeID, typ, desc, rawComments)
				case *ast.Field:
					if typ == "" {
						break
					}
					name := getFieldName(x)
					if !ast.IsExported(name) {
						break
					}
					typeID := fmt.Sprintf("%s.%s", pkg, typ)
					//TODO: Examples, traits.
					field := &Field{
						Name:        name,
						Description: doc.Synopsis(x.Doc.Text()),
						Type:        getFieldType(imports, pkg, x.Type),
						RawComments: x.Doc.Text(),
					}
					if x.Tag != nil {
						field.Tags = x.Tag.Value
					}
					fields[typeID] = append(fields[typeID], field)
				case *ast.ImportSpec:
					importPath := strings.TrimFunc(x.Path.Value, func(r rune) bool {
						return r == '"'
					})
					importName := importPath[strings.LastIndex(importPath, "/")+1:]
					if x.Name != nil {
						importName = x.Name.Name
					}
					imports[importName] = importPath
					//TODO: Work out what the `type PhoneNumber string` is.
				//case *ast.GenDecl:
					//// type PhoneType string
					//fmt.Println("GenDecl")
					//for _, spec := range x.Specs {
						//fmt.Println(spec)
						//fmt.Println(reflect.TypeOf(spec))
					//}
				//default:
					//fmt.Println(x, reflect.TypeOf(x))
				}
				return true
			})
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(fields)

	return
}

func getFullyQualifiedTypeName(pkg, typ string) TypeName {
	return TypeName(fmt.Sprintf("%s.%s", pkg, typ))
}

func getFieldType(imports map[string]string, pkg string, n ast.Expr) (ft FieldType) {
	switch x := n.(type) {
	case *ast.BadExpr:
		return FieldType{Name: TypeName("foundBadExpr")}
	case *ast.Ident:
		// A type within the same package, e.g. PhoneNumber or Address.
		if !IsKnownFieldType(x.Name) {
			return FieldType{
				Name: getFullyQualifiedTypeName(pkg, x.Name),
			}
		}
		// string, int etc.
		return FieldType{Name: TypeName(x.Name)}
	case *ast.Ellipsis:
		return FieldType{Name: TypeName("found Ellipsis")}
	case *ast.BasicLit:
		return FieldType{Name: TypeName("found BasicLit")}
	case *ast.FuncLit:
		return FieldType{Name: TypeName("found FuncLit")}
	case *ast.CompositeLit:
		return FieldType{Name: TypeName("found CompositeLit")}
	case *ast.ParenExpr:
		return FieldType{Name: TypeName("found ParenExpr")}
	case *ast.SelectorExpr:
		// A type within a different package, e.g. uuid.UUID.
		return FieldType{
			Name: TypeName(getPackageName(imports, x.X) + "." + x.Sel.Name),
		}
	case *ast.IndexExpr:
		return FieldType{Name: TypeName("found IndexExpr")}
	case *ast.IndexListExpr:
		return FieldType{Name: TypeName("found IndexListExpr")}
	case *ast.SliceExpr:
		return FieldType{Name: TypeName("found SliceExpr")}
	case *ast.TypeAssertExpr:
		return FieldType{Name: TypeName("found TypeAssertExpr")}
	case *ast.CallExpr:
		return FieldType{Name: TypeName("found CallExpr")}
	case *ast.StarExpr:
		// A type that has a pointer, e.g. *uuid.Google
		ft = getFieldType(imports, pkg, x.X)
		ft.Optional = true
		return ft
	case *ast.UnaryExpr:
		return FieldType{Name: TypeName("found UnaryExpr")}
	case *ast.BinaryExpr:
		return FieldType{Name: TypeName("found BinaryExpr")}
	case *ast.KeyValueExpr:
		return FieldType{Name: TypeName("found KeyValueExpr")}
	case *ast.ArrayType:
		// A slice or array, e.g. []uuid.Google
		ft = getFieldType(imports, pkg, x.Elt)
		ft.Multiple = true
		ft.Optional = true
		return ft
	case *ast.StructType:
		return FieldType{Name: TypeName("found StructType")}
	case *ast.FuncType:
		return FieldType{Name: TypeName("found FuncType")}
	case *ast.InterfaceType:
		return FieldType{Name: TypeName("found InterfaceType")}
	case *ast.MapType:
		return FieldType{Name: TypeName("found MapType")}
	case *ast.ChanType:
		return FieldType{Name: TypeName("found ChanType")}
	}
	return FieldType{
		Name: TypeName("unknown"),
	}
}

func getPackageName(imports map[string]string, n ast.Expr) string {
	switch x := n.(type) {
	case *ast.Ident:
		//TODO: Check import exists.
		return imports[x.Name]
	}
	panic("could not find package name")
}

func getFieldName(field *ast.Field) string {
	var names []string
	for _, name := range field.Names {
		names = append(names, name.Name)
	}
	return strings.Join(names, ".")
}
