package todo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHandlerWithMalformedJSON(t *testing.T) {
	handler := NewCreateHandler(&SessionDatastore{})
	r, err := recordedResponse(handler, "PUT", `{\"foo\": \"bar\"`)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 500, r.Code, "The status code should be 500")
	assert.Equal(
		t,
		r.HeaderMap.Get("X-Status-Reason"),
		"Malformed JSON request body",
		"The reason header is not set properly",
	)
}

func TestCreateHandlerWithoutRequestBody(t *testing.T) {
	handler := NewCreateHandler(&SessionDatastore{})
	r, err := recordedResponse(handler, "PUT", "")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 400, r.Code, "The status code should be 400 without a request body")
	assert.Equal(
		t,
		r.HeaderMap.Get("X-Status-Reason"),
		"Missing request body",
		"The reason header is not set properly",
	)
}

func TestCreateHandlerWithEmptyTodoTitle(t *testing.T) {
	handler := NewCreateHandler(&SessionDatastore{})
	r, err := recordedResponse(handler, "PUT", `{"title": "", "description": "A short description"}`)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 400, r.Code, "The status code should be 400 with an empty todo title")
	assert.Equal(
		t,
		r.HeaderMap.Get("X-Status-Reason"),
		"Validation failed; See body for reasons",
		"The reason header is not set properly",
	)
	assert.Equal(
		t,
		r.HeaderMap.Get("Content-Type"),
		"application/json; charset=utf-8",
		"The reason header is not set properly",
	)
	validationErrorResponse, err := validationErrorResponseFromBody(r.Body)
	if err != nil {
		t.Error(err)
	}
	if len(validationErrorResponse.Errors) != 1 {
		t.Error("Expected the response to have 1 error, but had", len(validationErrorResponse.Errors))
		t.Error("Expected the response to have the 'title required error' with the error code 1474574120")
	} else if validationErrorResponse.Errors[0].Code != 1474574120 {
		t.Error("Expected the response to have the 'title required error' with the error code 1474574120")
	}
	// TODO assert it hasn't been created
}

func TestCreateHandlerWithValidTodo(t *testing.T) {
	datastore := &SessionDatastore{}
	handler := NewCreateHandler(datastore)
	r, err := recordedResponse(handler, "PUT", `{"title": "My Test Todo", "description": "A short description"}`)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 201, r.Code, "The status code should be 201 with a valid todo")
	assert.Len(t, datastore.Index(), 1)
	todo := datastore.Index()[0]
	assert.Equal(t, "My Test Todo", todo.Title)
	assert.Equal(t, "A short description", todo.Description)
}

func validationErrorResponseFromBody(b *bytes.Buffer) (ValidationErrorResponse, error) {
	var validationErrorResponse ValidationErrorResponse
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return validationErrorResponse, err
	}
	err = json.Unmarshal(body, &validationErrorResponse)
	return validationErrorResponse, err
}

func recordedResponse(handler http.Handler, method string, bodyString string) (*httptest.ResponseRecorder, error) {
	resp := httptest.NewRecorder()

	req, err := http.NewRequest("PUT", "", bytes.NewBuffer([]byte(bodyString)))
	if err != nil {
		return nil, err
	}

	handler.ServeHTTP(resp, req)

	return resp, nil
}
