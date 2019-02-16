package jsonschema_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/tenntenn/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

func toJSON(t *testing.T, v interface{}) string {
	t.Helper()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatal("unexpected error", err)
	}
	return buf.String()
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
		name  string
		v     interface{}
		isErr bool
	}{
		{"int", 100, false},
		{"string", "example", false},
		{"struct", T{N: 100, S: ""}, false},
		{"nested struct", NT{T: T{N: 100, S: ""}}, false},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer
			err := Generate(&got, tt.v)
			switch {
			case tt.isErr && err == nil:
				t.Errorf("expected error does not occur")
			case !tt.isErr && err != nil:
				t.Errorf("unexpected error %v", err)
			}

			l := gojsonschema.NewStringLoader(got.String())
			s, err := gojsonschema.NewSchema(l)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			r, err := s.Validate(gojsonschema.NewStringLoader(toJSON(t, tt.v)))
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}

			if !r.Valid() {
				t.Errorf("invalid JSON Schema: %s", got.String())
			}
		})
	}
}
