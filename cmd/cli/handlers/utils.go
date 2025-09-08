package handlers

import "regexp"

func validateEmail(email string) bool {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if err != nil || !matched {
		return false
	}
	return true
}

func validatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}
