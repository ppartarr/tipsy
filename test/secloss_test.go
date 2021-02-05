package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/config"
	"github.com/ppartarr/tipsy/correctors"
)

func TestSecLossAlways(t *testing.T) {
	server := &config.Server{
		Checker: &config.Checker{
			Always: true,
		},

		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-all":       1698,
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

		Correctors:        []string{correctors.SwitchAll, correctors.RemoveLast, correctors.SwitchFirst},
		MinPasswordLength: 8,
	}

	checker := "always"
	q := 1000
	ballSize := 3
	attackerListFiles := []string{"../data/muslim-withcount.txt", "../data/rockyou-1m-withcount.txt", "../data/phpbb-withcount.txt"}
	// attackerListFile := "../data/muslim-withcount.txt"
	// defenderListFile := "../data/muslim-withcount.txt"

	for _, attackerListFile := range attackerListFiles {
		// convert results to json
		result := greedyMaxCoverageHeap(server, q, ballSize, attackerListFile, attackerListFile)
		bytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			t.Error(err.Error())
		}

		// save json to file
		filename := buildFilename(q, ballSize, server.MinPasswordLength, getDatasetFromFilename(attackerListFile))
		err = ioutil.WriteFile(filepath.Join(checker, filename), bytes, 0666)
		if err != nil {
			t.Error(err.Error())
		}
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
			"swc-all":       1698,
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

		Correctors:        []string{correctors.SwitchAll, correctors.RemoveLast, correctors.SwitchFirst},
		MinPasswordLength: 8,
	}

	checker := "blacklist"
	q := 1000
	ballSize := 3
	attackerListFiles := []string{"../data/muslim-withcount.txt", "../data/rockyou-1m-withcount.txt", "../data/phpbb-withcount.txt"}
	// attackerListFile := "../data/muslim-withcount.txt"
	// defenderListFile := "../data/muslim-withcount.txt"

	for _, attackerListFile := range attackerListFiles {
		// convert results to json
		result := greedyMaxCoverageHeap(server, q, ballSize, attackerListFile, attackerListFile)
		bytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			t.Error(err.Error())
		}

		// save json to file
		filename := buildFilename(q, ballSize, server.MinPasswordLength, getDatasetFromFilename(attackerListFile))
		err = ioutil.WriteFile(filepath.Join(checker, filename), bytes, 0666)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestSecLossOptimal(t *testing.T) {
	server := &config.Server{
		Checker: &config.Checker{
			Optimal: &config.OptimalChecker{
				File:                    "../data/rockyou-1k-withcount.txt",
				QthMostProbablePassword: 1000,
			},
		},

		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-all":       1698,
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

		Correctors:        []string{correctors.SwitchAll, correctors.RemoveLast, correctors.SwitchFirst},
		MinPasswordLength: 8,
	}

	checker := "optimal"
	q := 1000
	ballSize := 3
	attackerListFiles := []string{"../data/muslim-withcount.txt", "../data/rockyou-1m-withcount.txt", "../data/phpbb-withcount.txt"}
	// TODO comment out for estimating attakcer
	//defenderListFiles := []string{"../data/muslim-withcount.txt", "../data/rockyou-1m-withcount.txt", "../data/phpbb-withcount.txt"}

	for _, attackerListFile := range attackerListFiles {
		// convert results to json
		result := greedyMaxCoverageHeap(server, q, ballSize, attackerListFile, attackerListFile)
		bytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			t.Error(err.Error())
		}

		// save json to file
		filename := buildFilename(q, ballSize, server.MinPasswordLength, getDatasetFromFilename(attackerListFile))
		err = ioutil.WriteFile(filepath.Join(checker, filename), bytes, 0666)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestSecLoss(t *testing.T) {

	checker := "optimal"
	q := 10
	ballSize := 3
	minPasswordLength := 6
	defenderListFile := "../data/muslim-withcount.txt"

	filename := buildFilename(q, ballSize, minPasswordLength, getDatasetFromFilename(defenderListFile))

	filepath := filepath.Join(checker, filename)

	// open results
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("unable to unmarshal json")
	}

	result := &Result{}

	// get json
	err = json.Unmarshal(bytes, result)
	if err != nil {
		log.Fatal("unable to unmarshal json")
	}

	defenderList := checkers.LoadFrequencyBlacklist(result.DefenderListFile, minPasswordLength)
	blacklist := checkers.LoadBlacklist(server.Checker.Blacklist.File)

	// add points
	qs := []int{10}
	for _, rateLimit := range qs {
		guesses := result.GuessList[:rateLimit]
		naiveGuesses := result.NaiveGuessList[:rateLimit]

		// init checker
		checker := checkers.NewChecker(server.Typos, server.Correctors)

		guessListBall := make([]string, 0)
		for _, password := range guesses {

			union := unionBall(password, server, checker, defenderList, blacklist)
			guessListBall = append(guessListBall, union...)
			guessListBall = append(guessListBall, password)
		}
		fmt.Println(guessListBall)

		lambdaQGreedy := ballProbability(guessListBall, defenderList)
		lambdaQ := ballProbability(naiveGuesses, defenderList)
		secloss := (lambdaQGreedy - lambdaQ)
		fmt.Println(rateLimit)
		fmt.Println("lambda q greedy: ", lambdaQGreedy)
		fmt.Println("lambda q: ", lambdaQ)
		fmt.Println("secloss: ", secloss)
	}
}

