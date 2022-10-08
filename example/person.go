package example

import (
	"github.com/google/uuid"
	gid "github.com/google/uuid"
)

// Person that exists.
type Person struct {
	// ID is a UUID value.
	// Parsed as ast.SelectorExpr
	ID uuid.UUID
	// IDPtr is a pointer to a UUID.
	// Parsed as ast.StarExpr
	IDPtr *gid.UUID
	// Name of the person.
	Name string `json:"name"`
	// Age of the person.
	Age int `json:"age"`
	// PhoneNumbers is an array of numbers.
	// Parsed as asp.SliceExpr
	PhoneNumbers []PhoneNumber
	Address      Address
	// Random array of numbers.
	Random [3]int
}

type PhoneNumber struct {
	Type   PhoneType
	Number string
}

type Address struct {
	Line1 string
}

// Go doesn't have a specific enum type.
// But this idiom is the usual way of defining the allowed values.
//TODO: Represent these as possible values in the object model.
type PhoneType string

const (
	PhoneTypeMobile = "mobile"
	PhoneTypeLand   = "land"
)
