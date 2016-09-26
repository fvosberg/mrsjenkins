package main

import (
	"net/http"

	"github.com/fvosberg/mrsjenkins/todo"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type app struct {
	todoDatastore todo.Datastore
	mux           http.Handler
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func NewApp() *app {
	return NewAppWithTodoDatastore(todo.NewElasticDatastore())
}

func NewAppWithTodoDatastore(todoDatastore todo.Datastore) *app {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.Handle("/todos", todo.NewCreateHandler(todoDatastore)).Methods("PUT")

	app := &app{
		todoDatastore: todoDatastore,
		mux:           r,
	}
	return app
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
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
