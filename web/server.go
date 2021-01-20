package web

import (
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ppartarr/tipsy/web/session"
	"github.com/ppartarr/tipsy/web/users"
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
	UserService *users.UserService
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if isHTML(r.URL) {
		log.Println(r.URL.Path)
		// if request is for home.html, check if user has a session

		// if request is for home page, check the user session
		if strings.HasSuffix(r.URL.Path, "home.html") {
			_, err := session.GetUserIDInt(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		}

		// redirect reset.html to login.html if there is no token
		if r.URL.String() == "/reset.html" {
			http.Error(w, "cannot access reset page without a token", http.StatusUnauthorized)
			return
		}

		// if request is for login page & the user has a valid session, redirect to home page
		// if strings.HasSuffix(r.URL.Path, "login.html") {
		// 	userID, err := session.GetUserIDInt(w, r)

		// 	if userID != -1 && err == nil {
		// 		// redirect to home page
		// 		log.Println("redirecting to home.html")
		// 		http.Redirect(w, r, "/home.html", 301)
		// 	}
		// }
		s.FileHandler.ServeHTTP(w, r)
		return
	}

	if isAsset(r.URL.Path) {
		s.FileHandler.Handler.ServeHTTP(w, r)
		return
	}

	log.Println(r.URL.Path)

	switch r.URL.Path {
	// redirect to / => /login.hml
	case "/":
		log.Println("redirecting to login.html")
		http.Redirect(w, r, "/login.html", 301)
		return
	// serve robots file
	case "/robots.txt":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`User-agent: *
	Disallow: /`))
		return
	case "/login":
		form, err := s.UserService.Login(w, r)

		if err != nil {
			if err.Error() == "you must submit a valid form" {
				loginHTML := filepath.Join("static", "templates", "/login.html")
				render(w, loginHTML, form)
				return
			}
			log.Println(err.Error())
			loginHTML := filepath.Join("static", "templates", "/login.html")
			render(w, loginHTML, nil)
			return
		}
		homeHTML := filepath.Join("static", "templates", "/home.html")
		render(w, homeHTML, nil)
		return
	case "/register":
		form, err := s.UserService.Register(w, r)
		if err != nil {
			log.Println(err.Error())
			if err.Error() == "you must submit a valid form" {
				registerHTML := filepath.Join("static", "templates", "/registration.html")
				render(w, registerHTML, form)
				return
			}
		} else {
			homeHTML := filepath.Join("static", "templates", "/home.html")
			render(w, homeHTML, nil)
		}

		// http.Redirect(w, r, "/login.html", 301)
		return
	case "/logout":
		err := s.UserService.Logout(w, r)
		if err != nil {
			log.Println(err.Error())
		}
		loginHTML := filepath.Join("static", "templates", "/login.html")
		render(w, loginHTML, nil)
		return
	case "/recover":
		err := s.UserService.PasswordRecovery(w, r)
		if err != nil {
			log.Println(err.Error())
		}
		resetHTML := filepath.Join("static", "templates", "/reset.html")
		render(w, resetHTML, nil)
		return
	case "/reset":
		err := s.UserService.PasswordReset(w, r)
		if err != nil {
			log.Println(err.Error())
		}
		loginHTML := filepath.Join("static", "templates", "/login.html")
		render(w, loginHTML, nil)
		return
	default:
	}
}

func isHTML(url *url.URL) bool {
	slice := strings.Split(url.String(), "?")

	if len(slice) > 0 && strings.HasSuffix(slice[0], "html") {
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
