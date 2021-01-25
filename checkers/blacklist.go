package checkers

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/ppartarr/tipsy/correctors"
)

// CheckBlacklist uses a blacklist of high-probability passwords. It checks the password or any password in the ball only if it isn't in the blacklist
func (checker *Checker) CheckBlacklist(submittedPassword string, registeredPassword string, blacklist []string) bool {

	// check the submitted password first
	if CheckPasswordHash(submittedPassword, registeredPassword) {
		return true
	}

	// get the ball
	ball := correctors.GetBall(submittedPassword, checker.Correctors)

	// constant-time check of the remainder of the ball
	succcess := false
	for _, value := range ball {
		log.Println(value)
		// check password in the ball only if it isn't in the blacklist
		if !StringInSlice(value, blacklist) {
			if CheckPasswordHash(value, registeredPassword) {
				succcess = true
			}
		}
	}

	return succcess
}

// GetBlacklistBall returns the ball to check
func (checker *Checker) GetBlacklistBall(submittedPassword string, blacklist []string) []string {

	// get the ball
	ball := correctors.GetBall(submittedPassword, checker.Correctors)

	for _, password := range ball {
		// check password in the ball only if it isn't in the blacklist
		if !StringInSlice(password, blacklist) {
			ball = remove(ball, password)
		}
	}

	return ball
}

func remove(slice []string, s string) []string {
	for i := range slice {
		if slice[i] == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// LoadBlacklist loads a file of high-probability password e.g. ./data/rockyou1000.txt
func LoadBlacklist(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(content), "\n")
}
