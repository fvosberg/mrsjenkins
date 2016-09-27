package main

// Purpose of this testfile
// Test the correct linking of the app router with the todo router
// Test the correct 404s

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fvosberg/mrsjenkins/todo"
)

type fakeTodoDatastore struct {
	calledCreateCount int
}

func (f *fakeTodoDatastore) Create(t *todo.Todo) {
	f.calledCreateCount++
}

func (f *fakeTodoDatastore) assertCreateCalled(t *testing.T, expected int) {
	if f.calledCreateCount != expected {
		t.Error("TodoDatastore.Create should be called", expected, "times, but has been called", f.calledCreateCount)
	}
}

func assertStatusCode(t *testing.T, expected int, actual int, path string) {
	if expected != actual {
		t.Error("Status code of ", path, " is not ", expected, " but", actual)
	}
}

func assertStatusCodeNot(t *testing.T, expected int, actual int, path string) {
	if expected == actual {
		t.Error("Status code of ", path, " should not be ", actual)
	}
}

func TestMain(m *testing.M) {
	//logrus.SetOutput(ioutil.Discard)
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.Debugf("HALLO %+v\n", m)
	retCode := m.Run()
	//myTeardownFunction()
	os.Exit(retCode)
}

func TestRouter(t *testing.T) {
	app := NewApp()

	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL)
	if err != nil {
		t.Error(err)
	}

	assertStatusCode(t, 200, resp.StatusCode, "GET /")
}

func Test404(t *testing.T) {
	app := NewApp()
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL + "/notExistingRoute")
	if err != nil {
		t.Error(err)
	}

	assertStatusCode(t, 404, resp.StatusCode, "GET /notExistingRoute")
}

func TestSub404(t *testing.T) {
	app := NewApp()
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL + "/todos/notExistingRoute")
	if err != nil {
		t.Error(err)
	}

	assertStatusCode(t, 404, resp.StatusCode, "GET /todos/notExistingRoute")
}

func TestTodoRouting(t *testing.T) {
	for _, r := range todo.Routes {
		assertRouteExists(t, r)
	}

	notExistingMethodRoute := todo.Route{"/", "NOTEXISTINGMETHOD", func(w http.ResponseWriter, r *http.Request) {}}
	assertRouteDoesntExist(t, notExistingMethodRoute)
}

func assertRouteExists(t *testing.T, route todo.Route) {
	resp, err := requestRoute(route)
	if err != nil {
		t.Error(err)
	}

	routeS := route.Method + " /todos" + route.URL
	assertStatusCodeNot(t, 404, resp.StatusCode, routeS)
}

func assertRouteDoesntExist(t *testing.T, route todo.Route) {
	resp, err := requestRoute(route)
	if err != nil {
		t.Error(err)
	}

	routeS := route.Method + " /todos" + route.URL
	assertStatusCode(t, 404, resp.StatusCode, routeS)
}

func requestRoute(route todo.Route) (*http.Response, error) {
	s := httptest.NewServer(NewApp())
	defer s.Close()

	req, err := http.NewRequest(route.Method, s.URL+"/todos"+route.URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}
