package checkers

import (
	"github.com/ppartarr/tipsy/correctors"
)

// CheckAlways checks the password and the passwords in the ball by using the given correctors
func (checker *Checker) CheckAlways(submittedPassword string, registeredPasswordHash string) []string {

	// get the ball
	ball := correctors.GetBall(submittedPassword, checker.Correctors)

	return ball
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
