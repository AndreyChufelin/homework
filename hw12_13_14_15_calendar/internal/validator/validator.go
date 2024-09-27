package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v Validator) AddError(title, msg string) {
	v.Errors[title] = msg
}
func (v Validator) Check(condition bool, title, msg string) {
	if !condition {
		return
	}
	v.Errors[title] = msg
}
