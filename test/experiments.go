package test

import (
	"container/heap"
	"fmt"
	"log"
	"time"

	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/config"
	"github.com/ppartarr/tipsy/correctors"
	"github.com/thoas/go-funk"
)

// Ball alias
type Ball []string

var typoFrequencies map[string]int = map[string]int{
	"swc-all":   1698,
	"rm-last":   382,
	"swc-first": 209,
	"rm-first":  55,
	"sws-last1": 19,
	"sws-lastn": 14,
	"upncap":    13,
	"n2s-last":  9,
	"cap2up":    5,
	"add1-last": 5,
}

// Result contains all the vars necessary to reproduce a test
type Result struct {
	GuessList        []string
	NaiveGuessList   []string
	LambdaQ          float64
	LambdaQFuzzy     float64
	SecLoss          float64
	AttackerListFile string
	DefenderListFile string
	Correctors       []string
}

func greedyMaxCoverageHeap(config *config.Server, q int, ballSize int, minPasswordLength int, attackerListFile, defenderListFile string) *Result {

	// init checker
	checker := checkers.NewChecker(config.Typos, config.Correctors)

	// sample from password leaks
	now := time.Now()
	defenderList := checkers.LoadFrequencyBlacklist(defenderListFile, minPasswordLength)
	attackerList := checkers.LoadFrequencyBlacklist(attackerListFile, minPasswordLength)
	fmt.Println(time.Since(now))
	now = time.Now()

	var (
		guessList          []string
		naiveGuessList     []string
		priorityQueue          = make(PriorityQueue, 0)
		done                   = make(map[string]bool)
		length                 = 1
		sortedAttackerList     = correctors.ConvertMapToSortedSlice(attackerList)
		defenderListIndex  int = 0
		blacklist          []string
	)

	if config.Checker.Blacklist != nil {
		blacklist = checkers.LoadBlacklist(config.Checker.Blacklist.File)
	}

	for len(guessList) < q {
		// get the next most probable password from the attacker's password list
		if defenderListIndex < len(sortedAttackerList) {
			registeredPassword := sortedAttackerList[defenderListIndex].Key

			// check that it's longer than 6 chars
			if len(registeredPassword) < minPasswordLength {
				defenderListIndex++
				continue
			}

			// append the guess to the naive guess list
			if len(naiveGuessList) < q {
				naiveGuessList = append(naiveGuessList, registeredPassword)
			}

			if priorityQueue.Len() > 0 {
				item := heap.Pop(&priorityQueue).(*Item)
				// log.Println(item)
				item.weight = -item.weight
				if item.weight <= 0 {
					// log.Println(item)
					// log.Println(item.weight)
					log.Println("you have exhausted all the options")
					break
				}

				// printPriorityQueue(&priorityQueue, item.value, "rockyouG")

				for float64(item.weight) > float64(ballSize)*checkers.PasswordProbability(registeredPassword, attackerList) && len(guessList) < q {
					// add guess to guess list
					log.Println("Guess", len(guessList), "/", q, "password:", item.value, "weight:", float64(item.weight))
					guessList = append(guessList, item.value)

					// add password & ball to done
					killed := unionBallNotDone(item.value, done, config, checker, attackerList, q, blacklist)
					for _, password := range killed {
						done[password] = true
					}

					// log.Println(item.value)
					// log.Println(killed)
					// log.Println(done)

					// add all neighbours of this password to the priority queue
					for _, password := range killed {
						probability := checkers.PasswordProbability(registeredPassword, attackerList)
						neighbours := getNeighbours(password, config.Correctors, config, checker, attackerList, q, blacklist)
						neighbours = append(neighbours, password)
						for _, neighbour := range neighbours {
							// update neighbour weight in the priority queue
							neighbourItem := priorityQueue.Find(neighbour)
							if neighbourItem != nil {
								priorityQueue.update(neighbourItem, neighbourItem.value, neighbourItem.weight-probability)

								// remove neighbourItem from priorityQueue if it's weight is > 0
								if neighbourItem.weight <= 0 {
									// log.Println("removing from priority q:", neighbourItem.value)
									heap.Remove(&priorityQueue, neighbourItem.index)
								}
							}
						}
					}

					// pop new item off of priority queue
					if priorityQueue.Len() > 0 {
						item = heap.Pop(&priorityQueue).(*Item)
						item.weight = -item.weight
					} else {
						// log.Println("cannot pop item off of priority q")
						break
					}
				}

				// add item to priority q
				if item.weight > 0 && priorityQueue.Find(item.value) == nil && !correctors.StringInSlice(item.value, guessList) {
					// log.Println("push item", item)
					item.weight = -item.weight
					heap.Push(&priorityQueue, item)
				}
			}

			// insert neighbours & password into the priority queue
			neighbours := getNeighbours(registeredPassword, config.Correctors, config, checker, attackerList, q, blacklist)
			neighbours = append(neighbours, registeredPassword)

			allNeighbours := neighbours

			for _, neighbour := range neighbours {
				// don't add neighbour if it's already in the priority queue
				if priorityQueue.Find(neighbour) != nil {
					allNeighbours = remove(allNeighbours, neighbour)
				}
				// don't add neighbour to priority queue if it's already been tested
				_, ok := done[neighbour]
				if ok {
					allNeighbours = remove(allNeighbours, neighbour)
				}
			}

			// add items to the priority queue
			for _, neighbour := range allNeighbours {
				weight := power(neighbour, attackerList, done, config, checker, q, blacklist)
				item := &Item{
					value:  neighbour,
					weight: -weight,
				}
				heap.Push(&priorityQueue, item)
			}

			// print heap size update
			if priorityQueue.Len() > length {
				log.Println("Heap size:", priorityQueue.Len())
				length = priorityQueue.Len() * 2
			}

			defenderListIndex++
		} else {
			fmt.Println("out of options")
		}
	}

	// create guess list ball for lambda q calculation
	var guessListBall []string
	for _, password := range guessList {
		// get ball of passwords in guess list
		// ball := unionBall(password, done, config, checker, attackerList, q, blacklist)
		ball := correctors.GetBall(password, config.Correctors)
		ball = append(ball, password)
		guessListBall = append(guessListBall, ball...)
	}

	guessListBall = correctors.DeleteEmpty(guessListBall)
	log.Println(guessListBall)

	// remove duplicates from guess list
	guessListBall = funk.UniqString(guessListBall)

	lambdaQGreedy := ballProbability(guessListBall, defenderList)
	lambdaQ := ballProbability(naiveGuessList, defenderList)
	log.Println("typo guess list:", guessList)
	log.Println("normal guess list:", naiveGuessList)
	log.Println("lambda q", lambdaQ)
	log.Println("lambda q greedy", lambdaQGreedy)
	log.Println("sec loss", lambdaQGreedy-lambdaQ)
	result := &Result{
		GuessList:        guessList,
		NaiveGuessList:   naiveGuessList,
		LambdaQ:          lambdaQ,
		LambdaQFuzzy:     lambdaQGreedy,
		SecLoss:          lambdaQGreedy - lambdaQ,
		AttackerListFile: attackerListFile,
		DefenderListFile: defenderListFile,
		Correctors:       config.Correctors,
	}
	fmt.Println(time.Since(now))
	now = time.Now()
	return result
}

