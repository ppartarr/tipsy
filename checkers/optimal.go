package checkers

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ppartarr/tipsy/correctors"
	"gonum.org/v1/gonum/stat/combin"
)

// CheckOptimal use the given distribution of passwords and a distribution of typos to decide whether to correct the typo or not
func (checker *Checker) CheckOptimal(submittedPassword string, frequencyBlacklist map[string]int, q int) []string {

	var ball map[string]string = correctors.GetBallWithCorrectionType(submittedPassword, checker.Correctors)
	var ballProbability = make(map[string]float64)

	for passwordInBall, correctionType := range ball {
		// probability of guessing the password in the ball from the blacklist
		PasswordProbability := PasswordProbability(passwordInBall, frequencyBlacklist)

		// probability that the user made the user made the typo associated to the correction e.g. swc-all
		typoProbability := checker.CalculateTypoProbability(correctionType)

		// TODO change this to make probs customisable e.g. ngram vs pcfg vs historgram vs pwmodel
		// only add password to ball if PasswordProbability * typoProbability > 0
		ballProbability[passwordInBall] = PasswordProbability * typoProbability
	}

	// find the optimal set of passwords in the ball such that aggregate probability of each password in the ball
	// is lower than the probability of the qth most probable password in the blacklist
	probabilityOfQthPassword := FindFrequencyOfQthPassword(frequencyBlacklist, q)
	cutoff := float64(probabilityOfQthPassword) - PasswordProbability(submittedPassword, frequencyBlacklist)

	// get the set of passwords that maximises utility subject to completeness and security
	combinationToTry := CombinationProbability{}
	combinationToTry = FindOptimalSubset(ballProbability, cutoff)

	return combinationToTry.Passwords
}

// FindFrequencyOfQthPassword given the blacklist, find the probability of the qth password in the distribution
func FindFrequencyOfQthPassword(frequencyBlacklist map[string]int, q int) int {
	sortedSlice := correctors.ConvertMapToSortedSlice(frequencyBlacklist)
	return sortedSlice[q].Value
}

// FindProbabilityOfQthPassword given the blacklist, find the probability of the qth password in the distribution
func FindProbabilityOfQthPassword(frequencyBlacklist map[string]int, q int) float64 {

	if q >= len(frequencyBlacklist) {
		return 0
	}

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

	return 0
}

// FindOptimalSubset given the ball of a password, will solve a simple optimisation problem to find the
// set of password such that the total aggregate probability of the set is lower than that of the
// qth most probable password
// returns the set with the highest utility
func FindOptimalSubset(ballProbability map[string]float64, cutoff float64) CombinationProbability {
	passwordsInBall := make([]string, len(ballProbability))
	i := 0
	for word := range ballProbability {
		passwordsInBall[i] = word
		i++
	}
	passwordsInBall = correctors.DeleteEmpty(passwordsInBall)

	combinations := generateCombinations(passwordsInBall)
	combinationsProbability := make([]CombinationProbability, len(combinations))

	// calculate the aggregate probability of each password in a set
	for index, combination := range combinations {
		combinationProbability := CombinationProbability{}
		for _, password := range combination {
			combinationProbability.addPassword(password)
			combinationProbability.addProbability(ballProbability[password])
		}
		combinationsProbability[index] = combinationProbability
	}

	// build a new slice with combinations whose probability is smaller or equal to the cutoff
	filteredCombinations := []CombinationProbability{}

	for _, combinationProbability := range combinationsProbability {
		if combinationProbability.Probability <= cutoff {
			filteredCombinations = append(filteredCombinations, combinationProbability)
		}
	}

	return maxProbability(filteredCombinations)
}

// MaxProbability given a slice of filteredCombinations, will return the combination with the highest probability
func maxProbability(filteredCombinations []CombinationProbability) CombinationProbability {

	maxCombination := CombinationProbability{}
	maxFloat := 0.0
	for _, filteredCombination := range filteredCombinations {
		if filteredCombination.Probability != 0.0 && filteredCombination.Probability > maxFloat {
			// check if probability is the same
			if filteredCombination.Probability == maxFloat {
				// use the combinations with fewer passwords
				if len(filteredCombination.Passwords) < len(maxCombination.Passwords) {
					maxCombination = filteredCombination
					maxFloat = filteredCombination.Probability
				}
			}
			maxCombination = filteredCombination
			maxFloat = filteredCombination.Probability
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

// PasswordProbability calculates the probability of a password in a list
func PasswordProbability(password string, frequencies map[string]int) float64 {
	probability, ok := frequencies[password]

	if ok {
		return float64(probability) / float64(totalNumberOfPasswords(frequencies))
	}
	return 0
}

func totalNumberOfPasswords(frequencies map[string]int) int {
	sum := 0
	for _, frequency := range frequencies {
		sum += frequency
	}
	return sum
}

func (c *CombinationProbability) addProbability(probability float64) float64 {
	c.Probability += probability
	return probability
}

// GenerateCombinations given a slice of strings, will generate every combination of that slice
// e.g. given [a b c] will return [[a] [b] [c] [a b] [a c] [b c] [a b c]]
func generateCombinations(passwordsInBall []string) [][]string {

	combinations := make([][]string, int(math.Pow(2, float64(len(passwordsInBall)))))

	for i := 1; i <= len(passwordsInBall); i++ {
		intCombinations := combin.Combinations(len(passwordsInBall), i)
		for _, intCombination := range intCombinations {
			wordSlice := make([]string, i)
			for index, value := range intCombination {
				wordSlice[index] = passwordsInBall[value]
				// wordSlice = correctors.DeleteEmpty(wordSlice)
			}
			combinations = append(combinations, wordSlice)
		}
	}
	return combinations
}

// CalculateTypoProbability calculate the probability that a correction is being used using typo frequencies from Chatterjee et al.
func (checker *Checker) CalculateTypoProbability(correctionType string) float64 {
	// number of password fixed by the given corrector in Chatterjee et al.'s study
	typoFixProbability := make(map[string]float64)

	totalNumberOfCorrections := 0

	// calculate the total number of corrections
	for _, frequency := range checker.TypoFrequency {
		totalNumberOfCorrections += frequency
	}

	// convert frequency into probability
	for correction, frequency := range checker.TypoFrequency {
		typoFixProbability[correction] = float64(frequency) / float64(totalNumberOfCorrections)
	}

	return typoFixProbability[correctionType]
}

// LoadFrequencyBlacklist loads a file of frequency + high-probability password e.g. ./data/rockyou-1k-withcount.txt
func LoadFrequencyBlacklist(filename string, minPasswordLength int) map[string]int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data map[string]int = make(map[string]int)

	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		// line := strings.Split(scanner.Text(), " ")
		// log.Println(line)
		frequency, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}

		// TODO add support for whitespace passwords
		if len(line) > 1 {
			word := line[1]
			if len(word) >= minPasswordLength {
				data[word] = frequency
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
