package main

// TODO disable logging on running tests
// Purpose of this testfile
// Test the correct linking of the app router with the todo router
// So a test per route should be fine

import (
	"bytes"
	"encoding/json"
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

func newFakeTodoDatastore() *fakeTodoDatastore {
	return &fakeTodoDatastore{}
}

func newAppWithFakeTodoDatastore() (*app, *fakeTodoDatastore) {
	fakeTodoDatastore := newFakeTodoDatastore()
	return NewAppWithTodoDatastore(fakeTodoDatastore), fakeTodoDatastore
}

func requestTodo(method string, URL string, todo todo.Todo) (*http.Response, error) {
	json, err := json.Marshal(todo)
	if err != nil {
		return nil, err
	}
	jsonStr := []byte(json)
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
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
	fakeTodoDatastore := newFakeTodoDatastore()
	app := NewAppWithTodoDatastore(fakeTodoDatastore)
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Status code of / is not 200")
	}
}

func Test404(t *testing.T) {
	fakeTodoDatastore := newFakeTodoDatastore()
	app := NewAppWithTodoDatastore(fakeTodoDatastore)
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL + "/foobarNotFound")
	if err != nil {
		t.Error(err)
	}

	assertStatusCode(t, 404, resp.StatusCode, "GET /foobarNotFound")
}

func TestGetTodos(t *testing.T) {
	return
	fakeTodoDatastore := newFakeTodoDatastore()
	app := NewAppWithTodoDatastore(fakeTodoDatastore)
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL + "/todos")
	if err != nil {
		t.Error(err)
	}

	assertStatusCode(t, 200, resp.StatusCode, "GET /todos")
}

func assertHeader(t *testing.T, resp *http.Response, header string, expected string) {
	if resp.Header.Get(header) != expected {
		t.Error("The", header, "header should be", expected, "but is:", resp.Header.Get(header))
	}
}

func assertStatusCode(t *testing.T, expected int, actual int, path string) {
	if expected != actual {
		t.Error("Status code of ", path, " is not ", expected, " but", actual)
	}
}
