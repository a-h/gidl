package publictypes

import _ "embed"

//go:embed snapshot.json
var Expected string

type Public struct {
	A string
	B string
}
