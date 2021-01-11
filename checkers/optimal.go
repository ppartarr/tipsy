package checkers

import (
	"log"
	"sort"

	"github.com/ppartarr/tipsy/correctors"
	"gonum.org/v1/gonum/stat/combin"
)

// CheckOptimal use the given distribution of passwords and a distribution of typos to decide whether to correct the typo or not
func CheckOptimal(password string, frequencyBlacklist map[string]int) bool {
	// TODO get the password from db
	registeredPassword := "password"

	var ball map[string]string = getBallWithCorrectionType(password)
	var ballProbability = make(map[string]float64)

	for passwordInBall, correctionType := range ball {
		// probability of guessing the password in the ball from the blacklist
		passwordProbability := calculateProbabilityPasswordInBlacklist(passwordInBall, frequencyBlacklist)

		// probability that the user made the user made the typo associated to the correction e.g. swc-all
		typoProbability := calculateTypoProbability(correctionType)

		// TODO change this to make probs customisable e.g. ngram vs pcfg vs historgram vs pwmodel
		ballProbability[passwordInBall] = passwordProbability * typoProbability
	}

	// find the optimal set of passwords in the ball such that aggregate probability of each password in the ball
	// is lower than the probability of the qth most probable password in the blacklist
	// we use q = 10
	probabilityOfQthPassword := findProbabilityOfQthPassword(frequencyBlacklist, 10)
	cutoff := float64(probabilityOfQthPassword) - calculateProbabilityPasswordInBlacklist(password, frequencyBlacklist)

	// get the set of password that maximises utility subject to completeness and security
	combinationToTry := combinationProbability{}
	combinationToTry = findOptimalSubset(ballProbability, cutoff, frequencyBlacklist)

	log.Println(combinationToTry)
	for _, password := range combinationToTry.passwords {
		// log.Println(password)
		if registeredPassword == password {
			return true
		}
	}

	return false
}

// given the blacklist, find the probability of the qth password in the distribution
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

// given the ball of a password, will solve a simple optimisation problem to find the
// set of password such that the total aggregate probability of the set is lower than that of the
// qth most probable password
// returns the set with the highest utility
func findOptimalSubset(ballProbability map[string]float64, cutoff float64, frequencyBlacklist map[string]int) combinationProbability {
	log.Println("ball probability: ", ballProbability)
	log.Println("cutoff: ", cutoff)
	passwordsInBall := make([]string, len(ballProbability))

	for word := range ballProbability {
		passwordsInBall = append(passwordsInBall, word)
	}

	passwordsInBall = deleteEmpty(passwordsInBall)

	combinations := generateCombinations(passwordsInBall)
	combinationsProbability := []combinationProbability{}

	// calculate the aggregate probability of each password in a set
	for _, combination := range combinations {
		combinationProbability := combinationProbability{}
		for _, password := range combination {
			combinationProbability.addPassword(password)
			combinationProbability.addProbability(ballProbability[password])
		}
		combinationsProbability = append(combinationsProbability, combinationProbability)
	}

	// build a new slice with combinations whose probability is smaller or equal to the cutoff
	filteredCombinations := []combinationProbability{}

	for _, combinationProbability := range combinationsProbability {
		if combinationProbability.probability <= cutoff {
			filteredCombinations = append(filteredCombinations, combinationProbability)
		}
	}

	return maxProbability(filteredCombinations)
}

// given a slice of filteredCombinations, will return the combination with the highest probability
func maxProbability(filteredCombinations []combinationProbability) combinationProbability {

	maxCombination := combinationProbability{}
	maxFloat := 0.0
	for _, filteredCombination := range filteredCombinations {
		if filteredCombination.probability != 0 && filteredCombination.probability > maxFloat {
			maxCombination = filteredCombination
		}
	}
	return maxCombination
}

type combinationProbability struct {
	passwords   []string
	probability float64
}

func (c *combinationProbability) addPassword(password string) string {
	c.passwords = append(c.passwords, password)
	return password
}

func (c *combinationProbability) addProbability(probability float64) float64 {
	c.probability += probability
	return probability
}

// given a slice of strings, will generate every combination of that slice
// e.g. given [a b c] will return [[a] [b] [c] [a b] [a c] [b c] [a b c]]
func generateCombinations(passwordsInBall []string) (combinations [][]string) {
	passwordsInBall = deleteEmpty(passwordsInBall)

	for i := 1; i <= len(passwordsInBall); i++ {
		intCombinations := combin.Combinations(len(passwordsInBall), i)
		for _, intCombination := range intCombinations {
			wordSlice := make([]string, i)
			for _, value := range intCombination {
				wordSlice = append(wordSlice, passwordsInBall[value])
				wordSlice = deleteEmpty(wordSlice)
			}
			combinations = append(combinations, wordSlice)
		}
	}
	return combinations
}

// remove all empty strings from a slice
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// calculate the probability that the password is picked from the frequencyBlacklist
func calculateProbabilityPasswordInBlacklist(password string, frequencyBlacklist map[string]int) float64 {
	totalNumberOfPassword := 0

	for _, value := range frequencyBlacklist {
		totalNumberOfPassword += value
	}

	return float64(frequencyBlacklist[password]) / float64(totalNumberOfPassword)
}

// calculate the probability that a correction is being used using typo frequencies from Chatterjee et al.
func calculateTypoProbability(correctionType string) float64 {
	// number of password fixed by the given corrector in Chatterjee et al.'s study
	typoFixProbability := make(map[string]float64)

	// init frequencies total = 96963
	typoFixFrequency := map[string]int{
		"same":          90234,
		"other":         1918,
		"swc-all":       1698,
		"kclose":        1385,
		"keypress-edit": 1000,
		"rm-last":       382,
		"swc-first":     209,
		"rm-firstc":     55,
		"sws-last1":     19,
		"tcerror":       18,
		"sws-lastn":     14,
		"upncap":        13,
		"n2s-last":      9,
		"cap2up":        5,
		"add1-last":     5,
	}

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
func getBallWithCorrectionType(password string) map[string]string {
	var ball = make(map[string]string)

	ball[correctors.SwitchCaseAll(password)] = "swc-all"
	ball[correctors.SwitchCaseFirstLetter(password)] = "swc-first"
	ball[correctors.RemoveLastChar(password)] = "rm-last"

	return ball
}
