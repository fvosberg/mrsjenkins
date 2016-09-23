package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

type fakeTodoDatastore struct {
	calledCreateCount int
}

func (f *fakeTodoDatastore) Create(t *Todo) {
	f.calledCreateCount++
}

func newFakeTodoDatastore() *fakeTodoDatastore {
	return &fakeTodoDatastore{}
}

func newAppWithFakeTodoDatastore() (*App, *fakeTodoDatastore) {
	fakeTodoDatastore := newFakeTodoDatastore()
	return NewAppWithTodoDatastore(fakeTodoDatastore), fakeTodoDatastore
}

func TestPutTodos(t *testing.T) {
	app, datastore := newAppWithFakeTodoDatastore()
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := requestTodo("PUT", s.URL+"/todos", Todo{Title: "Test Todo", Description: "Description of Test Todo"})
	if err != nil {
		t.Error(err)
	}

	assertStatusCode(t, 200, resp.StatusCode, "PUT /todos")
	if datastore.calledCreateCount != 1 {
		t.Error("TodoDatastore.Create should be called", 1, "times, but has been called", datastore.calledCreateCount)
	}
}

func requestTodo(method string, URL string, todo Todo) (*http.Response, error) {
	// TODO refactor json generation
	jsonprep := `{"title":"` + todo.Title + `","description":"` + todo.Description + `"}`
	jsonStr := []byte(jsonprep)
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func TestPutTodosWithoutTitle(t *testing.T) {
	app, datastore := newAppWithFakeTodoDatastore()

	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := requestTodo("PUT", s.URL+"/todos", Todo{Title: "", Description: "Description of Test Todo"})
	if err != nil {
		t.Error(err)
	}
	assertStatusCode(t, 400, resp.StatusCode, "PUT /todos with invalid todo")
	if datastore.calledCreateCount != 0 {
		t.Error("TodoDatastore.Create should be called", 0, "times, but has been called", datastore.calledCreateCount)
	}
	assertHeader(t, resp, "Content-Type", "application/json; charset=utf-8")
	assertHeader(t, resp, "X-Status-Reason", "Validation failed; See body for reasons")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	var validationErrorResponse ValidationErrorResponse
	err = json.Unmarshal(body, &validationErrorResponse)
	if err != nil {
		t.Error(err)
	}
	if len(validationErrorResponse.Errors) != 1 {
		t.Error("Expected the response to have 1 error, but had", len(validationErrorResponse.Errors))
		t.Error("Expected the response to have the 'title required error' with the error code 1474574120")
	} else if validationErrorResponse.Errors[0].Code != 1474574120 {
		t.Error("Expected the response to have the 'title required error' with the error code 1474574120")
	}
}

func assertHeader(t *testing.T, resp *http.Response, header string, expected string) {
	if resp.Header.Get(header) != expected {
		t.Error("The", header, "header should be", expected, "but is:", resp.Header.Get(header))
	}
}

// TODO func TestPutTodosWithoutRequestBody(t *testing.T) {

func assertStatusCode(t *testing.T, expected int, actual int, path string) {
	if expected != actual {
		t.Error("Status code of ", path, " is not ", expected, " but", actual)
	}
}
