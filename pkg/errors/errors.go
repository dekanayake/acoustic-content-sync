package errors

import "github.com/pkg/errors"

func ErrorMessageWithStack(message string) error {
	return errors.WithStack(errors.New(message))
}

func ErrorWithStack(err error) error {
	return errors.WithStack(err)
}
