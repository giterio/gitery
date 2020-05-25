package validation

import (
	"errors"
	"regexp"
)

// ValidatePassword is to validate the password which should be a combination of uppercase and lowercase letters and numbers with a length of 8-32
func ValidatePassword(password string) (err error) {
	if len(password) < 8 {
		err = errors.New("Password must be 8 characters long")
		return
	}
	hasNumbers := regexp.MustCompile(`\d`)
	if !hasNumbers.MatchString(password) {
		err = errors.New("Password must have at least one number")
		return
	}
	hasLowerCase := regexp.MustCompile(`[a-z]`)
	if !hasLowerCase.MatchString(password) {
		err = errors.New("Password must have at least one lowercase character")
		return
	}
	hasUpperCase := regexp.MustCompile(`[A-Z]`)
	if !hasUpperCase.MatchString(password) {
		err = errors.New("Password must have at least one uppercase character")
		return
	}
	return
}