func printPriorityQueue(pq *PriorityQueue, password string, match string) {
	if password == match {
		// print priority queue
		f := 0
		for pq.Len() > 0 && f < 100 {
			item := heap.Pop(pq).(*Item)
			fmt.Println(item)
			f++
		}
		log.Fatal()
	}
}

func ballProbability(ball Ball, frequencies map[string]int) float64 {
	ballProbability := 0.0
	for _, password := range ball {
		ballProbability += checkers.PasswordProbability(password, frequencies)
	}
	return ballProbability
}

// returns the union ball of passwords
func unionBall(password string, done map[string]bool, config *config.Server, checker *checkers.Checker, attackerList map[string]int, q int, blacklist []string) []string {
	// TODO make get ball configurable according to the checker
	unionBall := make([]string, 0)

	// TODO add support for typtop
	if config.Checker.Always {
		unionBall = checker.CheckAlways(password)
		// log.Println("always")
	} else if config.Checker.Blacklist != nil {
		// log.Println("blacklist")
		unionBall = checker.CheckBlacklist(password, blacklist)
	} else if config.Checker.Optimal != nil {
		// log.Println("optimal")
		unionBall = checker.CheckOptimal(password, attackerList, q)
	}

	// check if passwords are in done
	return unionBall
}

// returns the union ball of passwords that are not in the done list
func unionBallNotDone(password string, done map[string]bool, config *config.Server, checker *checkers.Checker, attackerList map[string]int, q int, blacklist []string) []string {
	// TODO make get ball configurable according to the checker
	unionBallNotDone := make([]string, 0)
	temp := unionBall(password, done, config, checker, attackerList, q, blacklist)

	// check if passwords are in done
	for _, str := range temp {
		if done[str] != true {
			unionBallNotDone = append(unionBallNotDone, str)
		}
	}
	if done[password] != true {
		unionBallNotDone = append(unionBallNotDone, password)
	}

	return unionBallNotDone
}

