package main

import (
	"log"
	"net/http"
)

type Router struct {
	Message string
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled request on %s\n", r.URL.Path)
	w.Write([]byte(router.Message))
}

func main() {
	router := &Router{Message: "Hello, World"}

	log.Println("Starting webserver listening on :8080")
	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(s.ListenAndServe())
}
