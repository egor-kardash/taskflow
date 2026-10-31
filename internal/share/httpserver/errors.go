package httpserver

import "net/http"

type errorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteErrorFields(w, status, code, message, nil)
}

func WriteErrorFields(w http.ResponseWriter, status int, code, message string, fields map[string]string) {
	WriteJSON(w, status, errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}
