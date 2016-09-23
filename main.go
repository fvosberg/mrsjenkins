package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	elastic "gopkg.in/olivere/elastic.v3"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

type Todo struct {
	Title       string `validate:"nonzero",json:"title"`
	Description string `validate:"-",json:"description"`
}

type App struct {
	todoDatastore TodoDatastore
	mux           http.Handler
}

type TodoDatastore interface {
	Create(*Todo)
}

type TodoElasticsearchDatastore struct {
	client *elastic.Client
}

func (t *TodoElasticsearchDatastore) Create(todo *Todo) {
	log.Printf("TODO: should create Todo %+v\n", todo)
}

func NewApp() *App {
	app := NewAppWithTodoDatastore(
		NewTodoElasticDatastore(),
	)
	return app
}

func NewAppWithTodoDatastore(todoDatastore TodoDatastore) *App {
	app := &App{
		todoDatastore: todoDatastore,
	}
	app.initRouter()
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

func NewTodoElasticDatastore() TodoDatastore {
	return NewTodoElasticDatastoreWithClient(
		NewElasticsearchClient("http://elasticsearch.mrsjenkins.de:9200"),
	)
}

func NewTodoElasticDatastoreWithClient(elastic *elastic.Client) TodoDatastore {
	datastore := &TodoElasticsearchDatastore{
		client: elastic,
	}
	return datastore
}

func (a *App) initRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/todos", TodoHandler).Methods("GET")
	// TODO New TodoCreateHandler wird einfach a.TodosCreateHandler
	r.HandleFunc("/todos", NewTodoCreateHandler(a.todoDatastore)).Methods("PUT")
	a.mux = r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled request on %s\n", r.URL.Path)
	w.Write([]byte("Hallo Welt"))
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled request with %s on %s\n", r.Method, r.URL.Path)
	w.Write([]byte("Todos"))
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func todoFromRequest(r *http.Request) (Todo, error) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	return todo, err
}

func NewTodoCreateHandler(todoDatastore TodoDatastore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handled request with %s on %s\n", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		todo, err := todoFromRequest(r)
		if err == io.EOF {
			log.Println("Empty registration request body")
			w.Header().Set("X-Status-Reason", "Missing request body")
			w.WriteHeader(400)
		} else if err != nil {
			log.Printf("Decoding todo create request failed - %+v - %+v\n", err, r)
		} else {
			log.Printf("Decoding todo create request successful - %+v\n", todo)
			if err := validator.Validate(todo); err != nil {
				log.Printf("Validation error for todo: %+v\n", err)
				w.Header().Set("X-Status-Reason", "Validation failed; See body for reasons")
				w.WriteHeader(400)
				var responseBody ValidationErrorResponse
				for field, ves := range err.(validator.ErrorMap) {
					for _, ve := range ves {
						switch {
						default:
							log.Printf("TODO LOG STATUS CRITICAL: There is no error code definition for this validation error combination: %s, %s", field, ve)
						case field == "Title" && ve == validator.ErrZeroValue:
							validationError := ValidationError{Code: 1474574120, Description: "The field Title is required"}
							responseBody.Errors = append(responseBody.Errors, validationError)
						}
					}
				}
				jsonBody, err := json.Marshal(responseBody)
				if err != nil {
					log.Printf("Error while marshalling todos validation error response to json. - %+v - %+v\n", err, responseBody)
				}
				w.Write(jsonBody)
			} else {
				todoDatastore.Create(&todo)
				w.Write([]byte("Todo Created"))
			}
		}
	}
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
