package service

import (
	"context"
	"database/sql"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/models"
)

type TaskService struct {
	DB *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{DB: db}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (user_id, title, description, category_id, completed, due_date)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id, created_at`
	err := s.DB.QueryRowContext(ctx, query, task.UserID, task.Title, task.Description, task.CategoryID, task.Completed, task.DueDate).Scan(&task.ID, &task.CreatedAt)
	return err
}

func (s *TaskService) GetTasksByUser(ctx context.Context, userID int) ([]*models.Task, error) {
	query := `SELECT id, user_id, title, description, category_id, completed, created_at, due_date
	FROM tasks WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.CategoryID, &t.Completed, &t.CreatedAt, &t.DueDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, category_id = $3,completed = $4, due_date = $5 WHERE id = $6 AND user_id = $7`
	_, err := s.DB.ExecContext(ctx, query, task.Title, task.Description, task.CategoryID, task.Completed, task.DueDate, task.ID, task.UserID)
	return err
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID int, userID int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	_, err := s.DB.ExecContext(ctx, query, taskID, userID)
	return err
}

func (s *TaskService) MarkTaskCompletion(ctx context.Context, taskID int, userID int, completed bool) error {
	query := `UPDATE tasks SET completed = $1 WHERE id = $2 AND user_id = $3`
	_, err := s.DB.ExecContext(ctx, query, completed, taskID, userID)
	return err
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskID int, userID int) (*models.Task, error) {
	query := `SELECT id, user_id, title, description, category_id, completed, created_at, due_date FROM tasks WHERE id = $1 AND user_id = $2`
	row := s.DB.QueryRowContext(ctx, query, taskID, userID)

	var task models.Task
	err := row.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.CategoryID, &task.Completed, &task.CreatedAt, &task.DueDate)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
