package auth

import "fmt"

type AuthError struct {
	Reason string
	Inner  error
}

func (e AuthError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("auth error: %s: %v", e.Reason, e.Inner)
	}
	return fmt.Sprintf("auth error: %s", e.Reason)
}

func (e AuthError) Unwrap() error {
	return e.Inner
}
