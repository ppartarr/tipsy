package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func server() {
	fmt.Println(os.Args)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveHTTP)

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
