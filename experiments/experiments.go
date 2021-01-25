package main

import (
	"container/heap"
	"fmt"
	"log"
	mrand "math/rand"

	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/correctors"
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

var topCorrectors = []string{"swc-all", "swc-first", "rm-last"}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`1234567890-=[]\\;',./~!@#$%^&*()_+{}|:\"<>?")

func testMaxHeap(attackerList map[string]int, blacklist []string, done map[string]bool) {
	// Some items and their priorities.
	items := map[string]int{
		"password": 59462, "iloveyou": 49952, "princess": 33291,
	}

	// Create a weight queue, put the items in it, and
	// establish the weight queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value := range items {
		pq[i] = &Item{
			value:  value,
			weight: power(value, attackerList, blacklist, done),
			index:  i,
		}
		log.Println(pq[i])
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its weight.
	item := &Item{
		value:  "rockyou",
		weight: power("rockyou", attackerList, blacklist, done),
	}
	pq.Push(item)

	// Take the items out; they arrive in decreasing weight order.
	for pq.Len() > 0 {
		item := pq.Pop()
		fmt.Println(item)
		// fmt.Printf("%.2f:%s ", item.weight, item.value)
	}

}

func main() {
	// seed randomness
	// mrand.Seed(time.Now().UnixNano())

	// sample from password leaks
	blacklist := checkers.LoadBlacklist("../data/rockyou1000.txt")
	fblacklist := checkers.LoadFrequencyBlacklist("../data/rockyou-withcount1000.txt", 6)
	attackerList := checkers.LoadFrequencyBlacklist("../data/rockyou-withcount1000.txt", 6)

	// TODO get from args
	q := 10

	var (
		priorityQueue  PriorityQueue = make(PriorityQueue, 0)
		guessList      []string
		naiveGuessList []string
		done           map[string]bool = make(map[string]bool)
		length         int             = 1
		// startTime            time.Time      = time.Now().UTC()
		ballSize           float64 = 5
		sortedAttackerList         = correctors.ConvertMapToSortedSlice(attackerList)
		fblacklistIndex    int     = 0
		minPasswordLength  int     = 6
	)

	log.Println("starting loop")

	for len(guessList) < q {
		// get the next most probable password from the attacker's password list
		if fblacklistIndex < len(sortedAttackerList) {
			registeredPassword := sortedAttackerList[fblacklistIndex].Key

			// check that it's longer than 6 chars
			if len(registeredPassword) < minPasswordLength {
				fblacklistIndex++
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
					log.Println(item)
					log.Println(item.weight)
					log.Println("you have exhausted all the options")
					break
				}
				for float64(item.weight) > ballSize*passwordProbability(registeredPassword, attackerList) && len(guessList) < q {
					// add guess to guess list
					log.Println("Guess", len(guessList), "/", q, "password:", item.value, "weight:", float64(item.weight)/float64(totalFrequencies(fblacklist)))
					guessList = append(guessList, item.value)

					// add password & ball to done
					killed := unionBallNotDone(item.value, done)
					for _, password := range killed {
						done[password] = true
					}

					// add all neighbours of this password to the priority queue
					for _, password := range killed {
						probability := passwordProbability(registeredPassword, attackerList)
						neighbours := getNeighbours(password, topCorrectors)
						neighbours = append(neighbours, password)
						for _, neighbour := range neighbours {
							// update neighbour weight in the priority queue
							neighbourItem := priorityQueue.Find(neighbour)
							if neighbourItem != nil {
								priorityQueue.update(neighbourItem, neighbourItem.value, neighbourItem.weight-probability)

								// remove neighbourItem from priorityQueue if it's weight is > 0
								// log.Println("removing from priority q")
								if neighbourItem.weight <= 0 {
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
						break
					}
				}

				// add item to priority q
				if item.weight > 0 && priorityQueue.Find(item.value) == nil && !checkers.StringInSlice(item.value, guessList) {
					log.Println("push item", item)
					item.weight = -item.weight
					heap.Push(&priorityQueue, item)
				}
			}

			// insert neighbours & password into the priority queue
			neighbours := getNeighbours(registeredPassword, topCorrectors)
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
				weight := power(neighbour, attackerList, blacklist, done)
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

			fblacklistIndex++
		} else {
			fmt.Println("out of options")
		}
	}

	lambdaQ := ballProbability(naiveGuessList, fblacklist)
	lambdaQFuzzy := ballProbability(guessList, fblacklist)
	log.Println("typo guess list:", guessList)
	// log.Println("normal guess list:", naiveGuessList)
	log.Println("lambda q", lambdaQ)
	log.Println("lambda q fuzzy", lambdaQFuzzy)
	log.Println("sec loss", lambdaQFuzzy-lambdaQ)
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

func passwordProbability(password string, frequencies map[string]int) float64 {
	probability, ok := frequencies[password]

	if ok {
		return float64(probability) / float64(totalNumberOfPasswords(frequencies))
	}
	return 0
}

func ballProbability(ball Ball, frequencies map[string]int) float64 {
	ballProbability := 0.0
	for _, password := range ball {
		ballProbability += passwordProbability(password, frequencies)
	}
	return ballProbability
}

// returns the union ball of passwords that are not in the done list
func unionBallNotDone(password string, done map[string]bool) []string {
	// TODO make get ball configurable according to the checker
	unionBallNotDone := make([]string, 0)
	temp := correctors.GetBall(password, topCorrectors)
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

func power(password string, attackerList map[string]int, blacklist []string, done map[string]bool) float64 {
	probability := 0.0

	// add passwords in ball to done
	unionBall := unionBallNotDone(password, done)

	for _, pw := range unionBall {
		probability += passwordProbability(pw, attackerList)
	}
	// log.Println(unionBall)
	// log.Println(probability)
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

func getNeighbours(password string, bestCorrectors []string) []string {
	neighbours := make([]string, 0)
	for _, corrector := range bestCorrectors {
		edits := correctors.ApplyInverseCorrectionFunction(corrector, password)
		neighbours = append(neighbours, edits...)
	}
	neighbours = correctors.DeleteEmpty(neighbours)
	return neighbours
}

func max(a float64, b int) float64 {
	if a >= float64(b) {
		return a
	}

	return float64(b)
}

func totalNumberOfPasswords(frequencies map[string]int) int {
	sum := 0
	for _, frequency := range frequencies {
		sum += frequency
	}
	return sum
}

func applyEdits(password string) []string {
	edits := make([]string, 10)

	// apply all correctors to password & add to edits
	// ball := correctors.GetBall(password, correctors.GetNBestCorrectors(10, typoFrequencies))
	// for _, passwordInBall := range ball {
	// 	edits = append(edits, passwordInBall)
	// }

	for _, letter := range letterRunes {
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

func getRandomPasswordFromBlacklist(blacklist []string) string {
	return blacklist[mrand.Intn(len(blacklist))-1]
}

func getRandomPasswordFromFrequencyBlacklist(blacklist []correctors.KeyValue) string {
	return blacklist[mrand.Intn(len(blacklist))-1].Key
}
