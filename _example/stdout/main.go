package main

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/tenntenn/jsonschema"
	"github.com/tenntenn/jsonschema/handler"
)

func main() {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var v Person
	var buf bytes.Buffer
	if err := jsonschema.Generate(&buf, v); err != nil {
		log.Fatal(err)
	}

	h, err := handler.New(handler.ToWriter(os.Stdout), &buf)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", h))
}
