package checkers

import (
	"github.com/ppartarr/tipsy/correctors"
	"golang.org/x/crypto/bcrypt"
)

// CheckAlways checks the password and the passwords in the ball by using the given correctors
func (checker *CheckerService) CheckAlways(submittedPassword string, registeredPassword string) bool {
	// TODO make these run in constant time to avoid side-channels

	// check the submitted password first
	if CheckPasswordHash(submittedPassword, registeredPassword) {
		return true
	}

	// get n best correctors
	nBestCorrectors := correctors.GetNBestCorrectors(checker.NumberOfCorrectors, checker.TypoFrequency)

	// get the ball
	ball := correctors.GetBall(submittedPassword, nBestCorrectors)

	// constant-time check of the remainder of the ball
	success := false
	for _, value := range ball {
		if CheckPasswordHash(value, registeredPassword) {
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

// GetBallForN is the union of the ball and the submitted password
func GetBallForN(password string, numberOfCorrectors int) []string {
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
