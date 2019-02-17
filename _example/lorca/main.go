package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/tenntenn/jsonschema"
	"github.com/tenntenn/jsonschema/handler"
	"github.com/zserge/lorca"
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
	var schema bytes.Buffer
	if err := jsonschema.Generate(&schema, p); err != nil {
		log.Fatal(err)
	}

	var val bytes.Buffer
	if err := json.NewEncoder(&val).Encode(p); err != nil {
		log.Fatal(err)
	}

	var data = struct {
		Schema string
		JSON   string
	}{
		Schema: schema.String(),
		JSON:   val.String(),
	}

	var html bytes.Buffer
	if err := handler.Template.Execute(&html, data); err != nil {
		log.Fatal(err)
	}

	ui, err := lorca.New("", "", 480, 240)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Bind("submit", func(s string) {
		fmt.Println("submit", s)
	})

	url := fmt.Sprintf("data:text/html,%s", url.PathEscape(html.String()))
	ui.Load(url)

	<-ui.Done()
}
