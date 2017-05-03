package helpers

import (
	"github.com/pkg/errors"
	"strings"
)

var (
	MinLengthTitle = 10
	MinLengthContent = 16
	MinLengthPassword = 8
	MinLengthUsername = 6
	ErrTooShort = errors.New("Input too short")
	ErrContainSpaces = errors.New("Input contains spaces")
)

// TODO: Probably should add an NG check for these methods

func ValidateTitle(input string) error {
	trimmed := strings.Trim(input, " ")
	if len(trimmed) < MinLengthTitle {
		return ErrTooShort
	} else {
		return nil
	}
}

func ValidateContent(input string) error {
	trimmed := strings.Trim(input, " ")
	if len(trimmed) < MinLengthContent {
		return ErrTooShort
	} else {
		return nil
	}
}

func ValidatePassword(input string) error {

	if strings.Contains(input, " ") {
		return ErrContainSpaces
	} else if len(input) < MinLengthPassword {
		return ErrTooShort
	} else {
		return nil
	}

}

func ValidateUsername(input string) error {

	if strings.Contains(input, " ") {
		return ErrContainSpaces
	} else if len(input) < MinLengthUsername {
		return ErrTooShort
	} else {
		return nil
	}

}
