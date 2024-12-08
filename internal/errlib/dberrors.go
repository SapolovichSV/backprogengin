package errlib

import "github.com/jackc/pgx/v5"

type NotFoundErr struct {
	Where string
	What  string
}

func (e NotFoundErr) Error() string {
	return e.What + " not found in " + e.Where
}

type UnexpectedErr struct {
	Where string
	Err   error
}

func (e UnexpectedErr) Error() string {
	return "unexpected error in " + e.Where + " : " + e.Err.Error()
}
func CheckErrNotFoundInDB(err error) bool {
	if err == nil {
		return false
	}
	if err == pgx.ErrNoRows {
		return true
	}
	return false
}
func CheckErrUnexpectedInDB(err error) bool {
	return err != nil
}
func WrapError(err error, where string, what string) error {
	if CheckErrNotFoundInDB(err) {
		return NotFoundErr{Where: where, What: what}
	}
	if CheckErrUnexpectedInDB(err) {
		return UnexpectedErr{Where: where, Err: err}
	}
	return err
}
