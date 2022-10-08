package model

import "go/doc"

func New() *Model {
	return &Model{
		Types: make(map[string]*Type),
	}
}

type Model struct {
	Types map[string]*Type
}

func (m *Model) SetTypeComment(typeID, comment string) {
	t, ok := m.Types[typeID]
	if !ok {
		return
	}
	t.RawComments = comment
	t.Description = doc.Synopsis(t.RawComments)
}

func (m *Model) SetFieldComment(typeID, fieldName, comment string) {
	t, ok := m.Types[typeID]
	if !ok {
		return
	}
	for i, f := range t.Fields {
		if f.Name == fieldName {
			t.Fields[i].RawComments = comment
			t.Fields[i].Description = doc.Synopsis(comment)
		}
	}
}

type Type struct {
	// ID of the type.
	// For Go based sources, it is the fully qualified type name.
	// github.com/a-h/gidl/model.Type
	ID string `json:"id"`
	// Name of the wire representation of the type, e.g. field.
	Name        string   `json:"name"`
	Description string   `json:"desc,omitempty"`
	Fields      []*Field `json:"fields,omitempty"`
	Traits      []Trait  `json:"traits,omitempty"`
	RawComments string   `json:"rawComments,omitempty"`
}

type Trait string

// TypeName is either a built-in type, or a fully qualified package name.
type TypeName string

const (
	TypeNameInt     TypeName = "int"
	TypeNameString  TypeName = "string"
	TypeNameFloat64 TypeName = "float64"
)

var KnownFieldTypes = []TypeName{TypeNameInt, TypeNameString, TypeNameFloat64}

func IsKnownFieldType(t string) bool {
	for _, kft := range KnownFieldTypes {
		if t == string(kft) {
			return true
		}
	}
	return false
}

type FieldType struct {
	Name TypeName `json:"name"`
	// Optional is set to true if the field type is a pointer.
	// Since slices and maps are also pointers, they're optional by default too.
	Optional bool `json:"optional,omitempty"`
	// Multiple is set to true if the field type is an array.
	Multiple bool `json:"multiple,omitempty"`
}

type Field struct {
	ID string `json:"id"`
	// Name of the wire representation of the type, e.g. field.
	Name string `json:"name"`
	// Description of the field usage.
	Description string `json:"desc,omitempty"`
	// Type of the field.
	Type FieldType `json:"type,omitempty"`
	// Examples of the data stored in the field.
	// abc
	// 123
	// 0x4A
	Examples    []string `json:"examples,omitempty"`
	Traits      []Trait  `json:"traits,omitempty"`
	RawComments string   `json:"rawComments,omitempty"`
	Tags        string   `json:"tags,omitempty"`
}
