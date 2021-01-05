package main

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/stat/combin"
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

// CheckBlacklist uses a blacklist of high-probability passwords. It checks the password or any password in the ball
func CheckBlacklist(password string, blacklist []string) bool {
	// TODO get the password from db
	registeredPassword := "password"

	var ball []string = getBall(password)

	// check the submitted password first
	if password == registeredPassword {
		return true
	} else {
		for _, value := range ball {
			// check password in the ball only if it isn't in the blacklist
			if !stringInSlice(value, blacklist) {
				return registeredPassword == value
			}
		}
	}

	return false
}

// CheckOptimal use the given distribution of passwords and a distribution of typos to decide whether to correct the typo or not
func CheckOptimal(password string, frequencyBlacklist map[string]int) bool {
	// TODO get the password from db
	// registeredPassword := "password"

	var ball map[string]string = getBallWithCorrectionType(password)
	var ballProbability = make(map[string]float64)

	for passwordInBall, correctionType := range ball {
		// get probability of the password in the blacklist
		// fmt.Println(passwordInBall)
		// fmt.Println(correctionType)

		// probability of guessing the password in the ball from the blacklist
		passwordProbability := calculatePasswordProbability(passwordInBall, frequencyBlacklist)
		typoProbability := calculateTypoProbability(correctionType)
		// fmt.Println(calculatePasswordProbability(passwordInBall, frequencyBlacklist))
		// fmt.Println(calculateTypoProbability(correctionType))
		// fmt.Println(passwordProbability * typoProbability)

		ballProbability[passwordInBall] = passwordProbability * typoProbability
	}

	// find the optimal set of passwords in the ball such that aggregate probability of each password in the ball
	// is lower than the probability of the qth most probable password in the blacklist
	// we use q = 10
	probabilityOfQthPassword := findProbabilityOfQthPassword(frequencyBlacklist, 10)
	cutoff := float64(probabilityOfQthPassword) - calculatePasswordProbability(password, frequencyBlacklist)

	findOptimalSubset(ballProbability, cutoff)

	return false
}

func findProbabilityOfQthPassword(frequencyBlacklist map[string]int, q int) float64 {
	frequencies := make([]int, len(frequencyBlacklist))
	totalNumberOfPassword := 0

	for _, frequency := range frequencyBlacklist {
		frequencies = append(frequencies, frequency)
		totalNumberOfPassword += frequency
	}

	sort.Ints(frequencies)

	for _, frequency := range frequencyBlacklist {
		if frequency == frequencies[len(frequencies)-q] {
			return float64(frequency) / float64(totalNumberOfPassword)
		}
	}

	return -1
}

func findOptimalSubset(ballProbability map[string]float64, cutoff float64) {
	fmt.Println(ballProbability)
	fmt.Println(cutoff)
	passwordsInBall := make([]string, len(ballProbability))

	for word := range ballProbability {
		passwordsInBall = append(passwordsInBall, word)
	}

	passwordsInBall = deleteEmpty(passwordsInBall)
	generateCombinations(passwordsInBall)
	// continue implementing get_most_val_under_prob https://github.com/rchatterjee/mistypography/blob/master/typofixer/common.py
}

// given a slice of strings, will generate every combination of that slice
func generateCombinations(passwordsInBall []string) (subsets [][]string) {
	passwordsInBall = deleteEmpty(passwordsInBall)

	for i := 1; i <= len(passwordsInBall); i++ {
		combinations := combin.Combinations(len(passwordsInBall), i)
		for _, combination := range combinations {
			wordSlice := make([]string, i)
			fmt.Println(combination)
			for _, value := range combination {
				wordSlice = append(wordSlice, passwordsInBall[value])
				wordSlice = deleteEmpty(wordSlice)
			}
			subsets = append(subsets, wordSlice)
		}
	}
	return subsets
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func calculatePasswordProbability(password string, frequencyBlacklist map[string]int) float64 {
	totalNumberOfPassword := 0

	for _, value := range frequencyBlacklist {
		totalNumberOfPassword += value
	}

	return float64(frequencyBlacklist[password]) / float64(totalNumberOfPassword)
}

// calculate the probability that a correction is being used using typo frequencies from Chatterjee et al.
func calculateTypoProbability(correctionType string) float64 {
	// number of password fixed by the given corrector in Chatterjee et al.'s study
	typoFixFrequency := make(map[string]int)
	typoFixProbability := make(map[string]float64)

	// init frequencies total = 96963
	typoFixFrequency["same"] = 90234
	typoFixFrequency["other"] = 1918
	typoFixFrequency["swc-all"] = 1698
	typoFixFrequency["kclose"] = 1385
	typoFixFrequency["keypress-edit"] = 1000
	// combined all rm-last
	typoFixFrequency["rm-last"] = 382
	typoFixFrequency["swc-first"] = 209
	typoFixFrequency["rm-firstc"] = 55
	typoFixFrequency["sws-last1"] = 19
	typoFixFrequency["tcerror"] = 18
	typoFixFrequency["sws-lastn"] = 14
	typoFixFrequency["upncap"] = 13
	typoFixFrequency["n2s-last"] = 9
	typoFixFrequency["cap2up"] = 5
	typoFixFrequency["add1-last"] = 5

	totalNumberOfCorrections := 0

	// calculate the total number of corrections
	for _, frequency := range typoFixFrequency {
		totalNumberOfCorrections += frequency
	}

	// convert frequency into probability
	for correction, frequency := range typoFixFrequency {
		typoFixProbability[correction] = float64(frequency) / float64(totalNumberOfCorrections)
	}

	return typoFixProbability[correctionType]
}

// this is the union of the ball and the submitted password
func getBall(password string) []string {
	var ball []string

	return append(ball,
		SwitchCaseAll(password),
		SwitchCaseFirstLetter(password),
		RemoveLastChar(password),
	)
}

func getBallWithCorrectionType(password string) map[string]string {
	var ball = make(map[string]string)

	ball[SwitchCaseAll(password)] = "swc-all"
	ball[SwitchCaseFirstLetter(password)] = "swc-first"
	ball[RemoveLastChar(password)] = "rm-last"

	return ball
}

func stringInSlice(s string, list []string) bool {
	for _, value := range list {
		if value == s {
			return true
		}
	}
	return false
}
