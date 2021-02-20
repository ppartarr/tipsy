package checkers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFrequencyBlacklist = map[string]int{
	"password":        100,
	"hello":           90,
	"iloveyou":        80,
	"Iloveyou":        70,
	"there":           60,
	"foobar":          50,
	"duckduckgo":      40,
	"Password!":       30,
	"nomore":          20,
	"babydon'thurtme": 10,
}

func TestCheckOptimal(t *testing.T) {
	checker := NewChecker(testTypos, topCorrectors)

	ball := checker.CheckOptimal("password!", testFrequencyBlacklist, 5)
	if !assert.ElementsMatch(t, ball, []string{"Password!", "password"}) {
		t.Error("ball should be the set of strings containing the output of the correctors, unless it's in the blacklist")
	}

	// Add more tests for optimal checker
}

func TestCalculateTypoProbability(t *testing.T) {
	checker := NewChecker(testTypos, topCorrectors)

	prob := checker.CalculateTypoProbability("same")

	if prob != float64(testTypos["same"])/float64(96964) {
		t.Error("calculate typo probability is not working...")
	}

}

func TestCalculateProbabilityPasswordInBlacklist(t *testing.T) {
	prob := PasswordProbability("password", testFrequencyBlacklist)
	if prob != float64(100)/float64(550) {
		t.Error("calculate typo probability password in blacklist is not working...")
	}
}

func TestGenerateCombinations(t *testing.T) {
	combinations := generateCombinations([]string{"a", "b", "c"})
	comb := [][]string{{"a"}, {"b"}, {"c"}, {"a", "b"}, {"a", "c"}, {"b", "c"}, {"a", "b", "c"}}
	if !assert.ElementsMatch(t, combinations, comb) {
		t.Error("generate combinations is broken...")
	}
}

func TestMaxProbability(t *testing.T) {
	combinations := make([]CombinationProbability, 0)
	comb1 := CombinationProbability{
		Passwords:   []string{"hello", "world"},
		Probability: 0.8,
	}
	comb2 := CombinationProbability{
		Passwords:   []string{"world"},
		Probability: 0.5,
	}
	combinations = append(combinations, comb1, comb2)

	max := maxProbability(combinations)
	if !assert.ElementsMatch(t, max.Passwords, comb1.Passwords) || max.Probability != comb1.Probability {
		t.Error("max probability is broken")
	}
}

func TestFindProbabilityOfQthPassword(t *testing.T) {
	prob := FindProbabilityOfQthPassword(testFrequencyBlacklist, 3)
	fmt.Println(prob)
	if prob != float64(80)/float64(550) {
		t.Error("should return the corrector prob")
	}
}

func TestFindOptimalSubset(t *testing.T) {
	ballProbability := map[string]float64{
		"password": 0.8,
		"hello":    0.7,
		"world":    0.65,
		"foo":      0.1,
	}
	comb := FindOptimalSubset(ballProbability, 0.7)
	if !assert.ElementsMatch(t, comb.Passwords, []string{"hello"}) {
		t.Error("find optimal subset should return the max")
	}

	comb = FindOptimalSubset(ballProbability, 0.75)
	if !assert.ElementsMatch(t, comb.Passwords, []string{"foo", "world"}) {
		t.Error("find optimal subset should return the max <= cutoff")
	}
}
