package functions

import _ "embed"

//go:embed snapshot.json
var Expected string

func ThisShouldBeIgnored() {
}

func asShouldThis() {
}

type Data struct {
	A string
}

func (d Data) IgnoreMe() {
}

func (d Data) andMeToo() {
}

func (d *Data) Same() {
}

type DontIgnoreMe struct {
	Please string
}
