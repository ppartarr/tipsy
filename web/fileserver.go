package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// FileServer implements simple HTTP file server
// that can serve a custom 404 page
type FileServer struct {

	// http handler to serve files
	Handler http.Handler
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("static", "templates", filepath.Clean(r.URL.Path))

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

	// generate page from layoutif err !=
	render(w, filePath, nil)
}

// Render renders the HTML file based on the layout template & the interface
func render(w http.ResponseWriter, filename string, data interface{}) {
	layoutPage := filepath.Join("static", "templates", "layout.html")

	tmpl, err := template.ParseFiles(layoutPage, filename)
	if err != nil {
		log.Println(err)
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
