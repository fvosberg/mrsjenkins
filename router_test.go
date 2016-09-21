package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	app := NewApp()
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
	app := NewApp()
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL + "/foobarNotFound")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 404 {
		t.Error("Status code of /foobarNotFound is not 404")
	}
}

func TestGetTodos(t *testing.T) {
	app := NewApp()
	s := httptest.NewServer(app)
	defer s.Close()

	resp, err := http.Get(s.URL + "/todos")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Status code of GET /todos is not 200 but", resp.StatusCode)
	}
}

func TestPutTodos(t *testing.T) {
	app := NewAppWithTodoDatastore(
		NewTodoElasticDatastore(
			NewElasticsearchClient("http://192.168.99.100:9200"),
		),
	)

	s := httptest.NewServer(app)
	defer s.Close()

	client := &http.Client{}

	req, err := http.NewRequest("PUT", s.URL+"/todos", nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 404 {
		t.Error("Status code of GET /todos is not 404 but", resp.StatusCode)
	}
}