func power(password string, attackerList map[string]int, done map[string]bool, config *config.Server, checker *checkers.Checker, q int, blacklist []string) float64 {
	probability := 0.0

	// add passwords in ball to done
	unionBall := unionBallNotDone(password, done, config, checker, attackerList, q, blacklist)

	for _, pw := range unionBall {
		probability += checkers.PasswordProbability(pw, attackerList)
	}
	return probability
}

func remove(slice []string, s string) []string {
	for i, str := range slice {
		if str == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return nil
}

func getNeighbours(password string, bestCorrectors []string, conf *config.Server, checker *checkers.Checker, attackerList map[string]int, q int, blacklist []string) []string {
	neighbours := make([]string, 0)
	for _, corrector := range bestCorrectors {
		edits := make([]string, 0)

		if conf.Checker.Always {
			edits = correctors.ApplyInverseCorrectionFunction(corrector, password)
			neighbours = append(neighbours, edits...)
		} else if conf.Checker.Blacklist != nil {
			edits = correctors.ApplyInverseCorrectionFunction(corrector, password)
			for _, neighbour := range neighbours {
				blacklist := checkers.LoadBlacklist(conf.Checker.Blacklist.File)
				if !correctors.StringInSlice(neighbour, blacklist) {
					neighbours = append(neighbours, edits...)
				}
			}
		} else if conf.Checker.Optimal != nil {
			edits = CheckInverseOptimal(password, attackerList, q, checker)
		}

		neighbours = append(neighbours, edits...)
	}
	neighbours = correctors.DeleteEmpty(neighbours)
	return neighbours
}

func applyEdits(password string) []string {
	edits := make([]string, 10)

	for _, letter := range correctors.LetterRunes {
		for i := 0; i < len(password); i++ {
			// add every rune in every index
			edits = append(edits, password[:i]+string(letter)+password[i:])
		}
		edits = append(edits, password+string(letter))
	}
	for i := 0; i < len(password); i++ {
		edits = append(edits, password[:i]+password[i+1:])
	}
	return edits
}

func convertMapToSlice(in map[string]int) []string {
	var ss []string

	for key := range in {
		ss = append(ss, key)
	}

	return ss
}

// CheckInverseOptimal use the given distribution of passwords and a distribution of typos to decide whether to correct the typo or not
func CheckInverseOptimal(submittedPassword string, frequencyBlacklist map[string]int, q int, checker *checkers.Checker) []string {

	var ball map[string]string = getInverseBallWithCorrectionType(submittedPassword, checker.Correctors)
	var ballProbability = make(map[string]float64)

	for passwordInBall, correctionType := range ball {
		// probability of guessing the password in the ball from the blacklist
		PasswordProbability := checkers.PasswordProbability(passwordInBall, frequencyBlacklist)

		// probability that the user made the user made the typo associated to the correction e.g. swc-all
		typoProbability := checker.CalculateTypoProbability(correctionType)

		// TODO change this to make probs customisable e.g. ngram vs pcfg vs historgram vs pwmodel
		// only add password to ball if PasswordProbability * typoProbability > 0
		if PasswordProbability*typoProbability > 0 {
			ballProbability[passwordInBall] = PasswordProbability * typoProbability
		}
	}

	// find the optimal set of passwords in the ball such that aggregate probability of each password in the ball
	// is lower than the probability of the qth most probable password in the blacklist
	probabilityOfQthPassword := checkers.FindProbabilityOfQthPassword(frequencyBlacklist, q)
	cutoff := float64(probabilityOfQthPassword) - checkers.PasswordProbability(submittedPassword, frequencyBlacklist)

	// get the set of passwords that maximises utility subject to completeness and security
	combinationToTry := checkers.CombinationProbability{}
	combinationToTry = checkers.FindOptimalSubset(ballProbability, cutoff)

	return combinationToTry.Passwords
}

func getInverseBallWithCorrectionType(password string, corrections []string) map[string]string {
	var ballWithCorrectorName = make(map[string]string)

	for _, corrector := range corrections {
		correctedPasswords := correctors.ApplyInverseCorrectionFunction(corrector, password)
		for _, correctedPassword := range correctedPasswords {
			if correctedPassword != password {
				ballWithCorrectorName[correctedPassword] = corrector
			}
		}
	}

	return ballWithCorrectorName
}
