package domain

import "time"

type Task struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type TaskRepository interface {
    Create(task *Task) error
    Update(task *Task) error
    Delete(id, userID int64) error
    GetByID(id, userID int64) (*Task, error)
    GetAllByUserID(userID int64) ([]Task, error)
}

type TaskUsecase interface {
    Create(userID int64, title, description string) error
    Update(id, userID int64, title, description string, done bool) error
    Delete(id, userID int64) error
    GetByID(id, userID int64) (*Task, error)
    GetAllByUserID(userID int64) ([]Task, error)
}