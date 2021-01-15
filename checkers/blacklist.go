package checkers

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/ppartarr/tipsy/correctors"
)

// CheckBlacklist uses a blacklist of high-probability passwords. It checks the password or any password in the ball only if it isn't in the blacklist
func (checker *CheckerService) CheckBlacklist(submittedPassword string, registeredPassword string, blacklist []string) bool {

	// check the submitted password first
	if CheckPasswordHash(submittedPassword, registeredPassword) {
		return true
	}

	// get n best correctors
	nBestCorrectors := correctors.GetNBestCorrectors(checker.NumberOfCorrectors, checker.TypoFrequency)

	// get the ball
	ball := correctors.GetBall(submittedPassword, nBestCorrectors)

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

// LoadBlacklist loads a file of high-probability password e.g. ./data/blacklistRockYou1000.txt
func LoadBlacklist(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(content), "\n")
}
