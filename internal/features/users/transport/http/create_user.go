package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/egor-kardash/taskflow/internal/features/users/service"
	"github.com/egor-kardash/taskflow/internal/share/httpserver"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			httpserver.WriteError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large")
			return
		}
		httpserver.WriteError(w, http.StatusBadRequest, "invalid_json", "request body must be a single JSON object")
		return
	}

	user, err := h.users.CreateUser(r.Context(), service.CreateUserInput{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	w.Header().Set("Location", "/users/"+strconv.FormatInt(user.ID, 10))
	httpserver.WriteJSON(w, http.StatusCreated, toUserResponse(user))
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("unexpected data after JSON object")
	}

	return nil
}
