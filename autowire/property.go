package autowire

import (
	"fmt"
	"reflect"

	"github.com/modern-go/reflect2"
	"github.com/non1996/go-jsonobj/constraint"
	util "github.com/non1996/go-jsonobj/utils"
)

type propertyProvider struct {
	scope    string
	instance any
	provide  func() any
}

type properties struct {
	m      map[string]any // scope + "/" + key -> value
	scopes map[string]propertyProvider
}

func newProperties() properties {
	return properties{
		m:      map[string]any{},
		scopes: map[string]propertyProvider{},
	}
}

func (p *properties) add(provider propertyProvider) {
	p.scopes[provider.scope] = provider
}

func (p *properties) set(scope, key string, value any) {
	p.m[fmt.Sprintf("%s/%s", scope, key)] = value
}

func (p *properties) get(scope, key string) (any, bool) {
	provider, exist := p.scopes[scope]
	if !exist {
		return nil, false
	}

	if reflect2.IsNil(provider.instance) {
		provider.instance = provider.provide()
		kvs := objToKvPairs(provider.instance)

		for _, kv := range kvs {
			p.set(scope, kv.First, kv.Second)
		}
	}

	value, exist := p.m[fmt.Sprintf("%s/%s", scope, key)]
	return value, exist
}

func valueAsInt[I constraint.Int | constraint.Uint](v any) I {
	value := reflect.ValueOf(v)

	switch value.Type().Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return I(value.Int())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return I(value.Uint())
	case reflect.Float32,
		reflect.Float64:
		return I(value.Float())
	default:
		panic("")
	}
}

func valueAsFloat[F constraint.Float](v any) F {
	value := reflect.ValueOf(v)

	switch value.Type().Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return F(value.Int())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return F(value.Uint())
	case reflect.Float32,
		reflect.Float64:
		return F(value.Float())
	default:
		panic("")
	}
}

func valueAsBool[B constraint.Bool](v any) B {
	value := reflect.ValueOf(v)

	switch value.Type().Kind() {
	case reflect.Bool:
		return B(value.Bool())
	default:
		panic("")
	}
}

func valueAsString[S constraint.String](v any) S {
	value := reflect.ValueOf(v)
	switch value.Type().Kind() {
	case reflect.String:
		return S(value.String())
	default:
		panic("")
	}
}

func valueAsSlice[S ~[]E, E any](v any) S {
	value := reflect.ValueOf(v)

	switch value.Type().Kind() {
	case reflect.Slice, reflect.Array:
		elemType1 := value.Type().Elem()
		elemType2 := TypeOf[E]()

		if elemType1 != elemType2 {
			panic("")
		}

		return S(value.Interface().([]E))
	default:
		panic("")
	}
}

func objToKvPairs(obj any) (res []util.Pair[string, any]) {
	v := deRefValue(reflect.ValueOf(obj))
	if v.Kind() != reflect.Struct {
		panic("invalid object")
	}

	return objToKvImpl("", v)
}

func objToKvImpl(prefix string, v reflect.Value) (res []util.Pair[string, any]) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		if fieldValue.Type().Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				continue
			}

			fieldValue = deRefValue(fieldValue)
		}

		if fieldValue.Type().Kind() == reflect.Struct {
			if fieldType.Anonymous {
				res = append(res, objToKvImpl(joinKey(prefix, fieldValue.Type().Name()), fieldValue)...)
			} else {
				res = append(res, objToKvImpl(joinKey(prefix, fieldType.Name), fieldValue)...)
			}
		} else if validConfigFieldKind(fieldValue.Type().Kind()) {
			if fieldType.Anonymous {
				res = append(res, util.NewPair(joinKey(prefix, fieldValue.Type().Name()), fieldValue.Interface()))
			} else {
				res = append(res, util.NewPair(joinKey(prefix, fieldType.Name), fieldValue.Interface()))
			}
		}
	}

	return res
}

func joinKey(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

func deRefValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return v
}

func validConfigFieldKind(kind reflect.Kind) bool {
	if kind > reflect.Invalid && kind < reflect.Complex64 {
		return true
	}
	return kind == reflect.String || kind == reflect.Slice || kind == reflect.Array
}
