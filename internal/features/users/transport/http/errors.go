package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
	"github.com/egor-kardash/taskflow/internal/share/httpserver"
)

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	var validationErr *domain.ValidationError

	switch {
	case errors.As(err, &validationErr):
		httpserver.WriteErrorFields(w, http.StatusBadRequest,
			"validation_failed", "request validation failed", validationErr.Fields)
	case errors.Is(err, domain.ErrNotFound):
		httpserver.WriteError(w, http.StatusNotFound, "user_not_found", "user not found")
	case errors.Is(err, domain.ErrPhoneTaken):
		httpserver.WriteError(w, http.StatusConflict, "phone_number_taken", "phone number already in use")
	default:
		h.log.Error("users request failed",
			"request_id", middleware.GetReqID(r.Context()),
			"error", err,
		)
		httpserver.WriteError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}
