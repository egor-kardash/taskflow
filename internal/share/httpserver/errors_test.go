package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) errorResponse {
	t.Helper()

	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type = %q, want JSON", ct)
	}

	var got errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body %q: %v", rec.Body.String(), err)
	}
	if got.Error.Code == "" {
		t.Fatalf("empty error code in %q", rec.Body.String())
	}
	return got
}

func TestNotFoundIsJSON(t *testing.T) {
	router := NewRouter(newTestLogger(), stubPinger{})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if got := decodeError(t, rec); got.Error.Code != "not_found" {
		t.Fatalf("code = %q, want not_found", got.Error.Code)
	}
}

func TestMethodNotAllowedIsJSON(t *testing.T) {
	router := NewRouter(newTestLogger(), stubPinger{})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if got := decodeError(t, rec); got.Error.Code != "method_not_allowed" {
		t.Fatalf("code = %q, want method_not_allowed", got.Error.Code)
	}
}

func TestRecovererWritesJSONError(t *testing.T) {
	handler := Recoverer(newTestLogger())(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/panic", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if got := decodeError(t, rec); got.Error.Code != "internal_error" {
		t.Fatalf("code = %q, want internal_error", got.Error.Code)
	}
	if strings.Contains(rec.Body.String(), "boom") {
		t.Fatalf("response leaks panic details: %q", rec.Body.String())
	}
}
