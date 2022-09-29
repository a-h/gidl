package model

func New() *Model {
	return &Model{
		Types: make(map[string]*Type),
	}
}

type Model struct {
	Types map[string]*Type
}

func (m *Model) AddType(typeID, name string, desc string, rawComments string) {
	m.Types[typeID] = &Type{
		ID:          typeID,
		Name:        name,
		Description: desc,
		RawComments: rawComments,
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

func (t *Type) AddField(name, desc string) {
	t.Fields = append(t.Fields, &Field{
		Name:        name,
		Description: desc,
	})
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
