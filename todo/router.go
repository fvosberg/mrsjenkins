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
		{"/", "PUT", http.HandlerFunc(createHandle)},
	}
)

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

func createHandle(w http.ResponseWriter, r *http.Request) {
	logrus.Print("Create handler")
}
