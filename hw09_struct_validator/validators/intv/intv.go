package intv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrMinTooBigValue = errors.New("value is too small")

func MinValidator(i int, val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("invalid value for min validator")
	}

	if i < v {
		return fmt.Errorf("value is too small")
	}

	return nil
}

var ErrMaxTooBigValue = errors.New("value is too big")

func MaxValidator(i int, val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("invalid value for min validator")
	}

	if i > v {
		return fmt.Errorf("value is too big")
	}

	return nil
}

var ErrInNotContainValue = errors.New("value doesn't exist in set")

func InValidator(i int, val string) error {
	for _, v := range strings.Split(val, ",") {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid value for in validator")
		}

		if i == n {
			return nil
		}
	}

	return fmt.Errorf("%w %s", ErrInNotContainValue, val)
}
