package bind

import (
	"errors"
	"fmt"
	"reflect"
)

type Error interface {
	error

	bindError()
}

type bindError struct {
	error
}

func (e bindError) bindError() {}

func ErrorOf(message string) error {
	return bindError{errors.New(message)}
}

var (
	ErrNil              error = ErrorOf("target is nil")
	ErrNotPointer       error = ErrorOf("target must be a pointer")
	ErrNotStructPointer error = ErrorOf("target must be a struct")
)

type ConvertError struct {
	Value string
	To    reflect.Type
}

func (e ConvertError) Error() string {
	return fmt.Sprintf("cannot convert %q to %s", e.Value, e.To)
}

func (e ConvertError) bindError() {}

func RaiseConvertError(value string, to reflect.Type) error {
	return ConvertError{value, to}
}
