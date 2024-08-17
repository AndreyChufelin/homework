package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/AndreyChufelin/homework/hw09_struct_validator/validators"
	"github.com/AndreyChufelin/homework/hw09_struct_validator/validators/intv"
	"github.com/AndreyChufelin/homework/hw09_struct_validator/validators/str"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result strings.Builder

	l := len(v)
	for i, err := range v {
		result.WriteString(err.Error())
		if i < l-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

var ErrInvalidValidator = errors.New("invalid validator name")

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	errs := ValidationErrors{}

	if value.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := value.Type().Field(i)
		tags := strings.Split(structField.Tag.Get("validate"), "|")

		for _, tag := range tags {
			if tag == "" {
				continue
			}
			s := strings.Split(tag, ":")
			key := s[0]

			val := ""
			if len(s) >= 2 {
				val = s[1]
			}

			switch field.Kind() {
			case reflect.Struct:
				errs = append(errs, validateStruct(field, key)...)
			case reflect.Slice:
				err := validateSlice(&errs, field, structField, key, val)
				if err != nil {
					return err
				}
			default:
				err := validateType(&errs, field, structField.Name, key, val)
				if err != nil {
					return err
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateType(errs *ValidationErrors, value reflect.Value, name string, key string, val string) error {
	switch value.Kind() {
	case reflect.String:
		validator, ok := validators.ValidatorsStr[key]
		if !ok {
			return ErrInvalidValidator
		}
		if err := validator(value.String(), val); err != nil {
			if errors.Is(err, str.ErrInvalidValue) {
				return fmt.Errorf("%s: %w", name, err)
			}
			*errs = append(*errs, ValidationError{name, err})
		}
	case reflect.Int:
		validatorInt, ok := validators.ValidatorsInt[key]
		if !ok {
			return ErrInvalidValidator
		}
		if err := validatorInt(int(value.Int()), val); err != nil {
			if errors.Is(err, intv.ErrInvalidValue) {
				return fmt.Errorf("%s: %w", name, err)
			}
			*errs = append(*errs, ValidationError{name, err})
		}
	}

	return nil
}

func validateStruct(value reflect.Value, key string) ValidationErrors {
	var errs ValidationErrors

	if key == "nested" && value.CanInterface() {
		err := Validate(value.Interface())
		target := ValidationErrors{}

		if errors.As(err, &target) {
			errs = append(errs, target...)
		}
	}

	return errs
}

func validateSlice(
	errs *ValidationErrors, field reflect.Value, structField reflect.StructField, key string, val string,
) error {
	l := field.Len()
	for j := 0; j < l; j++ {
		err := validateType(errs, field.Index(j), structField.Name+"["+strconv.Itoa(j)+"]", key, val)
		if err != nil {
			return err
		}
	}

	return nil
}
