package helpers

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Signup struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=20"`
	LastName  string `json:"last_name" validate:"required,min=3,max=20"`
	Email     string `json:"email" validate:"required,email"`
	UserType  string `json:"user_type" validate:"required,oneof=ADMIN MODERATOR USER"`
	Password  string `json:"password" validate:"required,min=6"`
}
