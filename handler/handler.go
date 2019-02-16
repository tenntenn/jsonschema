package handler

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

// New creates a new http.Handler which provides editor of schema.
// action is a handler to receive JSONs which are edited by the editor.
func New(action http.Handler, schema io.Reader) (http.Handler, error) {
	return WithTemplate(action, schema, Template)
}

// WithTemplate creates a new http.Handler which provides editor of schema with given template.
func WithTemplate(action http.Handler, schema io.Reader, tmpl *template.Template) (http.Handler, error) {
	b, err := ioutil.ReadAll(schema)
	if err != nil {
		return nil, err
	}
	schemaStr := string(b)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			status := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(status), status)
		case http.MethodPost:
			action.ServeHTTP(w, r)
		case http.MethodGet:
			if err := tmpl.Execute(w, schemaStr); err != nil {
				status := http.StatusInternalServerError
				http.Error(w, err.Error(), status)
			}
		}
	}), nil
}

// ToWriter copies request body to the writer w with a new line.
func ToWriter(w io.Writer) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if _, err := io.Copy(w, r.Body); err != nil {
			status := http.StatusInternalServerError
			http.Error(rw, err.Error(), status)
		}
		fmt.Fprintln(w)
	})
}
