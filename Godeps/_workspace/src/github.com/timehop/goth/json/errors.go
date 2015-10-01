package json

import "fmt"

type MissingFieldError string

func (s MissingFieldError) Error() string {
	return fmt.Sprintf("goth/json: missing field '%v'", string(s))
}

func (s MissingFieldError) IsTransient() bool { return false }

func IsMissingFieldError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(MissingFieldError)
	return ok
}

func NewTypeError(field, expectedType string, actualValue interface{}) error {
	return &TypeError{field, expectedType, fmt.Sprintf("%T", actualValue)}
}

type TypeError struct {
	field    string
	expected string
	actual   string
}

func (s TypeError) Error() string {
	return fmt.Sprintf("goth/json: expected field '%v' to be of type %v but was %v", s.field, s.expected, s.actual)
}

func (s TypeError) IsTransient() bool { return false }
