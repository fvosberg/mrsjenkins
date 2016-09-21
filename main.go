package main

import (
	"log"
	"net/http"
	"os"
	"time"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/gorilla/mux"
)

type App struct {
	todoDatastore TodoDatastore
	mux           http.Handler
}

type TodoDatastore interface {
}

type TodoElasticsearchDatastore struct {
	client *elastic.Client
}

func NewApp() *App {
	app := &App{
		mux: NewRouter(),
	}
	return app
}

func NewAppWithTodoDatastore(todoDatastore TodoDatastore) *App {
	app := NewApp()
	app.todoDatastore = todoDatastore
	return app
}

func NewElasticsearchClient(URL string) *elastic.Client {
	var client *elastic.Client
	var err error
	log.Printf("Trying to initialize an Elasticsearch client on \"%s\"\n", URL)
	for {
		client, err = elastic.NewClient(
			elastic.SetURL(URL),
			elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		)
		if err != nil {
			log.Printf("Error while connecting to elasticsearch on \"%s\": %+v - retrying in 5 seconds\n", URL, err)
			time.Sleep(5 * time.Second)
		} else {
			log.Printf("Initialized Elasticsearch client on \"%s\"\n", URL)
			break
		}
	}
	return client
}

func NewTodoElasticDatastore(elastic *elastic.Client) TodoDatastore {
	datastore := &TodoElasticsearchDatastore{
		client: elastic,
	}
	return datastore
}

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/todos", CustomerHandler).Methods("GET")

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

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
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
