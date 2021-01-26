package checkers

import (
	"github.com/ppartarr/tipsy/correctors"
)

// CheckAlways checks the password and the passwords in the ball by using the given correctors
func (checker *Checker) CheckAlways(submittedPassword string) []string {

	// get the ball
	ball := correctors.GetBall(submittedPassword, checker.Correctors)

	return ball
}
