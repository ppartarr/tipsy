package test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ppartarr/tipsy/config"
)

func TestSecLossAlways(t *testing.T) {
	server := &config.Server{
		Checker: &config.Checker{
			Always: true,
		},

		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-al":        169,
			"kclose":        1385,
			"keypress-edit": 1000,
			"rm-last":       382,
			"swc-first":     209,
			"rm-first":      55,
			"sws-last":      19,
			"tcerror":       18,
			"sws-lastn":     14,
			"upncap":        13,
			"n2s-last":      9,
			"cap2up":        5,
			"add1-last":     5,
		},

		Correctors: []string{"swc-all", "rm-last", "swc-first"},
	}

	checker := "always"
	q := 10
	ballSize := 3
	maxPasswordLength := 6
	attackerListFile := "../data/rockyou-withcount.txt"
	defenderListFile := "../data/rockyou-withcount.txt"

	// convert results to json
	result := greedyMaxCoverageHeap(server, q, ballSize, maxPasswordLength, attackerListFile, defenderListFile)
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Error(err.Error())
	}

	// save json to file
	filename := strconv.Itoa(q) + "-" + strconv.Itoa(ballSize) + "-" + strconv.Itoa(maxPasswordLength) + ".json"
	err = ioutil.WriteFile(filepath.Join(checker, filename), bytes, 0666)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSecLossBlacklist(t *testing.T) {
	server := &config.Server{
		Checker: &config.Checker{
			Blacklist: &config.BlacklistChecker{
				File: "../data/rockyou-1k.txt",
			},
		},

		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-al":        169,
			"kclose":        1385,
			"keypress-edit": 1000,
			"rm-last":       382,
			"swc-first":     209,
			"rm-first":      55,
			"sws-last":      19,
			"tcerror":       18,
			"sws-lastn":     14,
			"upncap":        13,
			"n2s-last":      9,
			"cap2up":        5,
			"add1-last":     5,
		},

		Correctors: []string{"swc-all", "rm-last", "swc-first"},
	}

	checker := "blacklist"
	q := 10
	ballSize := 3
	maxPasswordLength := 6
	attackerListFile := "../data/rockyou-withcount.txt"
	defenderListFile := "../data/rockyou-withcount.txt"

	// convert results to json
	result := greedyMaxCoverageHeap(server, q, ballSize, maxPasswordLength, attackerListFile, defenderListFile)
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Error(err.Error())
	}

	// save json to file
	filename := strconv.Itoa(q) + "-" + strconv.Itoa(ballSize) + "-" + strconv.Itoa(maxPasswordLength) + ".json"
	err = ioutil.WriteFile(filepath.Join(checker, filename), bytes, 0666)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestSecLossOptimal(t *testing.T) {
	server := &config.Server{
		Checker: &config.Checker{
			Optimal: &config.OptimalChecker{
				File:                    "data/rockyou-1k.txt",
				QthMostProbablePassword: 10,
			},
		},

		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-al":        169,
			"kclose":        1385,
			"keypress-edit": 1000,
			"rm-last":       382,
			"swc-first":     209,
			"rm-first":      55,
			"sws-last":      19,
			"tcerror":       18,
			"sws-lastn":     14,
			"upncap":        13,
			"n2s-last":      9,
			"cap2up":        5,
			"add1-last":     5,
		},

		Correctors: []string{"swc-all", "rm-last", "swc-first"},
	}

	checker := "optimal"
	q := 10
	ballSize := 3
	maxPasswordLength := 6
	attackerListFile := "../data/rockyou-1m-withcount.txt"
	defenderListFile := "../data/rockyou-1m-withcount.txt"

	// convert results to json
	result := greedyMaxCoverageHeap(server, q, ballSize, maxPasswordLength, attackerListFile, defenderListFile)
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Error(err.Error())
	}

	// save json to file
	filename := strconv.Itoa(q) + "-" + strconv.Itoa(ballSize) + "-" + strconv.Itoa(maxPasswordLength) + ".json"
	err = ioutil.WriteFile(filepath.Join(checker, filename), bytes, 0666)
	if err != nil {
		t.Error(err.Error())
	}
}
