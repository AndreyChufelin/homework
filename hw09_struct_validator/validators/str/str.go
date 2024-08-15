package str

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLenIsTooLong      = errors.New("value is too long")
	ErrLenIsTooShort     = errors.New("value is too short")
	ErrRegexpNoMatch     = errors.New("value doesn't match regexp")
	ErrInNotContainValue = errors.New("value doesn't exist in set")
)

func LenValidator(s, val string) error {
	v, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("invalid value for length validator")
	}

	if len(s) > v {
		return ErrLenIsTooLong
	} else if len(s) < v {
		return ErrLenIsTooShort
	}

	return nil
}

func RegexpValidator(s, val string) error {
	r, err := regexp.Compile(val)
	if err != nil {
		return fmt.Errorf("invalid value for regexp validator")
	}

	if !r.MatchString(s) {
		return ErrRegexpNoMatch
	}

	return nil
}

func InValidator(s, val string) error {
	for _, v := range strings.Split(val, ",") {
		if s == v {
			return nil
		}
	}

	return fmt.Errorf("%w %s", ErrInNotContainValue, val)
}
