package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dreadl0ck/musig/intern/pkg/db"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

func server() {
	// setup & open bolt database
	var (
		boltDB    *bolt.DB
		usersPath = "db/users.bolt"
	)

	boltDB, err := bolt.Open(usersPath, 0666, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// init http handlers
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHTTP)
	http.HandleFunc("/login", login, boltDB)
	http.HandleFunc("/register", login, boltDB)

	log.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	layoutPage := filepath.Join("static", "templates", "layout.html")
	filePath := filepath.Join("static", "templates", filepath.Clean(r.URL.Path))

	fmt.Println(layoutPage)
	fmt.Println(filePath)
	fmt.Println(r.URL.Path)

	switch r.URL.Path {
	// redirect to / => /login.hml
	case "/":
		http.Redirect(w, r, "/login.html", 301)
		return
	// serve robots file
	case "/robots.txt":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`User-agent: *
	Disallow: /`))
		return
	// serve favicon
	case "/favicon.ico":
		c, err := ioutil.ReadFile("favicon.ico")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(c)

		return
	default:
	}

	// return a 404 if the page doesn't exist
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	// generate page from layout
	tmpl, err := template.ParseFiles(layoutPage, filePath)
	if err != nil {
		log.Println(err.Error())

		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

// Login allows a user to login to their account
func Login(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	fmt.Println("method: ", r.Method)

	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return
	}

	r.ParseForm()
	fmt.Println("username: ", r.Form["username"])
	fmt.Println("password: ", r.Form["password"])
}

// Register allows a user to register a new account
func Register(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// TODO handle HEAD, PUT, and PATCH separately
	if r.Method != "POST" {
		http.Redirect(w, r, "/login.html", 301)
		return
	}

	r.ParseForm()
	fmt.Println(r.Form)

	passwordHash, err := HashPassword(r.Form["password"])
	if err != nil {
		fmt.Println(err.Error())
	}

	// create new user from request
	user = &User{
		Username:     r.Form["username"],
		PasswordHash: passwordHash,
	}

	// db.Update()
}

type User struct {
	ID           int    `json:"id" storm:"id,increment"`
	Username     string `json:"username" storm:"unique"`
	PasswordHash string `json:"password"`
}

type UserRegistration struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
