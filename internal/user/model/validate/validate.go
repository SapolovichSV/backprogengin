package validate

type ValidateError struct {
	What string
}

func (e ValidateError) Error() string {
	return e.What + " is invalid"
}
func UserName(username string) error {
	if len(username) < 4 {
		return ValidateError{What: "username"}
	}
	return nil
}
func VPassword(password string) error {
	if len(password) < 4 {
		return ValidateError{What: "password"}
	}
	return nil
}
