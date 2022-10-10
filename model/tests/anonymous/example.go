package anonymous

import _ "embed"

//go:embed snapshot.json
var Expected string

// Anonymous structs. Nope.

type Data struct {
	A struct {
		B string
	}
}

