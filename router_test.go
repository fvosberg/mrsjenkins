package main

// Purpose of this testfile
// Test the correct linking of the app router with the todo router
// Test the correct 404s

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/fvosberg/mrsjenkins/todo"
)

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
	logrus.SetLevel(logrus.DebugLevel)
	// logrus.Debugf("HALLO %+v\n", m)
	retCode := m.Run()
	//myTeardownFunction()
	os.Exit(retCode)
}

const TodoRoutingPrefix = "/todos"

func TestRouter(t *testing.T) {
	cases := map[string]struct {
		Method             string
		URL                string
		ExpectedStatusCode int
	}{
		"Root":   {"GET", "/", 200},
		"404":    {"GET", "/notExistingRoute", 404},
		"Sub404": {"GET", TodoRoutingPrefix + "/notExistingRoute", 404},
	}

	for k, tc := range cases {
		resp := getResponse(t, tc.Method, tc.URL)
		if resp.StatusCode != tc.ExpectedStatusCode {
			t.Error("Failed", k, "- Status code of", tc.Method, tc.URL, "is not", tc.ExpectedStatusCode, "but", resp.StatusCode)
		}
	}
}

func TestTodoRoutes(t *testing.T) {
	for _, r := range todo.Routes {
		URL := TodoRoutingPrefix + r.URL
		resp := getResponse(t, r.Method, URL)
		if resp.StatusCode == 404 {
			t.Error("Failed - Status code of", r.Method, URL, "is not SUCCESS but", resp.StatusCode)
		}
	}

	resp := getResponse(t, "NOTEXISTINGMETHOD", TodoRoutingPrefix+"/")
	if resp.StatusCode != 404 {
		t.Error("Failed - Status code of the root rout with a non existing HTTP method is not 404, but", resp.StatusCode)
	}
}

func getResponse(t *testing.T, method string, URL string) *http.Response {
	app := NewApp()

	s := httptest.NewServer(app)
	defer s.Close()

	req, err := http.NewRequest(method, s.URL+URL, nil)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	return resp
}
