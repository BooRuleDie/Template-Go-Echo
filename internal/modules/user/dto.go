package user

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=20,alpha"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,phone"`
	// Age is an optional field. If not provided, it remains nil and is skipped during validation due to the "omitempty" tag.
	// When specified, its value must be between 15 and 90, inclusive, as enforced by the "gte" and "lte" rules.
	// This design allows 0 (zero value for int) to be accepted as a legitimate value if explicitly set,
	// rather than being treated as "not provided".
	Age *int `json:"age" validate:"omitempty,gte=15,lte=90"`
}
