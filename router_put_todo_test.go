package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/fvosberg/mrsjenkins/todo"
)

func TestMain(m *testing.M) {
	logrus.SetOutput(ioutil.Discard)
	// logrus.SetLevel(logrus.DebugLevel)
	// logrus.Debugf("HALLO %+v\n", m)
	retCode := m.Run()
	//myTeardownFunction()
	os.Exit(retCode)
}

func TestPutTodos(t *testing.T) {
	app, datastore := newAppWithFakeTodoDatastore()
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := requestTodo("PUT", s.URL+"/todos", todo.Todo{Title: "Test Todo", Description: "Description of Test Todo"})
	if err != nil {
		t.Error(err)
	}

	assertStatusCode(t, 200, resp.StatusCode, "PUT /todos")
	datastore.assertCreateCalled(t, 1)
}

func TestPutTodosWithoutTitle(t *testing.T) {
	app, datastore := newAppWithFakeTodoDatastore()

	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := requestTodo("PUT", s.URL+"/todos", todo.Todo{Title: "", Description: "Description of Test Todo"})
	if err != nil {
		t.Error(err)
	}
	assertStatusCode(t, 400, resp.StatusCode, "PUT /todos with invalid todo")
	assertHeader(t, resp, "X-Status-Reason", "Validation failed; See body for reasons")

	datastore.assertCreateCalled(t, 0)
	assertHeader(t, resp, "Content-Type", "application/json; charset=utf-8")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	var validationErrorResponse todo.ValidationErrorResponse
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

func TestPutTodosWithoutTodoBody(t *testing.T) {
	app, datastore := newAppWithFakeTodoDatastore()

	s := httptest.NewServer(app)
	defer s.Close()

	req, err := http.NewRequest("PUT", s.URL+"/todos", nil)
	if err != nil {
		t.Error(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	datastore.assertCreateCalled(t, 0)
	assertStatusCode(t, 400, resp.StatusCode, "PUT /todos without todo body")
	assertHeader(t, resp, "X-Status-Reason", "Missing request body")
}
