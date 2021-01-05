package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// get args from cli
	// args := os.Args[1:]
	// submittedPassword := args[0]
	// fmt.Println(args[0])

	// numberOfCorrectors := 3go
	// checker := "always"
	password := "Password"

	// load black list
	// blackList := loadBlackList("./data/blacklistRockYou1000.txt")
	// fmt.Println(blackList)
	// CheckBlacklist(password, blacklist)

	// load frequency black list
	fmt.Println(290729 / 14344391)
	frequencyBlacklist := loadFrequencyBlackList("./data/blacklistTest.txt")
	fmt.Println(frequencyBlacklist)
	CheckOptimal(password, frequencyBlacklist)

}

func loadBlackList(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(string(content), "\n")
}

func loadFrequencyBlackList(filename string) map[string]int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data map[string]int = make(map[string]int)

	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), " ")
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
