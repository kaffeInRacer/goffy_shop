package auth

import (
	"strings"
	"time"

	"github.com/kaffein/goffy/internal/interfaces/http/dto"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type AuthEntities struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Gender    string `json:"gender"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (a *AuthEntities) BeforeSave() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hash)
	return nil
}

func (a *AuthEntities) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plainPassword))
	return err == nil
}

func NewUserFromRegister(r *dto.AuthRegister) *AuthEntities {
	now := time.Now().Format(time.RFC3339)
	return &AuthEntities{
		Slug:      generateSlug(r.Username),
		Username:  r.Username,
		Password:  r.Password,
		Email:     r.Email,
		Role:      RoleUser,
		Gender:    r.Gender,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewAdminFromRegister(r *dto.AuthRegister) *AuthEntities {
	return &AuthEntities{
		Slug:     generateSlug(r.Username),
		Username: r.Username,
		Password: r.Password,
		Email:    r.Email,
		Role:     RoleAdmin,
		Gender:   r.Gender,
	}
}

func generateSlug(username string) string {
	return strings.ToLower(strings.ReplaceAll(username, " ", "-"))
}
