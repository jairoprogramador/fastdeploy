package common

import (
	"errors"
	"strings"
)

type StringValueObject struct {
	value string
}

func NewStringValueObject(value string, fieldName string) (StringValueObject, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return StringValueObject{}, errors.New(fieldName + " cannot be empty")
	}

	safeValue := makeSafeForFileSystem(trimmedValue)

	return StringValueObject{value: safeValue}, nil
}

func (s StringValueObject) Value() string {
	return s.value
}

func (s StringValueObject) String() string {
	return s.value
}

func (s StringValueObject) Equals(other StringValueObject) bool {
	return s.value == other.value
}

func makeSafeForFileSystem(value string) string {
	unsafeChars := map[string]string{
		"<":  "_",
		">":  "_",
		":":  "_",
		"\"": "_",
		"/":  "_",
		"\\": "_",
		"|":  "_",
		"?":  "_",
		"*":  "_",
	}

	result := value
	for unsafe, safe := range unsafeChars {
		result = strings.ReplaceAll(result, unsafe, safe)
	}

	return result
}
