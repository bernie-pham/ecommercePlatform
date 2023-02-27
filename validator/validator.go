package val

import (
	"fmt"
	"regexp"
)

var (
	isEmailFormat = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9]+@[a-zA-Z]+.[a-zA-Z]+$`).MatchString
)

func ValidateEmail(value string) error {
	result := isEmailFormat(value)
	if !result {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// func IsExist(value string) error {
// 	if
// }
