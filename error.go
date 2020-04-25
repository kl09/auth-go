package auth

import (
	"errors"
	"fmt"
)

const (
	// ErrInternal is an internal error.
	ErrInternal = "internal"
	// ErrCredNotFound is returned when credential not found.
	ErrCredNotFound = "credential_not_found"
	// ErrAuth is returned when auth is failed.
	ErrAuth = "auth_failed"
)

// Error represents an error within the context of Quoter service.
type Error struct {
	// Code is a machine-readable code.
	Code string `json:"code"`
	// Message is a human-readable message.
	Message string `json:"message"`
	// err is a previous error in error chain.
	err error
}

// Error returns the string representation of the error message.
func (e Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.err.Error())
	}

	return fmt.Sprintf("%s %s", e.Code, e.Message)
}

// Unwrap the error returning the error's reason.
func (e Error) Unwrap() error {
	return e.err
}

// ErrorCode returns the code of the error, if available.
func ErrorCode(err error) string {
	var e Error
	if errors.As(err, &e) {
		return e.Code
	}

	return "internal"
}

// ErrorMsg returns the message of the error, if available.
func ErrorMsg(err error) string {
	var e Error
	if errors.As(err, &e) {
		return e.Message
	}

	return "no message"
}

// NewError creates a new Error instance using provided code and message.
func NewError(code, message string) error {
	return Error{
		Code:    code,
		Message: message,
	}
}

// ErrorHas checks if one of error codes exist in errors chain and returns the first error that has one of the codes provided.
func ErrorHas(err error, codes ...string) error {
	var e Error
	if errors.As(err, &e) {
		for _, code := range codes {
			if e.Code == code {
				return e
			}
		}

		return ErrorHas(e.Unwrap(), codes...)
	}

	return nil
}

// WrapError wraps err into a new error with provided code and message.
// It will return nil if error is nil.
func WrapError(err error, code, message string) error {
	if err == nil {
		return nil
	}

	return Error{
		Code:    code,
		Message: message,
		err:     err,
	}
}
