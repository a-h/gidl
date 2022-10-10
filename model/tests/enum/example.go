package enum

import _ "embed"

//go:embed snapshot.json
var Expected string

type Public struct {
	A StringEnum
	B IntEnum
}

type StringEnum string

const (
	StringEnumA StringEnum = "A"
	StringEnumB StringEnum = "B"
	StringEnumC StringEnum = "C"
)

type IntEnum int

const (
	IntEnum0 IntEnum = iota
	IntEnum1
	IntEnum2
)
