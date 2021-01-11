package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ppartarr/tipsy/web/session"
)

var (
	// url extensions to identify a static asset
	assetExtensions = []string{
		"css",
		"jpg",
		"jpeg",
		"html",
		"woff",
		"js",
		"json",
		"svg",
		"png",
		"webp",
		"webmanifest",
		"ico",
		"mkv",
		"mp4",
		"wmv",
		"mov",
		"MOV", // default for videos created on iOS
		"flv",
		"avi",
		"webm",
	}
)

// Server handles HTTP traffic from client
type Server struct {
	FileHandler *FileServer
	UserService *UserService
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// init http handlers
	// fs := http.FileServer(http.Dir("./static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.HandleFunc("/", s.FileHandler.serveHTTP)
	// http.HandleFunc("/login", s.UserService.Login)
	// http.HandleFunc("/register", s.UserService.Register)

	if isHTML(r.URL.Path) {
		fmt.Println(r.URL.Path)
		// if request is for home.html, check if user has a session

		// if request is for home page, check the user session
		if strings.HasSuffix(r.URL.Path, "home.html") {
			_, err := session.GetUserIDInt(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		// if request is for login page & the user has a valid session, redirect to home page
		// if strings.HasSuffix(r.URL.Path, "login.html") {
		// 	userID, err := session.GetUserIDInt(w, r)

		// 	if userID != -1 && err == nil {
		// 		// redirect to home page
		// 		fmt.Println("redirecting to home.html")
		// 		http.Redirect(w, r, "/home.html", 301)
		// 	}
		// }
		s.FileHandler.ServeHTTP(w, r)
	}

	if isCSS(r.URL.Path) {
		s.FileHandler.Handler.ServeHTTP(w, r)
	}

	switch r.URL.Path {
	// redirect to / => /login.hml
	case "/":
		fmt.Println("redirecting to login.html")
		http.Redirect(w, r, "/login.html", 301)
		return
	// serve robots file
	case "/robots.txt":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`User-agent: *
	Disallow: /`))
		return
	case "/login":
		_, err := s.UserService.Login(w, r)

		fmt.Println(err)
		if err == nil {
			http.Redirect(w, r, "/home.html", 301)
		}

		// TODO if login successfull redirect to home
		return
	case "/register":
		_, err := s.UserService.Register(w, r)
		fmt.Println(err)
		if err == nil {
			http.Redirect(w, r, "/login.html", 301)
		}
		return
	case "/logout":
		s.UserService.Logout(w, r)
		// _, err := s.UserService.Logout(w, r)
		// if err == nil {
		// 	http.Redirect(w, r, "/login.html", 301)
		// }
		return
	// serve favicon
	// case "/favicon.ico":
	// 	c, err := ioutil.ReadFile("favicon.ico")
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = w.Write(c)

	// 	return
	default:
	}
}

func isHTML(path string) bool {
	if strings.HasSuffix(path, "html") {
		return true
	}
	return false
}

func isCSS(path string) bool {
	if strings.HasSuffix(path, "css") {
		return true
	}
	return false
}

func isAsset(path string) bool {
	for _, e := range assetExtensions {
		if strings.HasSuffix(path, e) {
			return true
		}
	}
	return false
}
