package routers

import (
	"boilerplate-sqlc/models"
	"encoding/json"
	"net/http"
)

type resultset struct{}

// Resultset ...
type Resultset interface {
	ResponsWithJSON(w http.ResponseWriter, code int, message interface{})
	ResponsWithError(w http.ResponseWriter, code int, message string)
}

// NewResultset func
func NewResultset() Resultset {
	return &resultset{}
}

// ResponsWithJSON ...
func (*resultset) ResponsWithJSON(w http.ResponseWriter, code int, data interface{}) {

	response, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ResponsWithError ...
func (r *resultset) ResponsWithError(w http.ResponseWriter, code int, message string) {
	var resp models.ErrorResponse
	errorCode := http.StatusInternalServerError
	msg := "Internal Server Error"

	if code != http.StatusInternalServerError {
		errorCode = code
		msg = message
	}

	if message == "context deadline exceeded" {
		errorCode = http.StatusRequestTimeout
		msg = "Request timeout"
	} else if message != "" {
		msg = message
	}
	resp.Status = errorCode
	resp.Error = http.StatusText(errorCode)
	resp.Message = msg
	r.ResponsWithJSON(w, errorCode, resp)
}
