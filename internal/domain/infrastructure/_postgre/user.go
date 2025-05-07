package _postgre

import (
	"database/sql"
	"github.com/kaffein/goffy/internal/domain/entities"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Create(user *entities.User) (*entities.User, error) {
	query := `INSERT INTO users (name, email, password, role, gender) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var newID string
	err := repo.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role, user.Gender).Scan(&newID)
	if err != nil {
		return nil, err
	}

	user.ID = newID
	return user, nil
}
