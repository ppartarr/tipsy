package checkers

import (
	"github.com/ppartarr/tipsy/correctors"
	"golang.org/x/crypto/bcrypt"
)

// CheckAlways checks the password and the passwords in the ball by using the given correctors
func (checker *Checker) CheckAlways(submittedPassword string, registeredPasswordHash string) bool {
	// check the submitted password first
	if CheckPasswordHash(submittedPassword, registeredPasswordHash) {
		return true
	}

	// get the ball
	ball := correctors.GetBall(submittedPassword, checker.Correctors)

	// constant-time check of the remainder of the ball
	success := false
	for _, value := range ball {
		if CheckPasswordHash(value, registeredPasswordHash) {
			success = true
		}
	}

	return success
}

// CheckPasswordHash verifies that the password arguments matches the given hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetBall is the set of strings obtained by applying correctors to the input password
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
