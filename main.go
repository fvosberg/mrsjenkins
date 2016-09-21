package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/todos", CustomerHandler)
	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled request on %s\n", r.URL.Path)
	w.Write([]byte("Hallo Welt"))
}

func CustomerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled request on %s\n", r.URL.Path)
	w.Write([]byte("Todos"))
}

func main() {
	r := NewRouter()

	log.Println("Starting webserver listening on :8080")
	s := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatal(s.ListenAndServe())
}
