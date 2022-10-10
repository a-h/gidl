package model

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

func Get(packageName string) (m *Model, err error) {
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

	m = New()
	// Read through the definitions.
	for _, pkg := range pkgs {
		identifiers := getSortedKeys(pkg.TypesInfo.Defs)
		for _, identifier := range identifiers {
			// Only process exported types.
			if !identifier.IsExported() {
				continue
			}
			definition := pkg.TypesInfo.Defs[identifier]
			// Anything to be considered must have a definition.
			if definition == nil {
				continue
			}
			// Named types.
			if identifier.Obj != nil && identifier.Obj.Kind == ast.Typ {
				n, isNamedType := definition.Type().(*types.Named)
				if !isNamedType {
					continue
				}
				m.AddType(n)
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
				case *ast.ValueSpec:
					// Identify constants that make up string and integer enums.
					for _, name := range x.Names {
						c, isConstant := pkg.TypesInfo.ObjectOf(name).(*types.Const)
						if !isConstant {
							continue
						}
						typeID := c.Type().String()
						if _, ok := m.Types[typeID]; !ok {
							// Cannot find a type that this constant belongs to.
							// So it's not an enum value.
							continue
						}
						switch c.Val().Kind() {
						case constant.String:
							m.SetEnumStringValue(typeID, constant.StringVal(c.Val()), x.Doc.Text())
						case constant.Int:
							v, _ := constant.Int64Val(c.Val())
							m.SetEnumIntValue(typeID, v, x.Doc.Text())
						default:
							m.warnf("Constant %q does not have a type of integer or string, and can't be included in an enum", name)
						}
					}
				case *ast.FuncDecl:
					// Skip functions.
					return false
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

