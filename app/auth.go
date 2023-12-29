package app

type Login struct {
	Email    string `json:"email"  valid:"required, email"`
	Password string `json:"password" valid:"required"`
}

type Register struct {
	Username string `json:"username" valid:"required"`
	Email    string `json:"email" valid:"required, email"`
	Password string `json:"password" valid:"required"`
}