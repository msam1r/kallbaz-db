package db

import (
	"fmt"
)

type Store interface {
	// Get - returns the value of the given key or any error occurs. if the
	// key was not found it will return a NotFoundError.
	Get(key string) ([]byte, error)
	// Put - stores the value. it will return a BadRequestError if the provided
	// data was invalid or any other error occured.
	Put(key string, value []byte) error
	// Delete - deletes the value of the given key.
	Delete(key string) error
	// Close - closes the database and returns when all internal processes
	// has stopped. it returns any error occurs.
	Close() error
	// IsNotFoundError - check if the error is of type NotFoundError.
	IsNotFoundError(err error) bool
	// IsBadRequestError - check if the error is of type BadRequestError.
	IsBadRequestError(err error) bool
}

// NotFoundError - indicates that no value was found for the given key.
type NotFoundError struct {
	error
	missingKey string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find value for key: %s", e.missingKey)
}

func NewNotFoundError(key string) error {
	return &NotFoundError{missingKey: key}
}

// BadRequestError - represents an error by the consumer of the database.
type BadRequestError struct {
	error
	message string
}

func (e *BadRequestError) Error() string {
	return e.message
}

func NewBadRequestError(message string) error {
	return &BadRequestError{message: message}
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func IsBadRequestError(err error) bool {
	_, ok := err.(*BadRequestError)
	return ok
}
