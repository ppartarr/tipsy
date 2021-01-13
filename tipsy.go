package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/ppartarr/tipsy/checkers"
	"github.com/ppartarr/tipsy/config"
	"github.com/ppartarr/tipsy/mail"
	"github.com/ppartarr/tipsy/web"
	"github.com/ppartarr/tipsy/web/session"
	bolt "go.etcd.io/bbolt"
)

const (
	version = "0.0.1"
	domain  = "typo.partarrieu.me"
	email   = "philippe@partarrieu.me"
)

var (
	sessionKey = os.Getenv("SESSION_KEY")
)

func main() {
	// get args from cli
	args := os.Args[1:]
	submittedPassword := args[0]
	configFile := "tipsy.yml"

	// load server config
	tipsyConfig, err := config.LoadServer(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(tipsyConfig.Checker.Always)

	// load black list
	// blackList := loadBlackList("./data/blacklistRockYou1000.txt")
	// log.Println(blackList)
	// CheckBlacklist(password, blacklist)

	// load frequency black list
	frequencyBlacklist := loadFrequencyBlackList("./data/rockyou-withcount1000.txt")
	// frequencyBlacklist := loadFrequencyBlackList("./data/blacklistTest.txt")
	checkers.CheckOptimal(submittedPassword, frequencyBlacklist)

	// setup & open bolt database
	var (
		sessionDB   bleve.Index
		boltDB      *bolt.DB
		usersDBPath = "db/users.bolt"
	)

	boltDB, err = bolt.Open(usersDBPath, 0666, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer boltDB.Close()

	// init mailer
	if stat, err := os.Stat(configFile); err == nil && !stat.IsDir() {
		if tipsyConfig.SMTP != nil {
			mail.InitMailer(
				tipsyConfig.SMTP.Server,
				tipsyConfig.SMTP.Username,
				tipsyConfig.SMTP.Password,
				tipsyConfig.SMTP.From,
				tipsyConfig.SMTP.Port,
			)
		}
	} else {
		log.Fatal("specified configuration file " + configFile + " does not exist or is a directory")
	}

	// init session
	log.Println("initializing cookie store, cookies will expire after: ", time.Duration(60)*time.Second)
	session.StorageDir = "./db"
	session.InitStore(sessionKey, sessionDB, 60)

	// create the server instance
	var server *web.Server

	server = &web.Server{
		FileHandler: &web.FileServer{
			Handler: http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
		},
		UserService: web.NewUserService(boltDB, tipsyConfig),
	}

	// start listening to requests
	log.Println("Listening on :8000...")

	err = http.ListenAndServe(":8000", server)
	if err != nil {
		log.Fatal(err)
	}
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
