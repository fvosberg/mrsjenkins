package main

import (
	"net/http"

	"github.com/fvosberg/mrsjenkins/todo"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// The App struct holds the todoDatastore and the multiplexer
type App struct {
	todoDatastore todo.Datastore
	mux           http.Handler
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NewRouter and Subrouter return Router which implements http.handler
	a.mux.ServeHTTP(w, r)
}

// NewApp creates a new App object
func NewApp() *App {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.NewRoute().PathPrefix("/todos").Handler(
		http.StripPrefix("/todos", todo.NewRouter()),
	)

	app := &App{
		mux: r,
	}
	return app
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled request on %s\n", r.URL.Path)
	w.Write([]byte("Hallo Welt"))
}

func main() {
	app := NewApp()

	log.Println("Starting webserver listening on :8080")
	s := &http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	log.Fatal(s.ListenAndServe())
}
