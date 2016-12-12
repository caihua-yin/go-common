package api


import (
	"fmt"
)

// Error is base type to report errors from HTTP requests
type Error struct {
	Code         int
	ExtendedCode int
	Message      string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d (%d): %s", e.Code, e.ExtendedCode, e.Message)
}

// RawError is an error that should passthrough its contents back
type RawError struct {
	Code        int
	ContentType string
	Body        []byte
	Headers     map[string]string
}

func (e *RawError) Error() string {
	return fmt.Sprintf("%d [%s]: %s", e.Code, e.ContentType, string(e.Body))
}

// Check panics with appropriate error
func Check(err error) {
	if err == nil {
		return
	}

	panic(err)
}
