package jsonschema

import (
	"bytes"
	"encoding/json"
	"io"
	"path"
	"reflect"
)

const (
	// RefRoot is root of JSON Schema reference.
	RefRoot = "#/"
)

// Generator generates a JSON Schema.
type Generator interface {
	JSONSchema(w io.Writer, opts ...Option) error
}

// Generate generates JSON Schema from a Go type.
// Channel, complex, and function values cannot be encoded in JSON Schema.
// Attempting to generate such a type causes Generate to return
// an UnsupportedTypeError.
func Generate(w io.Writer, v interface{}, opts ...Option) error {

	if g, ok := v.(Generator); ok {
		return g.JSONSchema(w, opts...)
	}

	var g gen
	o := &obj{
		m:   map[string]interface{}{},
		ref: RefRoot,
	}

	if err := g.do(o, reflect.ValueOf(v), opts...); err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(o.m)
}

type gen struct{}

func (g *gen) do(o Object, v reflect.Value, options ...Option) error {

	switch v.Kind() {
	case reflect.Interface, reflect.Chan, reflect.Func,
		reflect.Ptr, reflect.Map, reflect.Slice:
		if v.IsNil() {
			return nil
		}
	}

	if g1, ok := v.Interface().(Generator); ok {

		var buf bytes.Buffer
		if err := g1.JSONSchema(&buf, options...); err != nil {
			return err
		}

		var m map[string]interface{}
		if err := json.NewDecoder(&buf).Decode(&m); err != nil {
			return err
		}

		for k, v := range m {
			o.Set(k, v)
		}

		return nil
	}

	switch v.Kind() {
	// unsupported types
	case reflect.Complex64, reflect.Complex128, reflect.Interface,
		reflect.Chan, reflect.Func, reflect.Invalid, reflect.UnsafePointer:
		return &json.UnsupportedTypeError{v.Type()}
	case reflect.Ptr:
		return g.do(o, v.Elem(), options...)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64:
		o.Set("type", "number")
	case reflect.Bool:
		o.Set("type", "boolean")
	case reflect.String:
		o.Set("type", "string")
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return &json.UnsupportedTypeError{v.Type()}
		}
		o.Set("type", "object")
	case reflect.Array, reflect.Slice:
		if err := g.arrayGen(o, v, options...); err != nil {
			return err
		}
	case reflect.Struct:
		if err := g.structGen(o, v, options...); err != nil {
			return err
		}
	}

	for _, opt := range options {
		var err error
		o, err = opt(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *gen) arrayGen(parent Object, v reflect.Value, options ...Option) error {
	o := &obj{
		m:   map[string]interface{}{},
		ref: path.Join(parent.Ref(), "items"),
	}

	elm := reflect.Zero(v.Type().Elem())
	if v.Len() != 0 {
		elm = v.Index(0)
	}
	if err := g.do(o, elm, options...); err != nil {
		return err
	}

	parent.Set("type", "array")
	parent.Set("items", o.m)

	return nil
}

func (g *gen) structGen(parent Object, v reflect.Value, options ...Option) error {
	required := make([]string, v.NumField())
	properties := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		f, ft := v.Field(i), v.Type().Field(i)
		name := ft.Name

		if ft.Anonymous {
			name = ft.Type.Name()
		}

		if tag, ok := ft.Tag.Lookup("json"); ok {
			name = tag
		}

		required[i] = name

		o := &obj{
			m:   map[string]interface{}{},
			ref: path.Join(parent.Ref(), "properties", name),
		}

		opts := make([]Option, len(options)+1)
		copy(opts, options)
		opts[len(opts)-1] = ByReference(o.Ref(), PropertyOrder(i))

		if err := g.do(o, f, opts...); err != nil {
			return err
		}

		properties[name] = o.m
	}

	parent.Set("type", "object")
	if title := v.Type().Name(); title != "" {
		parent.Set("title", title)
	}
	parent.Set("required", required)
	parent.Set("properties", properties)

	return nil
}
