package user

import "go-echo-template/internal/shared/utils"

type GetUserResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone"`
	Role      string  `json:"role"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=20,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`

	// Phone is an optional field. If not provided, it remains nil and is skipped during validation due to the "omitempty" tag.
	// When specified, its value must match the "phone" rule. This allows an empty string, but the intent is to accept only valid phone numbers when set.
	Phone *string `json:"phone" validate:"omitempty,phone"`
}

// HashPassword hashes the Password field using bcrypt.
func (r *CreateUserRequest) HashPassword() (string, error) {
	return utils.HashPassword(r.Password)
}

type UpdateUserRequest struct {
	ID    int64
	Name  string `json:"name" validate:"required,min=3,max=20,alpha"`
	Email string `json:"email" validate:"required,email"`

	// optional
	Phone *string `json:"phone" validate:"omitempty,phone"`
}
