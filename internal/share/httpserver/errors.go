package httpserver

import "net/http"

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}
