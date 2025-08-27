package valueObjects

import (
	"errors"
	"strings"
)

var unsafeChars = map[string]string{
	" ":  "_",
	"<":  "_",
	">":  "_",
	":":  "_",
	"\"": "_",
	"/":  "_",
	"\\": "_",
	"|":  "_",
	"?":  "_",
	"*":  "_",
	"%":  "_",
	"&":  "_",
	"+":  "_",
	"-":  "_",
	"=":  "_",
	"@":  "_",
	"#":  "_",
	"$":  "_",
	"!":  "_",
	"^":  "_",
	"~":  "_",
	"`":  "_",
	"'":  "_",
	"(":  "_",
	")":  "_",
	"{":  "_",
	"}":  "_",
	"[":  "_",
	"]":  "_",
}

type BaseString struct {
	value string
}

func NewBaseString(value string, fieldName string) (BaseString, error) {
	if value == "" {
		return BaseString{}, errors.New(fieldName + " cannot be empty")
	}
	return BaseString{value: value}, nil
}

func NewBaseStringEmpty() BaseString{
	return BaseString{value: ""}
}

func (s BaseString) Value() string {
	return s.value
}

func (s BaseString) String() string {
	return s.value
}

func (s BaseString) Equals(other BaseString) bool {
	return s.value == other.value
}

func (s BaseString) IsEmpty() bool {
	return s.value == ""
}

func (s BaseString) IsValid() bool {
	return !s.IsEmpty() && !strings.ContainsFunc(s.value, func(r rune) bool {
		_, exists := unsafeChars[string(r)]
		return exists
	})
}

func (s *BaseString) MakeSafe() BaseString {
	safeValue := MakeSafeForFileSystem(s.value)
	return BaseString{value: safeValue}
}

func MakeSafeForFileSystem(value string) string {
	trimmedValue := strings.TrimSpace(value)

	result := trimmedValue
	for unsafe, safe := range unsafeChars {
		result = strings.ReplaceAll(result, unsafe, safe)
	}

	return result
}
