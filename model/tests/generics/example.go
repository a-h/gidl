package generics

import _ "embed"

//go:embed snapshot.json
var Expected string

type Data struct {
	AllowThis string
	String    DataOfT[string]
	Int       DataOfT[int]
}

type DataOfT[T string | int] struct {
	Field T
}
