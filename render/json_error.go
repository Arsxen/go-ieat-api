package render

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorJSON(w http.ResponseWriter, code int, msg string) {
	err := jsonError{Code: code, Message: msg}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
