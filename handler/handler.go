package handler

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

// New creates a new http.Handler which provides editor of schema.
func New(schema io.Reader) (http.Handler, error) {
	return WithTemplate(schema, defaultTemplate)
}

// WithTemplate creates a new http.Handler which provides editor of schema with given template.
func WithTemplate(schema io.Reader, tmpl *template.Template) (http.Handler, error) {
	b, err := ioutil.ReadAll(schema)
	if err != nil {
		return nil, err
	}
	schemaStr := string(b)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, schemaStr); err != nil {
			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
		}
	}), nil
}

// PostToWriter copies POST request body to the writer w with a new line.
func PostToWriter(w io.Writer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			defer r.Body.Close()
			if _, err := io.Copy(w, r.Body); err != nil {
				status := http.StatusInternalServerError
				http.Error(rw, err.Error(), status)
			}
			fmt.Fprintln(w)
			return
		}
		h.ServeHTTP(rw, r)
	})
}
