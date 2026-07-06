package profile

import "time"

type UserProfile struct {
	ID        string    `json:"id"`
	Name      *string   `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name  *string `json:"name"`
	Phone *string `json:"phone"`
}
