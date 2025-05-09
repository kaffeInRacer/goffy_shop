package db_postgre

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kaffein/goffy/internal/app/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

var _ user.UserRepository = &UserRepository{}

func (r *UserRepository) Create(ctx context.Context, d *user.UserEntities) (*user.UserEntities, error) {
	query := `INSERT INTO users (slug, username, password, email, role, gender, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	          RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, d.Slug, d.Username, d.Password, d.Email, d.Role, d.Gender).
		Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *UserRepository) ShowAll(ctx context.Context) ([]*user.UserEntities, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, slug, username, email, role, gender, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.UserEntities
	for rows.Next() {
		u := &user.UserEntities{}
		if err := rows.Scan(&u.ID, &u.Slug, &u.Username, &u.Email, &u.Role, &u.Gender, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Get(ctx context.Context, slug string) (*user.UserEntities, error) {
	u := &user.UserEntities{}
	query := `SELECT id, slug, username, email, role, gender, created_at, updated_at FROM users WHERE slug = $1`
	err := r.db.QueryRowContext(ctx, query, slug).
		Scan(&u.ID, &u.Slug, &u.Username, &u.Email, &u.Role, &u.Gender, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) Update(ctx context.Context, d *user.UserEntities) (*user.UserEntities, error) {
	query := `UPDATE users SET username = $1, email = $2, role = $3, gender = $4, updated_at = NOW() WHERE id = $5 RETURNING updated_at`
	err := r.db.QueryRowContext(ctx, query, d.Username, d.Email, d.Role, d.Gender, d.ID).
		Scan(&d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}
