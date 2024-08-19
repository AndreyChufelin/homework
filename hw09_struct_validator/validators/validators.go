package validators

import (
	"github.com/AndreyChufelin/homework/hw09_struct_validator/validators/intv"
	"github.com/AndreyChufelin/homework/hw09_struct_validator/validators/str"
)

var (
	ValidatorsStr = map[string]func(string, string) error{
		"len":    str.LenValidator,
		"regexp": str.RegexpValidator,
		"in":     str.InValidator,
	}
	ValidatorsInt = map[string]func(int, string) error{
		"max": intv.MaxValidator,
		"min": intv.MinValidator,
		"in":  intv.InValidator,
	}
)
