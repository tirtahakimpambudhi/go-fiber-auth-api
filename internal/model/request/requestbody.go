package request

// Struct representing a user creation or update request.
type User struct {
	Username string `json:"username" form:"username" validate:"required,min=5"`   // Required username with minimum length of 5
	Email    string `json:"email" form:"email" validate:"required,email,max=254"` // Required email with max length of 254
	Password string `json:"password" form:"password" validate:"required,min=8"`   // Required password with minimum length of 8
}

// Struct representing a user edit request with optional fields.
type UserEdit struct {
	Username string `json:"username" form:"username" validate:"min=5"`   // Required username with minimum length of 5
	Email    string `json:"email" form:"email" validate:"email,max=254"` // Optional email with max length of 254
	Password string `json:"password" form:"password" validate:"min=8"`   // Optional password with minimum length of 8
}

// Struct for authentication requests.
type Auth struct {
	Email    string `json:"email" form:"email" validate:"required,email,max=254"` // Required email with max length of 254
	Password string `json:"password" form:"password" validate:"required,min=8"`   // Required password with minimum length of 8
}

// Struct for reset password request
type ResetPassword struct {
	Password string `json:"password" form:"password" validate:"required,min=8"`
	Confirm  string `json:"confirm_password" form:"confirm_password" validate:"required,min=8,eqfield=Password"`
}

// Struct for forgot password request
type ForgotPassword struct {
	Email string `json:"email" form:"email" validate:"required,email,max=254"`
}

// Struct for update role request
type UpdateRole struct {
	Email    string `json:"email" form:"email" validate:"required,email,max=254"`
	RoleName string `json:"role_name" form:"role_name" validate:"required,min=1"`
}
