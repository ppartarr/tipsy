package checkers

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ppartarr/tipsy/correctors"
	"gonum.org/v1/gonum/stat/combin"
)

// CheckOptimal use the given distribution of passwords and a distribution of typos to decide whether to correct the typo or not
func (checker *Checker) CheckOptimal(submittedPassword string, registeredPassword string, frequencyBlacklist map[string]int, q int) bool {

	// check the submitted password first
	if CheckPasswordHash(submittedPassword, registeredPassword) {
		return true
	}

	var ball map[string]string = correctors.GetBallWithCorrectionType(submittedPassword, checker.Correctors)
	var ballProbability = make(map[string]float64)

	for passwordInBall, correctionType := range ball {
		// probability of guessing the password in the ball from the blacklist
		passwordProbability := CalculateProbabilityPasswordInBlacklist(passwordInBall, frequencyBlacklist)

		// probability that the user made the user made the typo associated to the correction e.g. swc-all
		typoProbability := CalculateTypoProbability(correctionType)

		// TODO change this to make probs customisable e.g. ngram vs pcfg vs historgram vs pwmodel
		// only add password to ball if passwordProbability * typoProbability > 0
		if passwordProbability*typoProbability > 0 {
			ballProbability[passwordInBall] = passwordProbability * typoProbability
		}
	}

	// find the optimal set of passwords in the ball such that aggregate probability of each password in the ball
	// is lower than the probability of the qth most probable password in the blacklist
	probabilityOfQthPassword := FindProbabilityOfQthPassword(frequencyBlacklist, q)
	cutoff := float64(probabilityOfQthPassword) - CalculateProbabilityPasswordInBlacklist(submittedPassword, frequencyBlacklist)

	// get the set of password that maximises utility subject to completeness and security
	combinationToTry := CombinationProbability{}
	combinationToTry = FindOptimalSubset(ballProbability, cutoff, frequencyBlacklist)

	log.Println(combinationToTry)
	log.Println(probabilityOfQthPassword)
	log.Println(cutoff)

	// constant-time check of the remainder of the ball
	success := false
	for _, passwordsInCombination := range combinationToTry.Passwords {
		// log.Println(password)
		if CheckPasswordHash(passwordsInCombination, registeredPassword) {
			success = true
		}
	}

	return success
}

// FindProbabilityOfQthPassword given the blacklist, find the probability of the qth password in the distribution
func FindProbabilityOfQthPassword(frequencyBlacklist map[string]int, q int) float64 {
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

// FindOptimalSubset given the ball of a password, will solve a simple optimisation problem to find the
// set of password such that the total aggregate probability of the set is lower than that of the
// qth most probable password
// returns the set with the highest utility
func FindOptimalSubset(ballProbability map[string]float64, cutoff float64, frequencyBlacklist map[string]int) CombinationProbability {
	log.Println("ball probability: ", ballProbability)
	log.Println("cutoff: ", cutoff)

	passwordsInBall := make([]string, len(ballProbability))
	for word := range ballProbability {
		passwordsInBall = append(passwordsInBall, word)
	}
	passwordsInBall = DeleteEmpty(passwordsInBall)

	combinations := GenerateCombinations(passwordsInBall)
	combinationsProbability := []CombinationProbability{}

	// calculate the aggregate probability of each password in a set
	log.Println(combinations)
	for _, combination := range combinations {
		CombinationProbability := CombinationProbability{}
		for _, password := range combination {
			log.Println(password)
			log.Println(ballProbability[password])
			CombinationProbability.addPassword(password)
			CombinationProbability.addProbability(ballProbability[password])
		}
		combinationsProbability = append(combinationsProbability, CombinationProbability)
	}

	// build a new slice with combinations whose probability is smaller or equal to the cutoff
	filteredCombinations := []CombinationProbability{}

	for _, CombinationProbability := range combinationsProbability {
		if CombinationProbability.Probability <= cutoff {
			filteredCombinations = append(filteredCombinations, CombinationProbability)
		}
	}

	return MaxProbability(filteredCombinations)
}

// MaxProbability given a slice of filteredCombinations, will return the combination with the highest probability
func MaxProbability(filteredCombinations []CombinationProbability) CombinationProbability {

	maxCombination := CombinationProbability{}
	maxFloat := 0.0
	for _, filteredCombination := range filteredCombinations {
		if filteredCombination.Probability != 0 && filteredCombination.Probability > maxFloat {
			maxCombination = filteredCombination
		}
	}
	return maxCombination
}

// CombinationProbability is the aggregate probability of a set of passwords
type CombinationProbability struct {
	Passwords   []string
	Probability float64
}

func (c *CombinationProbability) addPassword(password string) string {
	c.Passwords = append(c.Passwords, password)
	return password
}

func (c *CombinationProbability) addProbability(probability float64) float64 {
	c.Probability += probability
	return probability
}

// GenerateCombinations given a slice of strings, will generate every combination of that slice
// e.g. given [a b c] will return [[a] [b] [c] [a b] [a c] [b c] [a b c]]
func GenerateCombinations(passwordsInBall []string) (combinations [][]string) {
	passwordsInBall = DeleteEmpty(passwordsInBall)

	for i := 1; i <= len(passwordsInBall); i++ {
		intCombinations := combin.Combinations(len(passwordsInBall), i)
		for _, intCombination := range intCombinations {
			wordSlice := make([]string, i)
			for index, value := range intCombination {
				wordSlice[index] = passwordsInBall[value]
				wordSlice = DeleteEmpty(wordSlice)
			}
			combinations = append(combinations, wordSlice)
		}
	}
	return combinations
}

// DeleteEmpty remove all empty strings from a slice
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// CalculateProbabilityPasswordInBlacklist calculate the probability that the password is picked from the frequencyBlacklist
func CalculateProbabilityPasswordInBlacklist(password string, frequencyBlacklist map[string]int) float64 {
	totalNumberOfPassword := 0

	for _, value := range frequencyBlacklist {
		totalNumberOfPassword += value
	}

	return float64(frequencyBlacklist[password]) / float64(totalNumberOfPassword)
}

// CalculateTypoProbability calculate the probability that a correction is being used using typo frequencies from Chatterjee et al.
func CalculateTypoProbability(correctionType string) float64 {
	// number of password fixed by the given corrector in Chatterjee et al.'s study
	typoFixProbability := make(map[string]float64)

	// init frequencies total = 96963
	// for every corrector, we give the frequency of typos that were corrected by it in Chatterjee et al's study
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

// LoadFrequencyBlacklist loads a file of frequency + high-probability password e.g. ./data/rockyou-withcount1000.txt
func LoadFrequencyBlacklist(filename string) map[string]int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data map[string]int = make(map[string]int)

	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		// log.Println(line)
		frequency, err := strconv.Atoi(line[0])
		word := line[1]

		if err != nil {
			log.Fatal(err)
		}
		data[word] = frequency
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
