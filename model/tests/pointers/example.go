package pointers

import _ "embed"

//go:embed snapshot.json
var Expected string

type Public struct {
	A **string
}
