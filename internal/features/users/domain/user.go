package domain

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MaxFullNameLen = 70

var phonePattern = regexp.MustCompile(`^\+[1-9][0-9]{6,14}$`)

type User struct {
	ID          int64
	Version     int64
	FullName    string
	PhoneNumber *string
}

func NewUser(fullName string, phoneNumber *string) (User, error) {
	name := strings.TrimSpace(fullName)
	phone, phoneErr := normalizePhone(phoneNumber)

	fields := make(map[string]string, 2)
	switch {
	case name == "":
		fields["full_name"] = "must not be blank"
	case strings.IndexFunc(name, unicode.IsControl) >= 0:
		fields["full_name"] = "must not contain control characters"
	case utf8.RuneCountInString(name) > MaxFullNameLen:
		fields["full_name"] = fmt.Sprintf("must be at most %d characters", MaxFullNameLen)
	}
	if phoneErr != "" {
		fields["phone_number"] = phoneErr
	}
	if len(fields) > 0 {
		return User{}, &ValidationError{Fields: fields}
	}

	return User{FullName: name, PhoneNumber: phone}, nil
}

func normalizePhone(raw *string) (*string, string) {
	if raw == nil {
		return nil, ""
	}

	phone := strings.TrimSpace(*raw)
	if phone == "" {
		return nil, ""
	}
	if !phonePattern.MatchString(phone) {
		return nil, "must match +[1-9][0-9]{6,14}"
	}

	return &phone, ""
}
