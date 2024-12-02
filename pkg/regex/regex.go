package regex

import (
	"errors"
	"regexp"
)

func Password(pass string) error {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&#])[A-Za-z\d@$!%*?&#]{8,}$`

	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(pass) {
		return errors.New("Password must contain at least one uppercase letter, one lowercase letter, one number and one special character")
	}

	return nil
}
