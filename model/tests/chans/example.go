package chans

import _ "embed"

//go:embed snapshot.json
var Expected string

type Data struct {
	AndIgnoreThisToo chan string
}

var IgnoreThisChannel chan string
