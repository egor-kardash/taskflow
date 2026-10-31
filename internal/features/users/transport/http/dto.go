package http

import "github.com/egor-kardash/taskflow/internal/features/users/domain"

type createUserRequest struct {
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

type userResponse struct {
	ID          int64   `json:"id"`
	Version     int64   `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func toUserResponse(user domain.User) userResponse {
	return userResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}
