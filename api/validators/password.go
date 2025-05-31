package validators

import (
	"github.com/go-playground/validator/v10"
)

// all takes any number of boolean values and returns true if and only if all values are true.
func all(cnds ...bool) bool {
	for _, cond := range cnds {
		if !cond {
			return false
		}
	}
	return true
}

// PasswordValidator is a custom validator that checks if a password is valid.
// A valid password must have at least 8 characters, contain at least one upper 
// case letter, one lower case letter, one digit, and one special character.
func PasswordValidator(f validator.FieldLevel) bool {
	const MIN_PASSWORD_LENGTH = 8

	pw := f.Field().String()

	if len(pw) < MIN_PASSWORD_LENGTH {
		return false
	}

	hasUpper, hasLower, hasDigit, hasSpecial := false, false, false, false
	for _, char := range pw {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}

	return all(hasUpper, hasLower, hasDigit, hasSpecial)
}
