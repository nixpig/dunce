package user

type User struct {
	Id       int    `validate:"required"`
	Username string `validate:"required,min=5,max=16"`
	Email    string `validate:"required,email,max=100"`
	Link     string `validate:"omitempty,url,max=255"`
	Role     string `validate:"required,max=16"`
	Password string `validate:"required"`
}

type UserRequest struct {
	Username string `validate:"required,min=5,max=16"`
	Email    string `validate:"required,email,max=100"`
	Link     string `validate:"omitempty,url,max=255"`
	Role     string `validate:"required,max=16"`
	Password string `validate:"required"`
}

type UserResponse struct {
	Id       int    `validate:"required"`
	Username string `validate:"required,min=5,max=16"`
	Email    string `validate:"required,email,max=100"`
	Link     string `validate:"omitempty,url,max=255"`
	Role     string `validate:"required,max=16"`
}
