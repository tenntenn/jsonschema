package main

import (
	"bytes"
	"encoding/json"
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

	p := Person{
		Name: "tenntenn",
		Age:  33,
	}
	var buf bytes.Buffer
	if err := jsonschema.Generate(&buf, p); err != nil {
		log.Fatal(err)
	}

	var val bytes.Buffer
	if err := json.NewEncoder(&val).Encode(p); err != nil {
		log.Fatal(err)
	}

	h, err := handler.New(&buf, handler.WithJSON(&val))
	if err != nil {
		log.Fatal(err)
	}

	h = handler.PostToWriter(os.Stdout, h)
	log.Fatal(http.ListenAndServe(":8080", h))
}
