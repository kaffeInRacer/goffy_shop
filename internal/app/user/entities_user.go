package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type UserEntities struct {
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

// BeforeSave prepares the user entity before saving to the database.
// It generates a slug from the username and hashes the plain-text password.
func (u *UserEntities) BeforeSave() error {
	// Generate slug: lowercase username and replace spaces with dashes.
	u.Slug = strings.ToLower(strings.ReplaceAll(u.Username, " ", "-"))

	// Role: set default role if empty
	switch u.Role {
	case "":
		u.Role = RoleUser
	case "admin":
		u.Role = RoleAdmin
	default:
		return errors.New("invalid role")
	}

	// Hash the password using bcrypt with the default cost.
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}
