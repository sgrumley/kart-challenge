package web

import (
	"fmt"
	"net/http"
	"strings"
)

type APIError interface {
	GetData() (int, string, string)
}

type ErrorResponse struct {
	Error *ErrorPayload `json:"error"`
}

type ErrorPayload struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

type Error struct {
	Err         error
	Status      int
	Code        string
	Description string
}

func NewRequestError(err error, status int, code string, publicMsg string) *Error {
	return &Error{
		Err:         err,
		Status:      status,
		Code:        code,
		Description: publicMsg,
	}
}

func (er Error) Error() string {
	var causeStr string
	if er.Err != nil {
		causeStr = fmt.Sprintf(": %s", er.Err.Error())
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("http status = %v, publicMsg = %s", er.Status, er.Description))
	if er.Code != "" {
		sb.WriteString(fmt.Sprintf(", code = %s", er.Code))
	}
	if causeStr != "" {
		sb.WriteString(causeStr)
	}

	return sb.String()
}

func (er Error) GetData() (int, string, string) {
	return er.Status, er.Code, er.Description
}

// For use when dealing with external clients
// hides any internal information that should not be exposed publicly
var (
	// Err400Default ...
	Err400Default = &Error{
		Status:      http.StatusBadRequest,
		Description: http.StatusText(http.StatusBadRequest),
		Code:        "generic_bad_request",
	}

	// Err404Default ...
	Err404Default = &Error{
		Status:      http.StatusNotFound,
		Description: http.StatusText(http.StatusNotFound),
		Code:        "generic_not_found",
	}

	// Err409Default ...
	Err409Default = &Error{
		Status:      http.StatusConflict,
		Description: http.StatusText(http.StatusConflict),
		Code:        "generic_conflict",
	}

	// Err422Default ...
	Err422Default = &Error{
		Status:      http.StatusUnprocessableEntity,
		Description: "Your request has not been processed, some precondition failed",
		Code:        "generic_unprocessable_entity",
	}

	// Err500Default ...
	Err500Default = &Error{
		Status:      http.StatusInternalServerError,
		Description: http.StatusText(http.StatusInternalServerError),
		Code:        "generic_internal_server_error",
	}
)
