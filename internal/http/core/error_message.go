package core

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorMessage) ShowErrorMessage(w http.ResponseWriter) {
	w.WriteHeader(e.Code)

	bytes, _ := json.Marshal(e)

	http.Error(w, string(bytes), e.Code)
}
