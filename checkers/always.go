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
	var ball []string = getBall(password)
	var unionBall []string = append(ball, password)

	return stringInSlice(registeredPassword, unionBall)
}

// this is the union of the ball and the submitted password
func getBall(password string) []string {
	var ball []string

	return append(ball,
		correctors.SwitchCaseAll(password),
		correctors.SwitchCaseFirstLetter(password),
		correctors.RemoveLastChar(password),
	)
}

func stringInSlice(s string, list []string) bool {
	for _, value := range list {
		if value == s {
			return true
		}
	}
	return false
}
