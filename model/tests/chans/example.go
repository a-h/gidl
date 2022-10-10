package chans

import _ "embed"

//go:embed snapshot.json
var Expected string

type Data struct {
	AndIgnoreThisToo       chan string
	AndThisArrayOfChannels []chan string
	AndThisAlias           ChanType
}

var IgnoreThisChannel chan string

type ChanType chan string
