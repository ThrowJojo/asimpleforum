package errors

import "errors"

type UserError struct {
	Err error
	Code int
}

var (
	ErrNotExist = &UserError{errors.New("Record does not exist"), 1}
	ErrBadRecord = &UserError{errors.New("Bad record"), 2}
	ErrSystem = &UserError{errors.New("System error"), 3}
	ErrExists = &UserError{errors.New("Record already exists"), 4}
	ErrTooShort = &UserError{errors.New("Input too short"), 5}
	ErrContainSpaces = &UserError{errors.New("Input contains spaces"), 6}
)

func (msg *UserError) Error() string {
	return msg.Err.Error()
}
