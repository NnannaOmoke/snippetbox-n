package validator

import "strings"
import "regexp"
import "unicode/utf8"

const EmailRegexString = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var EmailRegex = regexp.MustCompile(EmailRegexString)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, msg string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}

func (v *Validator) AddNonFieldError(msg string) {
	v.NonFieldErrors = append(v.NonFieldErrors, msg)
}

func (v *Validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, limit int) bool {
	return utf8.RuneCountInString(value) <= limit
}

func PermittedVal[T comparable](value T, pValues ...T) bool {
	for i := range pValues {
		if value == pValues[i] {
			return true
		}
	}
	return false
}

func MinChars(value string, limit int) bool {
	return utf8.RuneCountInString(value) >= limit
}

func Matches(value string, regexp *regexp.Regexp) bool {
	return regexp.MatchString(value)
}

func FitsCategory(field string, fn func(string) bool) bool {
	return fn(field)
}
