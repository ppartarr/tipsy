package checkers

import (
	"io/ioutil"
	"log"
	"strings"
)

// CheckBlacklist uses a blacklist of high-probability passwords. It checks the password or any password in the ball only if it isn't in the blacklist
func CheckBlacklist(submittedPassword string, registeredPassword string, blacklist []string) bool {

	var ball []string = GetBall(submittedPassword)

	// check the submitted password first
	if submittedPassword == registeredPassword {
		return true
	}

	for _, value := range ball {
		// check password in the ball only if it isn't in the blacklist
		if !StringInSlice(value, blacklist) {
			return registeredPassword == value
		}
	}

	return false
}

// LoadBlacklist loads a file of high-probability password e.g. ./data/blacklistRockYou1000.txt
func LoadBlacklist(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(content), "\n")
}
