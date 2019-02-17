package handler

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

type Options struct {
	JSON     string
	Template *template.Template
}

type Option func(o *Options) error

func WithJSON(r io.Reader) Option {
	return func(o *Options) error {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		o.JSON = string(b)
		return nil
	}
}

// New creates a new http.Handler which provides editor of schema.
func New(schema io.Reader, options ...Option) (http.Handler, error) {
	var data struct {
		Schema string
		JSON   string
	}
	b, err := ioutil.ReadAll(schema)
	if err != nil {
		return nil, err
	}
	data.Schema = string(b)

	var opts Options
	for _, o := range options {
		if err := o(&opts); err != nil {
			return nil, err
		}
	}
	data.JSON = opts.JSON

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := opts.Template
		if tmpl == nil {
			tmpl = Template
		}

		if err := tmpl.Execute(w, data); err != nil {
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
