package cli

import "fmt"

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("exit status %d", e.Code)
}

func NewExitError(code int, format string, args ...any) error {
	return &ExitError{
		Code: code,
		Err:  fmt.Errorf(format, args...),
	}
}
