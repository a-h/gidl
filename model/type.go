package model

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
	Comments    string   `json:"comments,omitempty"`
	Is          Is       `json:"is"`
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

type Field struct {
	ID string `json:"id"`
	// Name of the wire representation of the type, e.g. field.
	Name string `json:"name"`
	// Description of the field usage.
	Description string `json:"desc,omitempty"`
	Is          Is     `json:"is"`
	// Examples of the data stored in the field.
	// abc
	// 123
	// 0x4A
	Examples []string `json:"examples,omitempty"`
	Traits   []Trait  `json:"traits,omitempty"`
	Comments string   `json:"comments,omitempty"`
	Tags     string   `json:"tags,omitempty"`
}

type Is struct {
	Scalar *Scalar `json:"scalar,omitempty"`
	Enum   *Enum   `json:"enum,omitempty"`
	Array  *Array  `json:"array,omitempty"`
	Map    *Map    `json:"map,omitempty"`
	// Nullable is set to true if the field type is a pointer.
	// Since slices and maps are also pointers, they're optional by default too.
	Nullable bool `json:"nullable"`
}
type Scalar struct {
	Of TypeName `json:"of,omitempty"`
}
type Enum struct {
	OfStrings []EnumValue[string] `json:"ofStrings,omitempty"`
	OfInts    []EnumValue[int64]  `json:"ofInts,omitempty"`
}
type EnumValue[T string | int64] struct {
	Value       T      `json:"value"`
	Description string `json:"desc"`
}
type Array struct {
	Of Is `json:"of,omitempty"`
}
type Map struct {
	FromKey Is `json:"fromKey"`
	ToValue Is `json:"toValue"`
}