func TestSecLossDataset(t *testing.T) {
	checkerz := []string{"always", "blacklist", "optimal"}
	q := 1000
	ballSize := 3
	attackerListFile := "../data/phpbb-withcount.txt"
	// TODO comment out for estimating attakcer
	//defenderListFiles := []string{"../data/muslim-withcount.txt", "../data/rockyou-1m-withcount.txt", "../data/phpbb-withcount.txt"}

	for _, checker := range checkerz {
		var server *config.Server
		switch checker {
		case "always":
			server = getAlwaysConfig()
		case "blacklist":
			server = getBlacklistConfig("../data/rockyou-1k.txt")
		case "optimal":
			server = getOptimalConfig("../data/rockyou-1k-withcount.txt", q)
		}
		fmt.Println("running secloss optimal for", attackerListFile)
		// convert results to json
		result := greedyMaxCoverageHeap(server, q, ballSize, attackerListFile, attackerListFile)
		bytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			t.Error(err.Error())
		}

		// save json to file
		filename := buildFilename(q, ballSize, server.MinPasswordLength, getDatasetFromFilename(attackerListFile))
		err = ioutil.WriteFile(filepath.Join(checker, filename), bytes, 0666)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestSecLossDatasets(t *testing.T) {
	checkerz := []string{"always", "blacklist", "optimal"}
	q := 10
	ballSize := 3
	attackerDefenderListFiles := map[string][]string{
		"../data/rockyou-1m-withcount.txt": {"../data/phpbb-withcount.txt", "../data/muslim-withcount.txt"},
		"../data/muslim-withcount.txt":     {"../data/phpbb-withcount.txt", "../data/rockyou-1m-withcount.txt"},
		"../data/phpbb-withcount.txt":      {"../data/muslim-withcount.txt", "../data/rockyou-1m-withcount.txt"},
	}

	for attackerFile, defenderList := range attackerDefenderListFiles {
		for _, defenderFile := range defenderList {

			for _, checker := range checkerz {
				var server *config.Server

				switch checker {
				case "always":
					server = getAlwaysConfig()
				case "blacklist":

					if defenderFile == "../data/rockyou-1m-withcount.txt" {
						server = getBlacklistConfig("../data/rockyou-1k.txt")
					} else if defenderFile == "../data/phpbb-withcount.txt" {
						server = getBlacklistConfig("../data/phpbb-1k.txt")
					} else if defenderFile == "../data/muslim-withcount.txt" {
						server = getBlacklistConfig("../data/muslim-1k.txt")
					}

				case "optimal":
					if defenderFile == "../data/rockyou-1m-withcount.txt" {
						server = getOptimalConfig("../data/rockyou-1k-withcount.txt", q)
					} else if defenderFile == "../data/phpbb-withcount.txt" {
						server = getOptimalConfig("../data/phpbb-1k-withcount.txt", q)
					} else if defenderFile == "../data/muslim-withcount.txt" {
						server = getOptimalConfig("../data/muslim-1k-withcount.txt", q)
					}
				}
				fmt.Println("running secloss optimal for attacker", attackerFile, "and defender", defenderFile)
				// convert results to json
				result := greedyMaxCoverageHeap(server, q, ballSize, attackerFile, defenderFile)
				bytes, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					t.Error(err.Error())
				}

				// save json to file
				filename := buildFilename(q, ballSize, server.MinPasswordLength, getDatasetFromFilename(defenderFile))
				err = ioutil.WriteFile(filepath.Join("estimating", checker, filename), bytes, 0666)
				if err != nil {
					t.Error(err.Error())
				}
			}
		}
	}
}

func buildFilename(q int, ballsize int, minPasswordLength int, dataset string) string {
	return strconv.Itoa(q) + "-" + strconv.Itoa(ballsize) + "-" + strconv.Itoa(minPasswordLength) + "-" + dataset + ".json"
}

func getDatasetFromFilename(filename string) string {
	slice := strings.Split(filename, "/")
	slice = strings.Split(slice[len(slice)-1], "-")
	return slice[0]
}

func getAlwaysConfig() *config.Server {
	return &config.Server{
		Checker: &config.Checker{
			Always: true,
		},
		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-all":       1698,
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
		Correctors:        []string{correctors.SwitchAll, correctors.RemoveLast, correctors.SwitchFirst},
		MinPasswordLength: 8,
	}
}

func getBlacklistConfig(blacklist string) *config.Server {
	return &config.Server{
		Checker: &config.Checker{
			Blacklist: &config.BlacklistChecker{
				File: blacklist,
			},
		},
		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-all":       1698,
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
		Correctors:        []string{correctors.SwitchAll, correctors.RemoveLast, correctors.SwitchFirst},
		MinPasswordLength: 8,
	}
}

func getOptimalConfig(optimal string, q int) *config.Server {
	return &config.Server{
		Checker: &config.Checker{
			Optimal: &config.OptimalChecker{
				File:                    optimal,
				QthMostProbablePassword: q,
			},
		},
		Typos: map[string]int{
			"same":          90234,
			"other":         1918,
			"swc-all":       1698,
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
		Correctors:        []string{correctors.SwitchAll, correctors.RemoveLast, correctors.SwitchFirst},
		MinPasswordLength: 8,
	}
}
