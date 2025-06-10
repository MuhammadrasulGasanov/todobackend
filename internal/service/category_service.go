package service

import (
	"context"
	"database/sql"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/models"
)

type CategoryService struct {
	DB *sql.DB
}

func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{
		DB: db,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id, created_at`
	err := s.DB.QueryRowContext(ctx, query, category.UserID, category.Name).Scan(&category.ID, &category.CreatedAt)
	return err
}

func (s *CategoryService) DeleteCategory(ctx context.Context, categoryId int, userId int) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	_, err := s.DB.ExecContext(ctx, query, categoryId, userId)
	return err
}

func (s *CategoryService) GetCategoriesByUser(ctx context.Context, userID int) ([]*models.Category, error) {
	query := `SELECT id, user_id, name, created_at FROM categories WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (s *CategoryService) GetCategoryById(ctx context.Context, categoryId int, userID int) (*models.Category, error) {
	query := `SELECT id, user_id, name, created_at FROM categories WHERE id = $1 AND user_id = $2`
	row := s.DB.QueryRowContext(ctx, query, categoryId, userID)
	var category models.Category
	err := row.Scan(&category.ID, &category.UserID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
