package helper

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Cinemaniac/slg"
	"net/http"
)

func LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	slg.Logger.Error(err.Error(), "method", method, "uri", uri)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	env := Envelope{"error": message}

	if err := WriteJSON(w, status, env, nil); err != nil {
		LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, r, http.StatusNotFound, message)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	env := Envelope{"errors": errors}
	if err := WriteJSON(w, http.StatusUnprocessableEntity, env, nil); err != nil {
		LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	ErrorResponse(w, r, http.StatusConflict, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}
