package internal

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	ErrorCodeUnknown ErrorCode = iota
	// Define other error codes here
)

func WrapErrorf(err error, code ErrorCode, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

func NewErrorf(code ErrorCode, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), errors.New(fmt.Sprintf(format, args...)))
}
