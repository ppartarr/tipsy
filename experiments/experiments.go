package main

import (
	"container/heap"
	"log"
	mrand "math/rand"
	"time"

	"github.com/ppartarr/tipsy/checkers"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`1234567890-=[]\\;',./~!@#$%^&*()_+{}|:\"<>?")

func main() {
	// seed randomness
	// mrand.Seed(time.Now().UnixNano())

	// sample from password leaks
	// blacklist := checkers.LoadBlacklist("../data/blacklistRockYou1000.txt")
	fblacklist := checkers.LoadFrequencyBlacklist("../data/rockyou-withcount1000.txt")

	// TODO get from args
	q := 10

	var (
		priorityQueue        PriorityQueue = make(PriorityQueue, 0)
		guessList            []string
		done                 []string
		passwordFrequencies  map[string]int = fblacklist
		length               int            = 1
		startTime            time.Time      = time.Now().UTC()
		ballSize             float64        = 4
		sortedBlacklistSlice                = correctors.ConvertMapToSortedSlice(fblacklist)
	)

	for registeredPassword, frequency := range fblacklist {
		// TODO don't check for password under 6 chars
		if len(registeredPassword) < 6 {
			continue
		}

		// TODO make this customisable for typtop
		// neighbours := applyEdits(registeredPassword)

		neighbours := correctors.GetBall(registeredPassword, []string{"swc-all", "rm-last", "swc-first"})

		// log.Println(priorityQueue.Len())

		// iterate over priorityQueue
		for priorityQueue.Len() > 0 {
			item := heap.Pop(&priorityQueue).(*Item)
			item.priority = -item.priority
			// ball := applyEdits(item.value)
			log.Println("here")
			log.Println(item)
			ball := correctors.GetBall(item.value, []string{"swc-all", "rm-last", "swc-first"})
			// TODO verify that this is what's needed
			ballFrequencySum := sum(ball, passwordFrequencies)

			if item.priority == ballFrequencySum {

				// log.Println("priority:", item.priority)
				// log.Println("frequency:", frequency)
				// log.Println("ballsize:", int(ballSize))
				// log.Println("freq * ballsize:", frequency*int(ballSize))

				if item.priority >= frequency*int(ballSize) {
					log.Println("Guess ", len(guessList), "/", q, "password: ", item.value, "weight: ", item.priority/totalFrequencies(fblacklist))

					done = append(done, item.value)
					guessList = append(guessList, item.value)
					// TODO set frequency of ball to 0
					// passwordFrequencies[] = 0
					if len(guessList) >= q {
						break
					}
				} else {
					priorityQueue.update(item, item.value, -ballFrequencySum)
					break
				}
			} else {
				priorityQueue.update(item, item.value, -ballFrequencySum)
			}
		}

		ballMax := 0.0

		log.Println(neighbours)
		for _, neighbour := range neighbours {
			// neighbourBall := applyEdits(neighbour)
			neighbourBall := correctors.GetBall(registeredPassword, []string{"swc-all", "rm-last", "swc-first"})
			// log.Println("adding to priority queue")
			// if priorityQueue.Len() == 0 {
			// 	// init priority queue
			// 	priorityQueue[0] = &Item{
			// 		value:    neighbour,
			// 		priority: -sum(neighbourBall, passwordFrequencies),
			// 		index:    0,
			// 	}

			// 	heap.Init(&priorityQueue)
			// } else {
			// 	item := &Item{
			// 		value:    neighbour,
			// 		priority: -sum(neighbourBall, passwordFrequencies),
			// 	}
			// 	priorityQueue.Push(item)
			// }
			item := &Item{
				value:    neighbour,
				priority: -sum(neighbourBall, passwordFrequencies),
			}
			priorityQueue.Push(item)
			log.Println("add to priority queue:", item)
			ballMax = max(ballMax, len(neighbourBall))
		}
		ballSize = ballSize*0.9 + ballMax*0.1
		// log.Println("updated ball size to:", ballSize)

		if len(priorityQueue) > length {
			// print update whenever heap size doubles
			log.Println(">< (", time.Now().Local().UTC().Sub(startTime), ") heap size: ", len(priorityQueue), " ballsize: ", ballSize)
			length = len(priorityQueue) * 2
		}
		// if index%10 == 0 {
		// 	log.Println("%t %d",
		// 		time.Now().Local().UTC().Sub(startTime),
		// 		index,
		// 		registeredPassword,
		// 		frequency)
		// }
		if len(guessList) >= q {
			break
		}
	}

	lambdaQ := float64(sumQMostProbablePasswords(sortedBlacklistSlice, q)) / float64(totalFrequencies(fblacklist))
	guessedPassword := make([]string, 0)
	for _, guess := range guessList {
		// for _, edit := range applyEdits(guess) {
		for _, edit := range correctors.GetBall(guess, []string{"swc-all", "rm-last", "swc-first"}) {
			guessedPassword = append(guessedPassword, edit)
		}
	}
	guessedPassword = funk.UniqString(guessedPassword)

	lambdaQFuzzy := float64(sum(guessedPassword, passwordFrequencies)) / float64(totalFrequencies(fblacklist))
	log.Println("lambdaQ:", lambdaQ, "lambdaQFuzzy: ", lambdaQFuzzy)
	log.Println(guessList)
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
