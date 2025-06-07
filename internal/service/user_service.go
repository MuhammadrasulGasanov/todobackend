package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/MuhammadrasulGasanov/go-tasks/internal/models"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(username, password string) (*models.User, error) {
	existing, err := s.Repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashed),
	}

	err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
