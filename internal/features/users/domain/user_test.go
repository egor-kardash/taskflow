package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewUserTrimsAndNormalizes(t *testing.T) {
	phone := "  +79991234567  "

	user, err := NewUser("  Egor  ", &phone)
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	if user.FullName != "Egor" {
		t.Fatalf("full name = %q", user.FullName)
	}
	if user.PhoneNumber == nil || *user.PhoneNumber != "+79991234567" {
		t.Fatalf("phone = %v", user.PhoneNumber)
	}
}

func TestNewUserEmptyPhoneIsNil(t *testing.T) {
	for _, raw := range []*string{nil, ptr(""), ptr("   ")} {
		user, err := NewUser("Egor", raw)
		if err != nil {
			t.Fatalf("NewUser(%v): %v", raw, err)
		}
		if user.PhoneNumber != nil {
			t.Fatalf("phone = %v, want nil", *user.PhoneNumber)
		}
	}
}

func TestNewUserValidation(t *testing.T) {
	tests := []struct {
		name  string
		user  string
		phone *string
		field string
	}{
		{name: "blank name", user: "   ", field: "full_name"},
		{name: "long name", user: strings.Repeat("я", MaxFullNameLen+1), field: "full_name"},
		{name: "short phone", user: "Egor", phone: ptr("+123"), field: "phone_number"},
		{name: "no plus", user: "Egor", phone: ptr("79991234567"), field: "phone_number"},
		{name: "leading zero", user: "Egor", phone: ptr("+09991234567"), field: "phone_number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser(tt.user, tt.phone)

			var validationErr *ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("err = %v, want ValidationError", err)
			}
			if _, ok := validationErr.Fields[tt.field]; !ok {
				t.Fatalf("fields = %v, want %q", validationErr.Fields, tt.field)
			}
		})
	}
}

func TestNewUserReportsAllFields(t *testing.T) {
	bad := "x"

	_, err := NewUser("", &bad)

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("err = %v, want ValidationError", err)
	}
	if len(validationErr.Fields) != 2 {
		t.Fatalf("fields = %v, want both", validationErr.Fields)
	}
}

func ptr(s string) *string { return &s }
