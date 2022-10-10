package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/a-h/gidl/model"
)

func main() {
	m, err := model.Get("github.com/a-h/gidl/example")
	if err != nil {
		log.Fatalf("failed to parse: %v", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(m)
	if err != nil {
		fmt.Printf("error encoding: %v\n", err)
		os.Exit(1)
	}
	for _, warning := range m.Warnings {
		fmt.Println(warning)
	}
}
