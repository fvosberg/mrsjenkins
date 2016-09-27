package todo

import (
	"encoding/json"
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
	} else {
		logrus.Printf("Decoding todo create request successful - %+v\n", todo)
		if err := validator.Validate(todo); err != nil {
			logrus.Printf("Validation errors for todo - %+v - %+v\n", todo, err)
			w.Header().Set("X-Status-Reason", "Validation failed; See body for reasons")
			w.WriteHeader(400)
			responseBody := validationErrorResponse(w, err)
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

func validationErrorResponse(w http.ResponseWriter, err error) ValidationErrorResponse {
	var response ValidationErrorResponse
	for field, ves := range err.(validator.ErrorMap) {
		for _, ve := range ves {
			switch {
			default:
				// TODO should this return an error?
				logrus.Printf("TODO LOG STATUS CRITICAL: There is no error code definition for this validation error combination: %s, %s", field, ve)
			case field == "Title" && ve == validator.ErrZeroValue:
				validationError := ValidationError{Code: 1474574120, Description: "The field Title is required"}
				response.Errors = append(response.Errors, validationError)
			}
		}
	}
	return response
}
