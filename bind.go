package bind

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type Getter interface {
	Get(key string) []string
}

type requestGetter struct {
	*http.Request
}

func (g requestGetter) Get(key string) []string {
	return g.Form[key]
}

func FromRequest(r *http.Request, target interface{}) error {
	return FromGetter(requestGetter{r}, target)
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

func FromGetter(getter Getter, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return errors.New("target is not pointer or nil")
	}

	targetType := reflect.TypeOf(target).Elem()
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		tag, ok := field.Tag.Lookup("bind")
		if !ok {
			tag = field.Name
		}
		values := getter.Get(tag)
		err := assignValue(values, targetValue.Elem().Field(i))
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
	if value == "" {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		target.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		target.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		target.SetFloat(f)
	case reflect.String:
		target.SetString(value)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		target.SetBool(b)
	default:
		return fmt.Errorf("cannot set to %s", target.Kind().String())
	}

	return nil
}
