package model

import (
	"fmt"
	"go/doc"
	"go/types"
)

func New() *Model {
	m := &Model{
		Types: make(map[string]*Type),
	}
	m.warnf = func(format string, a ...any) {
		m.Warnings = append(m.Warnings, fmt.Sprintf(format, a...))
	}
	return m
}

type Model struct {
	Types    map[string]*Type `json:"types"`
	warnf    func(format string, a ...any)
	Warnings []string `json:"warnings,omitempty"`
}

func (m *Model) addType(n *types.Named) {
	t := Type{
		ID:     n.Origin().String(),
		Name:   n.Obj().Name(),
		Fields: m.getFields(n),
	}
	m.Types[t.ID] = &t
}

func (m *Model) getFields(n *types.Named) (fields []*Field) {
	s, ok := n.Underlying().(*types.Struct)
	if !ok {
		return
	}
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)
		name := f.Name()
		id := fmt.Sprintf("%s.%s", n.Origin().String(), name)
		is, desc, ok := getFieldType(f.Type())
		if !ok {
			m.warnf("%s - field type %q cannot be mapped", id, desc)
			continue
		}
		fields = append(fields, &Field{
			ID:   id,
			Name: name,
			Is:   is,
			Tags: s.Tag(i),
		})
	}
	return
}

func getFieldType(t types.Type) (is Is, desc string, ok bool) {
	switch t := t.(type) {
	case *types.Basic:
		desc = t.String()
		is.Scalar = &Scalar{
			Of: TypeName(t.String()),
		}
		ok = true
		return
	case *types.Slice:
		of, desc, ok := getFieldType(t.Elem())
		if !ok {
			return is, desc, false
		}
		is.Array = &Array{Of: of}
		is.Nullable = true
		return is, t.String(), true
	case *types.Array:
		of, desc, ok := getFieldType(t.Elem())
		if !ok {
			return is, desc, false
		}
		is.Array = &Array{Of: of}
		is.Nullable = false
		return is, t.String(), true
	case *types.Interface:
		desc = "interface: " + t.String()
		return
	case *types.TypeParam:
		desc = "typeparam: " + t.String()
		return
	case *types.Pointer:
		is, desc, ok = getFieldType(t.Elem())
		if !ok {
			return is, desc, false
		}
		is.Nullable = true
		return is, t.String(), true
	case *types.Map:
		var k, v Is
		if k, desc, ok = getFieldType(t.Key()); !ok {
			return is, desc, false
		}
		if v, desc, ok = getFieldType(t.Elem()); !ok {
			return is, desc, false
		}
		is.Map = &Map{
			FromKey: k,
			ToValue: v,
		}
		is.Nullable = true
		return is, t.String(), true
	case *types.Named:
		desc = t.String()
		// Disallow generics.
		if t.TypeParams().Len() > 0 {
			return
		}
		// Disallow function types.
		if _, isFunction := t.Underlying().(*types.Signature); isFunction {
			return
		}
		// Disallow channels.
		if _, isChan := t.Underlying().(*types.Chan); isChan {
			return
		}

		// Allowed type.
		is.Scalar = &Scalar{
			Of: TypeName(t.Origin().String()),
		}
		return is, desc, true
	}
	desc = t.String()
	return
}

func (m *Model) setTypeComment(typeID, comment string) {
	t, ok := m.Types[typeID]
	if !ok {
		return
	}
	t.Comments = comment
	t.Description = doc.Synopsis(t.Comments)
}

func (m *Model) setFieldComment(typeID, fieldName, comment string) {
	t, ok := m.Types[typeID]
	if !ok {
		return
	}
	for i, f := range t.Fields {
		if f.Name == fieldName {
			t.Fields[i].Comments = comment
			t.Fields[i].Description = doc.Synopsis(comment)
			return
		}
	}
	m.warnf("%s - failed to set comment on field name %q", typeID, fieldName)
}

func (m *Model) setEnumStringValue(typeID, value, comment string) {
	t, ok := m.Types[typeID]
	if !ok {
		return
	}
	t.EnumStringValues = append(t.EnumStringValues, EnumValue[string]{Value: value, Description: comment})
}

func (m *Model) setEnumIntValue(typeID string, value int64, comment string) {
	t, ok := m.Types[typeID]
	if !ok {
		return
	}
	t.EnumIntValues = append(t.EnumIntValues, EnumValue[int64]{Value: value, Description: comment})
}
