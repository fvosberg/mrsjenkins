package todo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	validator "gopkg.in/validator.v2"
)

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type createHandler struct {
	datastore Datastore
}

func (c *createHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Printf("Handled request with %s on %s\n", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	todo, err := newFromRequest(r)
	if err == io.EOF {
		logrus.Println("Empty registration request body")
		w.Header().Set("X-Status-Reason", "Missing request body")
		w.WriteHeader(400)
	} else if err != nil {
		logrus.Printf("Decoding todo create request failed - %+v - %+v\n", err, r)
		w.Header().Set("X-Status-Reason", "Malformed JSON request body")
		w.WriteHeader(500)
	} else {
		logrus.Printf("Decoding todo create request successful - %+v\n", todo)
		if err := validator.Validate(todo); err != nil {
			logrus.Printf("Validation errors for todo - %+v - %+v\n", todo, err)
			responseBody, err := validationErrorResponse(w, err)
			if err != nil {
				logrus.Errorf("Todo: %+v - Error: %+v\n", todo, err)
				w.Header().Set("X-Status-Reason", "Validation failed with an unsupported combination of field and validation error. Please contact the admin.")
				w.WriteHeader(500)
			}
			w.Header().Set("X-Status-Reason", "Validation failed; See body for reasons")
			w.WriteHeader(400)
			jsonBody, err := json.Marshal(responseBody)
			if err != nil {
				// should this end up in another response? 500?
				logrus.Printf("Error while marshalling todos validation error response to json. - %+v - %+v\n", err, responseBody)
			}
			w.Write(jsonBody)
		} else {
			logrus.Printf("Created todo successfully - %+v\n", todo)
			c.datastore.Create(&todo)
			w.WriteHeader(201)
			w.Write([]byte("Todo Created"))
		}
	}
}

func NewCreateHandler(datastore Datastore) *createHandler {
	return &createHandler{datastore: datastore}
}

func validationErrorResponse(w http.ResponseWriter, err error) (ValidationErrorResponse, error) {
	var response ValidationErrorResponse
	for field, ves := range err.(validator.ErrorMap) {
		for _, ve := range ves {
			switch {
			default:
				return response, errors.New("Invalid validation error. Field: %s, Error: %s - Please contact your administrator")
			case field == "Title" && ve == validator.ErrZeroValue:
				validationError := ValidationError{Code: 1474574120, Description: "The field Title is required"}
				response.Errors = append(response.Errors, validationError)
			}
		}
	}
	return response, nil
}
