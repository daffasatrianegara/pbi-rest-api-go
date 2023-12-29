package app

type UpdateUser struct {
	Username string `json:"username" valid:"required"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password"`
}

type GetUser struct {
	ID       uint    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}