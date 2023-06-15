package errors

import "fmt"

type NotFoundError struct {
	message string
}

func NewNotFoundError(msg string, args ...any) NotFoundError {
	return NotFoundError{
		message: fmt.Sprintf(msg, args...),
	}
}

func (n NotFoundError) Error() string {
	return n.message
}

func IsNotFound(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}
