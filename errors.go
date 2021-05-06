package QuickHTTPServerLauncher

import (
	cer "github.com/kaatinga/const-errs"
)

const (
	errValidationError cer.Error = "validation failed"
)

type validationError string

func (err validationError) Error() string {
	return string(err)
}

func (err validationError) Is(target error) bool {
	return target == errValidationError //nolint:errorlint
}
