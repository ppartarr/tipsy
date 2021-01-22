package main

import (
	"container/heap"
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

func main() {
	// seed randomness
	// mrand.Seed(time.Now().UnixNano())

	// sample from password leaks
	blacklist := checkers.LoadBlacklist("../data/rockyou1000.txt")
	fblacklist := checkers.LoadFrequencyBlacklist("../data/rockyou-withcount1000.txt")
	attackerList := checkers.LoadFrequencyBlacklist("../data/rockyou-withcount1000.txt")

	// TODO get from args
	q := 10

	var (
		priorityQueue       PriorityQueue = make(PriorityQueue, 0)
		guessList           []string
		naiveGuessList      []string
		tdone               map[string]bool = make(map[string]bool)
		rdone               map[string]bool = make(map[string]bool)
		passwordFrequencies map[string]int  = fblacklist
		length              int             = 1
		// startTime            time.Time      = time.Now().UTC()
		ballSize           float64 = 3
		sortedAttackerList         = correctors.ConvertMapToSortedSlice(attackerList)
		fblacklistIndex    int     = 0
		ballWeight         float64 = 0
	)

	log.Println("starting loop")

	for len(guessList) < q {
		// get the next most probable password from the attacker's password list
		registeredPassword := sortedAttackerList[fblacklistIndex].Key

		if len(registeredPassword) < 6 {
			fblacklistIndex++
			continue
		}

		if len(naiveGuessList) < q {
			naiveGuessList = append(naiveGuessList, registeredPassword)
		}

		if ballWeight <= 0 {
			ballWeight = float64(passwordFrequencies[registeredPassword]) * ballSize
		}

		if priorityQueue.Len() > 0 {
			item := heap.Pop(&priorityQueue).(*Item)
			weight := item.weight
			for weight > int(ballSize)*passwordFrequencies[registeredPassword] && len(guessList) < q {
				// add guess to guess list
				log.Println("Guess", len(guessList), "/", q, "password:", item.value, "weight:", float64(item.weight)/float64(totalFrequencies(fblacklist)))
				guessList = append(guessList, item.value)
				tdone[item.value] = true

				// add password & ball to done
				killed := unionBallNotDone(item.value, rdone)
				for _, password := range killed {
					rdone[password] = true
				}

				// add all neighbours of this password to the priority queue
				for _, password := range killed {
					probability := attackerList[password]
					neighbours := getNeighbours(password, topCorrectors)
					neighbours = append(neighbours, password)
					for _, neighbour := range neighbours {
						// update neighbour weight in the priority queue
						neighbourItem := priorityQueue.Find(neighbour)
						if neighbourItem != nil {
							priorityQueue.update(neighbourItem, neighbourItem.value, neighbourItem.weight-probability)

							// remove from priority queue if weight <= 0
							// log.Println("removing from priority q")
							// if neighbourItem.weight <= 0 {
							// 	heap.Remove(&priorityQueue, neighbourItem.index)
							// }
						}
					}
				}

				if priorityQueue.Len() > 0 {
					item = heap.Pop(&priorityQueue).(*Item)
					weight = item.weight
				}
			}
			// TODO add check here
		}

		// insert neighbours into the priority queue
		neighbours := getNeighbours(registeredPassword, topCorrectors)
		allNeighbours := neighbours
		neighbours = append(neighbours, registeredPassword)
		for _, neighbour := range neighbours {
			// don't add neighbour if it's already in the priority queue
			if priorityQueue.Find(neighbour) != nil {
				allNeighbours = remove(allNeighbours, neighbour)
			}
			// don't add neighbour to priority queue if it's already been tested
			_, ok := rdone[neighbour]
			if ok {
				allNeighbours = remove(allNeighbours, neighbour)
			}
		}

		// add items to the priority queue
		for _, neighbour := range allNeighbours {
			weight := power(neighbour, attackerList, blacklist, tdone)
			item := &Item{
				value:  neighbour,
				weight: weight,
			}
			heap.Push(&priorityQueue, item)
		}

		// print heap size update
		if priorityQueue.Len() > length {
			log.Println("Heap size:", priorityQueue.Len())
			length = priorityQueue.Len() * 2
		}

		fblacklistIndex++
	}

	log.Println("typo guess list:", guessList)
	log.Println("normal guess list:", naiveGuessList)
}

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

func unionBall(password string) []string {
	// TODO make get ball configurable according to the checker
	ball := correctors.GetBall(password, topCorrectors)
	return append(ball, password)
}

func power(password string, attackList map[string]int, blacklist []string, done map[string]bool) int {
	probability := 0

	// add passwords in ball to done
	unionBall := unionBall(password)

	for _, pw := range unionBall {
		_, ok := done[pw]
		if !ok {
			probability += attackList[pw]
		}
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

func getNeighbours(password string, bestCorrectors []string) []string {
	neighbours := make([]string, 0)
	for _, corrector := range bestCorrectors {
		edits := correctors.ApplyInverseCorrectionFunction(corrector, password)
		neighbours = append(neighbours, edits...)
	}
	neighbours = checkers.DeleteEmpty(neighbours)
	return neighbours
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func find(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func max(a float64, b int) float64 {
	if a >= float64(b) {
		return a
	}

	return float64(b)
}

// sum all frequencies of the passwords in the ball
func sum(ball Ball, frequencies map[string]int) int {
	sum := 0
	for _, password := range ball {
		sum += frequencies[password]
	}

	return sum
}

func totalFrequencies(frequencies map[string]int) int {
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

func sumQMostProbablePasswords(blacklist []correctors.KeyValue, q int) int {
	sum := 0
	for i := 0; i < q; i++ {
		// fmt.Println(blacklist[i])
		if i >= q {
			return sum
		}
		sum += blacklist[i].Value
	}

	return sum
}

func getRandomPasswordFromBlacklist(blacklist []string) string {
	return blacklist[mrand.Intn(len(blacklist))-1]
}

func getRandomPasswordFromFrequencyBlacklist(blacklist []correctors.KeyValue) string {
	return blacklist[mrand.Intn(len(blacklist))-1].Key
}
