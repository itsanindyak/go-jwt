package helpers

type messageType string

var (
	Error   messageType = "error"
	Success messageType = "success"
)

type Response struct {
	StatusCode  int         `json:"status_code" validate:"required"`
	MessageType messageType `json:"message_type" validate:"required"`
	Message     string      `json:"message" validate:"required"`
	Data        interface{} `json:"data,omitempty" `
}

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

type OTPReq struct {
	OTP string `json:"otp" validate:"required"`
}

type NameUpdate struct {
	FirstName string `json:"first_name" validate:"required,min=3,max=20"`
	LastName  string `json:"last_name" validate:"required,min=3,max=20"`
}
type EmailUpdate struct {
	Email string `json:"email" validate:"required,email"`
}
type TypeUpdate struct {
	UserType string `json:"user_type" validate:"required,oneof=ADMIN MODERATOR USER"`
}

type UpdateData struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=3,max=20"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=3,max=20"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	UserType  string `json:"user_type,omitempty" validate:"omitempty,oneof=ADMIN MODERATOR USER"`
}
