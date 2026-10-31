package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
	"github.com/egor-kardash/taskflow/internal/features/users/service"
)

type fakeService struct {
	created   domain.User
	createErr error
	user      domain.User
	getErr    error
}

func (f *fakeService) CreateUser(context.Context, service.CreateUserInput) (domain.User, error) {
	return f.created, f.createErr
}

func (f *fakeService) GetUser(context.Context, int64) (domain.User, error) {
	return f.user, f.getErr
}

func newTestRouter(svc UsersService) http.Handler {
	r := chi.NewRouter()
	New(svc, slog.New(slog.NewTextHandler(io.Discard, nil))).Register(r)
	return r
}

func TestCreateUser201(t *testing.T) {
	svc := &fakeService{created: domain.User{ID: 7, Version: 1, FullName: "Egor"}}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"full_name":"Egor"}`))
	newTestRouter(svc).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "/users/7" {
		t.Fatalf("location = %q", got)
	}

	var body userResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.ID != 7 {
		t.Fatalf("id = %d", body.ID)
	}
}

func TestCreateUserInvalidJSON(t *testing.T) {
	tests := map[string]string{
		"malformed":     `{"full_name":`,
		"empty":         ``,
		"unknown field": `{"full_name":"Egor","typo":1}`,
		"trailing data": `{"full_name":"Egor"}{"full_name":"X"}`,
		"trailing junk": `{"full_name":"Egor"} x`,
	}

	for name, payload := range tests {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
			newTestRouter(&fakeService{}).ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", rec.Code)
			}
		})
	}
}

func TestCreateUserValidationFields(t *testing.T) {
	svc := &fakeService{createErr: &domain.ValidationError{
		Fields: map[string]string{"full_name": "must not be blank"},
	}}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"full_name":""}`))
	newTestRouter(svc).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}

	var body struct {
		Error struct {
			Code   string            `json:"code"`
			Fields map[string]string `json:"fields"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Error.Code != "validation_failed" || body.Error.Fields["full_name"] == "" {
		t.Fatalf("body = %+v", body)
	}
}

func TestCreateUserConflict(t *testing.T) {
	svc := &fakeService{createErr: domain.ErrPhoneTaken}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"full_name":"Egor"}`))
	newTestRouter(svc).ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
}

func TestCreateUserTooLarge(t *testing.T) {
	payload := `{"full_name":"` + strings.Repeat("a", maxBodyBytes) + `"}`

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	newTestRouter(&fakeService{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestGetUser200(t *testing.T) {
	svc := &fakeService{user: domain.User{ID: 7, Version: 1, FullName: "Egor"}}

	rec := httptest.NewRecorder()
	newTestRouter(svc).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/7", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestGetUserNotFound(t *testing.T) {
	svc := &fakeService{getErr: domain.ErrNotFound}

	rec := httptest.NewRecorder()
	newTestRouter(svc).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/999", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestGetUserInvalidID(t *testing.T) {
	for _, id := range []string{"abc", "0", "-1", "1.5"} {
		t.Run(id, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestRouter(&fakeService{}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/"+id, nil))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d", rec.Code)
			}
		})
	}
}

func TestGetUserInternalErrorHidesDetails(t *testing.T) {
	svc := &fakeService{getErr: errors.New("connection refused")}

	rec := httptest.NewRecorder()
	newTestRouter(svc).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/7", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "connection refused") {
		t.Fatalf("leaked: %s", rec.Body.String())
	}
}
