package detek

import "fmt"

func NewError(cause error, reason ErrorType) error {
	return &DetekError{cause: cause, reason: reason}
}

type ErrorType string

const (
	ErrNotEnoughConfig ErrorType = "not enough config provided for producer"
	ErrKeyNotFound     ErrorType = "requested key not found"
)

type DetekError struct {
	cause  error
	reason ErrorType
}

func (d *DetekError) Error() string {
	if d.cause != nil {
		return fmt.Sprintf("detek failed: %s: %s", d.reason, d.cause.Error())
	} else {
		return fmt.Sprintf("detek failed: %s", d.reason)
	}
}

func (d *DetekError) Unwrap() error {
	return d.cause
}
