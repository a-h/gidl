package model

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

func GetPackageInfo(packageName string) (m *Model, err error) {
	config := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
	}
	pkgs, err := packages.Load(config, packageName)
	if err != nil {
		err = fmt.Errorf("error loading package %s: %w", packageName, err)
		return
	}

	//TODO: Ensure that imports are loaded, and load the correct version of modules based on that.
	// Then, create documentation for any 3rd party modules too.

	m = &Model{
		Types: map[string]*Type{},
	}
	// Read through the definitions.
	for _, pkg := range pkgs {
		identifiers := getSortedKeys(pkg.TypesInfo.Defs)
		for _, identifier := range identifiers {
			definition := pkg.TypesInfo.Defs[identifier]
			if identifier != nil {
				fmt.Println("id:", identifier.Name)
				if identifier.Obj == nil || identifier.Obj.Kind != ast.Typ {
					continue
				}
				fmt.Println("obj kind:", identifier.Obj.Kind.String())
			}
			if definition != nil {
				n, isNamedType := definition.Type().(*types.Named)
				if !isNamedType {
					continue
				}
				t := getType(n)
				m.Types[t.ID] = &t
			}
		}
	}
	// Add the comments to the definitions.
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			var lastComment string
			var typ string
			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					typ = x.Name.String()
					if !ast.IsExported(typ) {
						break
					}
					typeID := fmt.Sprintf("%s.%s", packageName, typ)
					m.SetTypeComment(typeID, lastComment)
				case *ast.GenDecl:
					lastComment = x.Doc.Text()
				case *ast.Field:
					if typ == "" {
						break
					}
					typeID := fmt.Sprintf("%s.%s", packageName, typ)
					m.SetFieldComment(typeID, getFieldName(x), x.Doc.Text())
				}
				return true
			})
		}
	}
	return
}

func getFieldName(field *ast.Field) string {
	var names []string
	for _, name := range field.Names {
		names = append(names, name.Name)
	}
	return strings.Join(names, ".")
}

func getSortedKeys(defs map[*ast.Ident]types.Object) []*ast.Ident {
	op := make([]*ast.Ident, len(defs))
	var i int
	for k := range defs {
		k := k
		op[i] = k
		i++
	}
	sort.Slice(op, func(i, j int) bool {
		return op[i].Name < op[j].Name
	})
	return op
}

func getType(n *types.Named) (t Type) {
	return Type{
		ID:     n.Origin().String(),
		Name:   n.Obj().Name(),
		Fields: getFields(n),
	}
}

func getFields(n *types.Named) (fields []*Field) {
	s, ok := n.Underlying().(*types.Struct)
	if !ok {
		return
	}
	fields = make([]*Field, s.NumFields())
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		name := f.Name()
		fields[i] = &Field{
			ID:   fmt.Sprintf("%s.%s", n.Origin().String(), name),
			Name: name,
			Type: getFieldType(f.Type()),
			Tags: s.Tag(i),
		}
	}
	return
}

func getFieldType(t types.Type) FieldType {
	switch t := t.(type) {
	case *types.Struct:
		panic("struct")
	case *types.Basic:
		return FieldType{
			Name: TypeName(t.String()),
		}
	case *types.Chan:
		panic("chan")
	case *types.Slice:
		return FieldType{
			Name:     TypeName(t.Elem().String()),
			Multiple: true,
			Optional: true,
		}
	case *types.Tuple:
		panic("tuple")
	case *types.Array:
		return FieldType{
			Name:     TypeName(t.Elem().String()),
			Optional: false,
			Multiple: false,
		}
	case *types.Interface:
		panic("interface")
	case *types.TypeParam:
		panic("type param")
	case *types.Pointer:
		return FieldType{
			Name:     TypeName(t.Elem().String()),
			Optional: true,
			Multiple: false,
		}
	case *types.Union:
		panic("union")
	case *types.Map:
		panic("map")
	case *types.Signature:
		panic("signature")
	case *types.Named:
		return FieldType{
			Name: TypeName(t.Origin().String()),
		}
	default:
		panic(fmt.Sprintf("unknown: %s", reflect.TypeOf(t)))
	}
}
