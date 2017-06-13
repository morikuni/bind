package bind

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type Getter interface {
	Get(key string) []string
}

func FromRequest(r *http.Request, target interface{}) error {
	r.ParseForm()
	return FromValues(r.Form, target)
}

type mapGetter struct {
	m map[string]string
}

func (g mapGetter) Get(key string) []string {
	if v, ok := g.m[key]; ok {
		return []string{v}
	}
	return nil
}

func FromMap(m map[string]string, target interface{}) error {
	return FromGetter(mapGetter{m}, target)
}

type valuesGetter struct {
	v url.Values
}

func (g valuesGetter) Get(key string) []string {
	return g.v[key]
}

func FromValues(v url.Values, target interface{}) error {
	return FromGetter(valuesGetter{v}, target)
}

func FromGetter(getter Getter, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	if targetValue.IsNil() {
		return ErrNil
	}
	if targetValue.Elem().Kind() != reflect.Struct {
		return ErrNotStructPointer
	}

	targetType := reflect.TypeOf(target).Elem()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		tag, ok := field.Tag.Lookup("bind")
		if !ok {
			tag = field.Name
		}
		values := getter.Get(tag)
		target := targetValue.Elem().Field(i)
		if !target.CanSet() {
			continue
		}
		err := assignValue(values, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func assignValue(values []string, target reflect.Value) error {
	if len(values) == 0 {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	if target.Kind() == reflect.Ptr {
		if target.IsNil() {
			target.Set(reflect.New(target.Type().Elem()))
		}
		target = target.Elem()
	}

	value := values[0]
	switch target.Kind() {
	case reflect.Slice:
		l := len(values)
		s := reflect.MakeSlice(target.Type(), l, l)
		for i := 0; i < l; i++ {
			assignValue(values[i:], s.Index(i))
		}
		target.Set(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return RaiseConvertError(value, target.Type())
		}
		target.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if value == "" {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return RaiseConvertError(value, target.Type())
		}
		target.SetUint(n)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return RaiseConvertError(value, target.Type())
		}
		target.SetFloat(f)
	case reflect.String:
		if value == "" {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		target.SetString(value)
	case reflect.Bool:
		if value == "" {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		b, err := strconv.ParseBool(value)
		if err != nil {
			return RaiseConvertError(value, target.Type())
		}
		target.SetBool(b)
	default:
		return RaiseConvertError(value, target.Type())
	}

	return nil
}
