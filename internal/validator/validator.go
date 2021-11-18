package validator

import (
	"regexp"
)

var (
	EmailRx = regexp.MustCompile( "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\. [a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)
// map of validation errors
type Validator struct {
	Errors map[string]string
}
// helper which creates a new Validator instance with an empty errors map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}
// returns true if the errors maps doesnt contain entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}
// adds an error message to the map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}
// adds an error message to the map only if a validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
// returns true if a specific value is in a list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}
// returns true if string value matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
// returns true if all string values in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
