package privatetypes

import _ "embed"

//go:embed snapshot.json
var Expected string

type private struct {
	A string
	B string
}
