package errlib

import "fmt"

func WrapErr(err error, msg string) error {
	return fmt.Errorf(" %s: %w ", msg, err)
}
