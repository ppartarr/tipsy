package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/ppartarr/tipsy/config"
	"github.com/ppartarr/tipsy/mail"
	"github.com/ppartarr/tipsy/web"
	"github.com/ppartarr/tipsy/web/session"
	"github.com/ppartarr/tipsy/web/users"
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
	configFile := "tipsy.yml"

	// load server config
	tipsyConfig, err := config.LoadServer(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(tipsyConfig.Checker.Always)

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
	maxAge := int(tipsyConfig.HTTPSessionValidity / time.Second)
	log.Println("initializing cookie store, cookies will expire after: ", maxAge)
	session.StorageDir = "./db"
	session.InitStore(sessionKey, sessionDB, maxAge)

	// create the server instance
	var server *web.Server

	server = &web.Server{
		FileHandler: &web.FileServer{
			Handler: http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
		},
		UserService: users.NewUserService(boltDB, tipsyConfig),
	}

	// start listening to requests
	log.Println("Listening on :8000...")

	err = http.ListenAndServe(":8000", server)
	if err != nil {
		log.Fatal(err)
	}
}
