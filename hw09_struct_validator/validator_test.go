package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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

	Request struct {
		IDs []int `validate:"in:1,2,3"`
	}

	UnsupportedTypes struct {
		Price float64 `validate:"max:50"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "1111",
				Name:   "Vasya Pupkin",
				Age:    25,
				Email:  "vasya@test.com",
				Role:   "admin",
				Phones: []string{"+79111111111"},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.10.0",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errors.New("length is greater than 5"),
				},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 600,
				Body: "OK",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errors.New("value 600 is not in [200 404 500]"),
				},
			},
		},
		{
			in: Request{
				IDs: []int{1, 2, 5},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "IDs",
					Err:   errors.New("value 5 is not in [1 2 3]"),
				},
			},
		},
		{
			in:          "test",
			expectedErr: ErrorNotStruct,
		},
		{
			in: UnsupportedTypes{
				Price: 64.59,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Price",
					Err:   ErrorUnsupportedType,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
