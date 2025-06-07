package repository

import (
	"database/sql"
	"errors"

	"github.com/MuhammadrasulGasanov/go-tasks/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.DB.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}