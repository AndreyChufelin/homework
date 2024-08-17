package intv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrMinTooSmallValue  = errors.New("value is too small")
	ErrInNotContainValue = errors.New("value doesn't exist in set")
	ErrMaxTooBigValue    = errors.New("value is too big")
	ErrInvalidValue      = errors.New("invalid value")
)

func MinValidator(i int, val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("min validator: %w", ErrInvalidValue)
	}

	if i < v {
		return ErrMinTooSmallValue
	}

	return nil
}

func MaxValidator(i int, val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("max validator: %w", ErrInvalidValue)
	}

	if i > v {
		return ErrMaxTooBigValue
	}

	return nil
}

func InValidator(i int, val string) error {
	for _, v := range strings.Split(val, ",") {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("in validator: %w", ErrInvalidValue)
		}

		if i == n {
			return nil
		}
	}

	return fmt.Errorf("%w %s", ErrInNotContainValue, val)
}
