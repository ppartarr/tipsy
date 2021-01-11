package web

import (
	"net/http"
	"strings"
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
		s.FileHandler.ServeHTTP(w, r)
	}

	if isCSS(r.URL.Path) {
		s.FileHandler.Handler.ServeHTTP(w, r)
	}

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
	case "/login":
		s.UserService.Login(w, r)
		return
	case "/register":
		s.UserService.Register(w, r)
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
