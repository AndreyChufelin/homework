package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/AndreyChufelin/homework/hw09_struct_validator/validators/str"
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

	Nested struct {
		User User `validate:"nested"`
		app  App  `validate:"nested"`
	}
)

type validateTests struct {
	in          interface{}
	expectedErr error
}

func TestValidate(t *testing.T) {
	tests := []validateTests{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    18,
				Email:  "john@go.dev",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   json.RawMessage(`{"foo":"bar"}`),
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte(`{"alg":"HS256"}`),
				Payload:   []byte(`{"foo":"bar"}`),
				Signature: []byte(`foobar`),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "ok",
			},
			expectedErr: nil,
		},
		{
			in: Nested{
				User: User{
					ID:     "123456789012345678901234567890123456",
					Name:   "John",
					Age:    18,
					Email:  "john@go.dev",
					Role:   "admin",
					Phones: []string{"12345678901"},
					meta:   json.RawMessage(`{"foo":"bar"}`),
				},
				app: App{Version: "1.0.0"},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.NoError(t, err, tt.expectedErr)
		})
	}
}

type wrongValue struct {
	Name string `validate:"len:wrong"`
}

func TestValidateValidatorErrors(t *testing.T) {
	tests := []validateTests{
		{
			in: User{
				ID:     "1234567890123456789012345678901234562",
				Name:   "John",
				Age:    55,
				Email:  "johngo.dev",
				Role:   "user",
				Phones: []string{"1234567890"},
				meta:   json.RawMessage(`{"foo":"bar"}`),
			},
			expectedErr: fmt.Errorf("ID: value is too long\nAge: value is too big\nEmail: value doesn't match regexp\n" +
				"Role: value doesn't exist in set admin,stuff\nPhones[0]: value is too short"),
		},
		{
			in: Nested{
				User: User{
					ID:     "1234567890123456789012345678901234562",
					Name:   "John",
					Age:    55,
					Email:  "johngo.dev",
					Role:   "user",
					Phones: []string{"1234567890"},
					meta:   nil,
				},
				app: App{Version: "1.0"},
			},
			expectedErr: fmt.Errorf("ID: value is too long\nAge: value is too big\nEmail: value doesn't match regexp\n" +
				"Role: value doesn't exist in set admin,stuff\nPhones[0]: value is too short"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

type wrongStruct struct {
	Name string `validate:"wrong:validator"`
}

func TestValidateInternalErrors(t *testing.T) {
	tests := []validateTests{
		{
			in: wrongStruct{
				Name: "John",
			},
			expectedErr: ErrInvalidValidator,
		},
		{
			in: wrongValue{
				Name: "John",
			},
			expectedErr: str.ErrInvalidValue,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
