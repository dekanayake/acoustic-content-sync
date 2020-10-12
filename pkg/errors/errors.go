package errors

import "github.com/pkg/errors"

type retryableError struct {
	error
}

func ErrorMessageWithStack(message string) error {
	return errors.WithStack(errors.New(message))
}

func ErrorWithStack(err error) error {
	return errors.WithStack(err)
}

func RetryableError(err error) error {
	return &retryableError{
		err,
	}
}

func isRetryableError(err error, level int) bool {
	_, ok := err.(*retryableError)
	if !ok && level < 5 {
		unwrappedError := errors.Unwrap(err)
		if unwrappedError != nil {
			return isRetryableError(unwrappedError, level+1)
		}
	}
	return ok
}

func IsRetryableError(err error) bool {
	return isRetryableError(err, 1)
}
