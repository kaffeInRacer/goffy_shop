package entities

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Gender   string `json:"gender"`
}

func (u *User) ShowID() string {
	return u.ID
}

func (u *User) ShowUsername() string {
	return u.Username
}

func (u *User) ShowEmail() string {
	return u.Email
}
