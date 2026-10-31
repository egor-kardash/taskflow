package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/egor-kardash/taskflow/internal/share/httpserver"
)

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid_user_id", "user id must be a positive integer")
		return
	}

	user, err := h.users.GetUser(r.Context(), id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	httpserver.WriteJSON(w, http.StatusOK, toUserResponse(user))
}
