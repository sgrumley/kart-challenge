package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func Respond(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("unable to encode response data with error: ", err)
	}
}

// RespondNoContent acknowleges success but has nothing to return
func RespondNoContent(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func jsonError(w http.ResponseWriter, status int, code, msg string) {
	w.WriteHeader(status)
	w.Header().Add("content-Type", "application/json")

	err := ErrorResponse{
		Error: &ErrorPayload{
			Code:    code,
			Message: msg,
		},
	}
	encErr := json.NewEncoder(w).Encode(err)
	if encErr != nil {
		log.Println(err)
	}
}

func RespondJSONError(w http.ResponseWriter, err error) {
	var apiErr APIError
	if !errors.As(err, &apiErr) {
		apiErr = Err500Default
	}

	status, code, msg := apiErr.GetData()
	jsonError(w, status, code, msg)
}
