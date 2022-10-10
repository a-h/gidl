package functiontypes

import _ "embed"

//go:embed snapshot.json
var Expected string

type Data struct {
	A func(test string)
	B FuncType
}

type FuncType func()
