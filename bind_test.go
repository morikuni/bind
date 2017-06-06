package bind

import (
	"testing"

	"github.com/morikuni/pointer"
	"github.com/stretchr/testify/assert"
)

type Tester struct {
	Int    int     `bind:"int"`
	Uint   uint    `bind:"uint"`
	Float  float64 `bind:"float"`
	String string  `bind:"string"`
	Bool   bool    `bind:"bool"`

	IntP    *int
	UintP   *uint
	FloatP  *float64
	StringP *string
	BoolP   *bool

	unexported *uintptr `bind:"unexported"`

	Dummy map[string]string
}

func TestFromGetter(t *testing.T) {
	type TestCase struct {
		Description string
		Input       map[string]string
		Expect      Tester
	}

	table := []TestCase{
		{
			Description: "full",
			Input: map[string]string{
				"int":    "123",
				"uint":   "123",
				"float":  "123.4",
				"string": "hello",
				"bool":   "true",

				"IntP":    "-123",
				"UintP":   "123",
				"FloatP":  "-123.4",
				"StringP": "hello",
				"BoolP":   "false",

				"unexported": "123",
			},
			Expect: Tester{
				Int:    123,
				Uint:   123,
				Float:  123.4,
				String: "hello",
				Bool:   true,

				IntP:    pointer.Int(-123),
				UintP:   pointer.Uint(123),
				FloatP:  pointer.Float64(-123.4),
				StringP: pointer.String("hello"),
				BoolP:   pointer.Bool(false),

				unexported: nil,
			},
		},
		{
			Description: "no input",
			Input:       map[string]string{},
			Expect: Tester{
				Int:    0,
				Uint:   0,
				Float:  0,
				String: "",
				Bool:   false,

				IntP:    nil,
				UintP:   nil,
				FloatP:  nil,
				StringP: nil,
				BoolP:   nil,

				unexported: nil,
			},
		},
		{
			Description: "empty string",
			Input: map[string]string{
				"int":    "",
				"uint":   "",
				"float":  "",
				"string": "",
				"bool":   "",

				"IntP":    "",
				"UintP":   "",
				"FloatP":  "",
				"StringP": "",
				"BoolP":   "",

				"unexported": "",
			},
			Expect: Tester{
				Int:    0,
				Uint:   0,
				Float:  0,
				String: "",
				Bool:   false,

				IntP:    pointer.Int(0),
				UintP:   pointer.Uint(0),
				FloatP:  pointer.Float64(0),
				StringP: pointer.String(""),
				BoolP:   pointer.Bool(false),

				unexported: nil,
			},
		},
	}

	for _, test := range table {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)

			tester := Tester{}
			err := FromMap(test.Input, &tester)
			assert.NoError(err)
			assert.Equal(test.Expect, tester)
		})
	}
}

func TestInvalidTargetTypes(t *testing.T) {
	type TestCase struct {
		Description string
		Input       interface{}
		Expect      error
	}

	table := []TestCase{
		{
			Description: "not pointer",
			Input:       Tester{},
			Expect:      ErrNotPointer,
		},
		{
			Description: "nil",
			Input:       (*Tester)(nil),
			Expect:      ErrNil,
		},
		{
			Description: "int",
			Input:       1,
			Expect:      ErrNotPointer,
		},
		{
			Description: "slice",
			Input:       []string{},
			Expect:      ErrNotPointer,
		},
		{
			Description: "chan",
			Input:       make(chan float32),
			Expect:      ErrNotPointer,
		},
		{
			Description: "non struct pointer",
			Input:       pointer.Bool(true),
			Expect:      ErrNotStructPointer,
		},
	}

	for _, test := range table {
		t.Run(test.Description, func(t *testing.T) {
			assert := assert.New(t)

			err := FromMap(map[string]string{}, test.Input)
			assert.Equal(test.Expect, err)
		})
	}
}
