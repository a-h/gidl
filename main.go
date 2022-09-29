package main

import (
	"log"

	"github.com/a-h/gidl/model"
)

func main() {
	_, err := model.Parse("github.com/a-h/gidl/model/example", "./example")
	if err != nil {
		log.Fatalf("failed to parse: %v", err)
	}
}
