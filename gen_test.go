package jsonschema_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	jd "github.com/josephburnett/jd/lib"
	"github.com/tenntenn/jsonschema"
	. "github.com/tenntenn/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func toJSON(t *testing.T, v interface{}) string {
	t.Helper()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatal("unexpected error", err)
	}
	return buf.String()
}

func jsonCompact(t *testing.T, src string) string {
	t.Helper()
	var dst bytes.Buffer
	if err := json.Compact(&dst, []byte(src)); err != nil {
		t.Fatal("unexpected error:", err)
	}
	return dst.String()
}

func jsonDiff(t *testing.T, a, b string) string {
	t.Helper()
	jsonA, err := jd.ReadJsonString(a)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	jsonB, err := jd.ReadJsonString(b)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	return jsonA.Diff(jsonB).Render()
}

type generator struct {
	json   string
	schema string
}

func (g *generator) JSONSchema(w io.Writer, opts ...jsonschema.Option) error {
	fmt.Fprint(w, g.schema)
	return nil
}

func (g *generator) MarshalJSON() ([]byte, error) {
	return []byte(g.json), nil
}

func TestGenerate(t *testing.T) {

	type T struct {
		N int
		S string
	}

	type NT struct {
		T T
	}

	cases := []struct {
		name   string
		v      interface{}
		expect string
		isErr  bool
	}{
		{
			name:   "int",
			v:      100,
			expect: `{"type":"number"}`,
		},
		{
			name:   "string",
			v:      "example",
			expect: `{"type":"string"}`,
		},
		{
			name: "int array",
			v:    []int{10, 20, 30},
			expect: `{
				"type":"array",
				"items": {"type": "number"}
			}`,
		},
		{
			name: "empty array",
			v:    []int{},
			expect: `{
				"type":"array",
				"items": {"type": "number"}
			}`,
		},
		{
			name:   "nil array",
			v:      []int(nil),
			expect: `{}`,
		},
		{
			name: "struct",
			v:    T{N: 100, S: ""},
			expect: `{
				"title": "T",
				"type":"object",
				"required": ["N", "S"],
				"properties":{
					"N":{
						"type":"number",
						"propertyOrder": 0
					},
					"S":{
						"type":"string",
						"propertyOrder": 1
					}
				}
			}`,
		},
		{
			name: "nested struct",
			v:    NT{T: T{N: 100, S: ""}},
			expect: `{
				"type":"object",
				"title": "NT",
				"required": ["T"],
				"properties": {
					"T": {
						"title": "T",
						"type":"object",
						"propertyOrder": 0,
						"required": ["N", "S"],
						"properties":{
							"N":{
								"type":"number",
								"propertyOrder": 0
							},
							"S":{
								"type":"string",
								"propertyOrder": 1
							}
						}
					}
				}
			}`,
		},
		{
			name: "generator",
			v: &generator{
				json:   `100`,
				schema: `{"type":"number"}`,
			},
			expect: `{"type":"number"}`,
		},
		{
			name: "generator in struct",
			v: struct {
				V *generator `json:"v"`
			}{
				V: &generator{
					json:   `100`,
					schema: `{"type":"number"}`,
				},
			},
			expect: `{
				"type":"object",
				"required": ["v"],
				"properties": {
					"v": {
						"type": "number"
					}
				}
			}`,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if err, isErr := r.(error); isErr {
					switch {
					case tt.isErr && err == nil:
						t.Errorf("expected error does not occur")
					case !tt.isErr && err != nil:
						t.Errorf("unexpected error: %v", err)
					}
				} else if r != nil {
					panic(r)
				}
			}()
			var buf bytes.Buffer
			errCheck(Generate(&buf, tt.v))
			got := buf.String()

			if diff := jsonDiff(t, got, tt.expect); diff != "" {
				t.Fatalf("generated JSON Schema does not match to expected one: %v", diff)
			}

			l := gojsonschema.NewStringLoader(got)
			s, err := gojsonschema.NewSchema(l)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			r, err := s.Validate(gojsonschema.NewStringLoader(toJSON(t, tt.v)))
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			if !r.Valid() {
				t.Errorf("invalid JSON Schema: %s", got)
			}
		})
	}
}
