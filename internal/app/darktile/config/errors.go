package config

import (
	"errors"
	"fmt"
)

type RecoverableError struct {
	msg   string
	inner error
}

func NewRecoverableError(msg string, cause error) *RecoverableError {
	return &RecoverableError{
		inner: cause,
		msg:   msg,
	}
}

func IsErrRecoverable(err error) bool {
	var rec *RecoverableError
	return errors.As(err, &rec)
}

func (e *RecoverableError) Error() string {
	if e.inner == nil {
		return e.msg
	}

	return fmt.Sprintf("%s: %s", e.msg, e.inner)
}
