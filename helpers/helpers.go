package helpers

import (
	"strings"
	"ForumDatabase/errors"
)

var (
	MinLengthTitle = 10
	MinLengthContent = 16
	MinLengthPassword = 8
	MinLengthUsername = 6
)

// TODO: Probably should add an NG check for these methods

func ValidateTitle(input string) *errors.UserError {
	trimmed := strings.Trim(input, " ")
	if len(trimmed) < MinLengthTitle {
		return errors.ErrTooShort
	} else {
		return nil
	}
}

func ValidateContent(input string) *errors.UserError {
	trimmed := strings.Trim(input, " ")
	if len(trimmed) < MinLengthContent {
		return errors.ErrTooShort
	} else {
		return nil
	}
}

func ValidatePassword(input string) *errors.UserError {

	if strings.Contains(input, " ") {
		return errors.ErrContainSpaces
	} else if len(input) < MinLengthPassword {
		return errors.ErrTooShort
	} else {
		return nil
	}

}

func ValidateUsername(input string) *errors.UserError {

	if strings.Contains(input, " ") {
		return errors.ErrContainSpaces
	} else if len(input) < MinLengthUsername {
		return errors.ErrTooShort
	} else {
		return nil
	}

}
