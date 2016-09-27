package todo

import (
	"bytes"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type Route struct {
	URL     string
	Method  string
	Handler http.Handler
}

func (r *Route) String() string {
	var b bytes.Buffer

	b.WriteString(r.Method)
	b.WriteString(" ")
	b.WriteString(r.URL)

	return b.String()
}

var (
	Routes = []Route{
		{"/", "GET", http.HandlerFunc(listHandle)},
		{"/", "PUT", NewCreateHandler(&SessionDatastore{})},
	}
)

// NewRouter returns a new router for the todo service
// this router routes on / level, so you have to remove the prefix if it is used embedded
func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(false)
	for _, route := range Routes {
		r.Handle(route.URL, route.Handler).Methods(route.Method)
	}
	return r
}

func listHandle(w http.ResponseWriter, r *http.Request) {
	logrus.Print("List handler")
}
