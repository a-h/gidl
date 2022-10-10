package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/a-h/gidl/model"
)

func main() {
	//_, err := model.Parse("github.com/a-h/gidl/model/example", "./example")
	m, err := model.Get("github.com/a-h/gidl/tests/complete")
	if err != nil {
		log.Fatalf("failed to parse: %v", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(m)
}
