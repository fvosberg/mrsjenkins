package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// this file does not test the handlers
// it just tests the routing
// If all tests of this file pass we should be shure that:
// * The routes of this package end up in the right handler
// * Routes that doesn't exist replying with a 404

func TestRoutes(t *testing.T) {
	for _, r := range Routes {
		assertRouteExists(t, r)
	}

	notExistingMethodRoute := Route{"/", "NOTEXISTINGMETHOD", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})}
	assertRouteDoesntExist(t, notExistingMethodRoute)
}

func TestRouteString(t *testing.T) {
	r := &Route{Method: "GET", URL: "/foobar"}
	if r.String() != "GET /foobar" {
		t.Error("The String() method of routes should return it's method and it's URL.")
	}
}

func assertRouteExists(t *testing.T, route Route) {
	resp, err := requestRoute(route)
	if err != nil {
		t.Error(err)
	}

	routeS := route.Method + " /todos" + route.URL
	assertStatusCodeNot(t, 404, resp.StatusCode, routeS)
}

func assertRouteDoesntExist(t *testing.T, route Route) {
	resp, err := requestRoute(route)
	if err != nil {
		t.Error(err)
	}

	routeS := route.Method + " /todos" + route.URL
	assertStatusCode(t, 404, resp.StatusCode, routeS)
}

func requestRoute(route Route) (*http.Response, error) {
	s := httptest.NewServer(NewRouter())
	defer s.Close()

	req, err := http.NewRequest(route.Method, s.URL+route.URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
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
