package dto

import "github.com/kaffein/goffy/internal/app/user"

type AuthRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
}

func (a *AuthRegister) ParsingData() *user.UserEntities {
	return &user.UserEntities{
		Username: a.Username,
		Password: a.Password,
		Email:    a.Email,
		Gender:   a.Gender,
	}
}

type AuthLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthLogin) ParsingData() *user.UserEntities {
	return &user.UserEntities{
		Email:    a.Email,
		Password: a.Password,
	}
}
