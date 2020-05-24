package validation

import (
	"errors"
	"regexp"
)

// ValidateEmail is to validate email format
func ValidateEmail(email string) (err error) {
	validEmail := regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	if !validEmail.MatchString(email) {
		err = errors.New("Invalid email format")
	}
	return
}
