package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	TestMinMax struct {
		IntVal  int    `validate:"min:10|max:20"`
		UintVal uint   `validate:"min:10|max:20"`
		StrVal  string `validate:"min:10|max:20"`
		StrVal2 string `validate:"min:10|max:20"`
	}

	TestInterface struct {
		IntVal  interface{} `validate:"min:10|max:20"`
		UintVal interface{} `validate:"min:10|max:20"`
		StrVal  interface{} `validate:"min:10|max:20"`
	}
	TestUnsupported struct {
		FloatVal float32 `validate:"min:10|max:20"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in           interface{}
		expectedErrs []error
	}{
		{ // invalid
			App{
				Version: "01.15.02",
			},
			[]error{ErrStringLength},
		},
		{ // valid
			App{
				Version: "01.00",
			},
			[]error{},
		},
		{ // invalid
			User{
				ID:     "4321343JH4KJ3H4K2341JK14324787236458",
				Name:   "Jon",
				Age:    35,
				Email:  "test@mail.net",
				Role:   "admin",
				Phones: []string{"12315478545", "12345678901"},
				meta:   nil,
			},
			[]error{},
		},
		{ // valid
			User{
				ID:     "4321343JH4KJ3H4K2341JK1432",
				Name:   "Jon",
				Age:    15,
				Email:  "test@mail",
				Role:   "user",
				Phones: []string{"12345", "12345678901"},
				meta:   nil,
			},
			[]error{ErrStringLength, ErrStringNotInSet, ErrStringContent, ErrNumberValue},
		},
		{ // ignore
			Token{
				Header:    []byte("12345"),
				Payload:   []byte("6576"),
				Signature: []byte("8878"),
			},
			[]error{},
		},
		{ // invalid
			Response{
				Code: 400,
				Body: "test",
			},
			[]error{ErrNumberNotInSet},
		},
		{ // valid
			Response{
				Code: 200,
				Body: "test",
			},
			[]error{},
		},
		{ // invalid
			TestMinMax{
				IntVal:  400,
				UintVal: 1000,
				StrVal:  "test",
				StrVal2: "hkjhfkjsdfkjdfjkasdhfjksdhfkjsadhfskjd",
			},
			[]error{ErrNumberValue, ErrStringLength},
		},
		{ // valid
			TestMinMax{
				IntVal:  15,
				UintVal: 15,
				StrVal:  "test string",
				StrVal2: "test string",
			},
			[]error{},
		},
		{ // invalid
			TestInterface{
				IntVal:  400,
				UintVal: 1000,
				StrVal:  "test",
			},
			[]error{ErrNumberValue, ErrStringLength},
		},
		{ // valid
			TestInterface{
				IntVal:  15,
				UintVal: 15,
				StrVal:  "test string",
			},
			[]error{},
		},
		{ // unsupported
			TestUnsupported{
				FloatVal: 15.23,
			},
			[]error{ErrNotValidatebleType},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			if len(tt.expectedErrs) > 0 {
				for _, tErr := range tt.expectedErrs {
					require.ErrorIs(t, err, tErr)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
