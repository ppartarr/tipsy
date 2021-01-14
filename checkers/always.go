package checkers

import (
	"github.com/ppartarr/tipsy/correctors"
)

// CheckAlways checks the password and the passwords in the ball by using the given correctors
func CheckAlways(password string, numberOfCheckers int) bool {
	// TODO make these run in constant time to avoid side-channels

	// TODO get the password from db
	registeredPassword := "password"

	// perform the check
	var ball []string = GetBall(password)
	var unionBall []string = append(ball, password)

	return StringInSlice(registeredPassword, unionBall)
}

// GetBall is the union of the ball and the submitted password
func GetBall(password string) []string {
	var ball []string

	return append(ball,
		correctors.SwitchCaseAll(password),
		correctors.SwitchCaseFirstLetter(password),
		correctors.RemoveLastChar(password),
	)
}

// StringInSlice returns true if s in is list
func StringInSlice(s string, list []string) bool {
	success := false
	for _, value := range list {
		if value == s {
			success = true
		}
	}
	return success
}
