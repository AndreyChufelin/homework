package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/AndreyChufelin/homework/hw09_struct_validator/validators"
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
		tags := strings.Split(reflect.TypeOf(v).Field(i).Tag.Get("validate"), "|")

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

			switch value.Field(i).Kind() {
			case reflect.Struct:
				validateStruct(&errs, value.Field(i), key)
			case reflect.Slice:
				l := value.Field(i).Len()
				for j := 0; j < l; j++ {
					err := validateType(&errs, value.Field(i).Index(j), value.Type().Field(i).Name+"["+strconv.Itoa(j)+"]", key, val)
					if err != nil {
						return err
					}
				}
			default:
				err := validateType(&errs, value.Field(i), value.Type().Field(i).Name, key, val)
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
			*errs = append(*errs, ValidationError{name, err})
		}
	case reflect.Int:
		validatorInt, ok := validators.ValidatorsInt[key]
		if !ok {
			return ErrInvalidValidator
		}
		if err := validatorInt(int(value.Int()), val); err != nil {
			*errs = append(*errs, ValidationError{name, err})
		}
	default:
	}

	return nil
}

func validateStruct(errs *ValidationErrors, value reflect.Value, key string) {
	if key == "nested" && value.CanInterface() {
		err := Validate(value.Interface())
		target := ValidationErrors{}

		if errors.As(err, &target) {
			*errs = append(*errs, target...)
		}
	}
}
