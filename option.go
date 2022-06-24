package jsonschema

import "github.com/IGLOU-EU/go-wildcard"

// Object is interface of JSON object.
type Object interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Ref() string
}

type obj struct {
	m   map[string]interface{}
	ref string
}

func (o *obj) Set(key string, value interface{}) {
	o.m[key] = value
}

func (o *obj) Get(key string) (value interface{}, ok bool) {
	value, ok = o.m[key]
	return
}

func (o *obj) Ref() string {
	return o.ref
}

// Option is options for JSON Schema.
type Option func(o Object) (Object, error)

// ByReference explicits refrence of adding option.
// It only supports refs which begins "#/".
func ByReference(pattern string, opt Option) Option {
	return func(o Object) (Object, error) {
		if wildcard.MatchSimple(pattern, o.Ref()) {
			return opt(o)
		}
		return o, nil
	}
}

// PropertyOrder is add propertyOrder to schema.
func PropertyOrder(order int) Option {
	return func(o Object) (Object, error) {
		o.Set("propertyOrder", order)
		return o, nil
	}
}

type refWrapper struct {
	obj Object
	ref string
}

func (o *refWrapper) Set(key string, value interface{}) {
	o.obj.Set(key, value)
}

func (o *refWrapper) Get(key string) (interface{}, bool) {
	return o.obj.Get(key)
}

func (o *refWrapper) Ref() string {
	return o.ref
}

// Ref replaces to given ref.
func Ref(ref string) Option {
	return func(o Object) (Object, error) {
		return &refWrapper{
			obj: o,
			ref: ref,
		}, nil
	}
}
