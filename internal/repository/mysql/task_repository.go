package repository

import (
	"database/sql"
	"fmt"
	"time"
	"todo-app/internal/domain"
)

type mysqlTaskRepository struct {
    db *sql.DB
}

func NewMysqlTaskRepository(db *sql.DB) domain.TaskRepository {
    return &mysqlTaskRepository{db}
}

func (r *mysqlTaskRepository) Create(task *domain.Task) error {
    query := `
        INSERT INTO tasks (user_id, title, description, done, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    
    now := time.Now()
    result, err := r.db.Exec(query,
        task.UserID,
        task.Title,
        task.Description,
        task.Done,
        now,
        now,
    )
    if err != nil {
        return fmt.Errorf("error creating task: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("error getting last insert id: %v", err)
    }

    task.ID = id
    task.CreatedAt = now
    task.UpdatedAt = now
    return nil
}

func (r *mysqlTaskRepository) Update(task *domain.Task) error {
    query := `
        UPDATE tasks 
        SET title = ?, description = ?, done = ?, updated_at = ?
        WHERE id = ? AND user_id = ?
    `
    
    now := time.Now()
    result, err := r.db.Exec(query,
        task.Title,
        task.Description,
        task.Done,
        now,
        task.ID,
        task.UserID,
    )
    if err != nil {
        return fmt.Errorf("error updating task: %v", err)
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %v", err)
    }
    if affected == 0 {
        return fmt.Errorf("task not found or unauthorized")
    }

    task.UpdatedAt = now
    return nil
}

func (r *mysqlTaskRepository) Delete(id, userID int64) error {
    query := `DELETE FROM tasks WHERE id = ? AND user_id = ?`
    
    result, err := r.db.Exec(query, id, userID)
    if err != nil {
        return fmt.Errorf("error deleting task: %v", err)
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting rows affected: %v", err)
    }
    if affected == 0 {
        return fmt.Errorf("task not found or unauthorized")
    }

    return nil
}

func (r *mysqlTaskRepository) GetByID(id, userID int64) (*domain.Task, error) {
    query := `
        SELECT id, user_id, title, description, done, created_at, updated_at
        FROM tasks
        WHERE id = ? AND user_id = ?
    `

    task := &domain.Task{}
    err := r.db.QueryRow(query, id, userID).Scan(
        &task.ID,
        &task.UserID,
        &task.Title,
        &task.Description,
        &task.Done,
        &task.CreatedAt,
        &task.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("error getting task: %v", err)
    }

    return task, nil
}

func (r *mysqlTaskRepository) GetAllByUserID(userID int64) ([]domain.Task, error) {
    query := `
        SELECT id, user_id, title, description, done, created_at, updated_at
        FROM tasks
        WHERE user_id = ?
        ORDER BY created_at DESC
    `

    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("error querying tasks: %v", err)
    }
    defer rows.Close()

    var tasks []domain.Task
    for rows.Next() {
        var task domain.Task
        err := rows.Scan(
            &task.ID,
            &task.UserID,
            &task.Title,
            &task.Description,
            &task.Done,
            &task.CreatedAt,
            &task.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning task: %v", err)
        }
        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating tasks: %v", err)
    }

    return tasks, nil
}