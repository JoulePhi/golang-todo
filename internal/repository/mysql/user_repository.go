package repository

import (
	"database/sql"
	"fmt"
	"time"
	"todo-app/internal/domain"
)

type mysqlUserRepository struct {
    db *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) domain.UserRepository {
    return &mysqlUserRepository{db}
}

func (r *mysqlUserRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (username, password, created_at, updated_at)
        VALUES (?, ?, ?, ?)
    `
    
    now := time.Now()
    result, err := r.db.Exec(query, 
        user.Username,
        user.Password,
        now,
        now,
    )
    if err != nil {
        return fmt.Errorf("error creating user: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("error getting last insert id: %v", err)
    }

    user.ID = id
    user.CreatedAt = now
    user.UpdatedAt = now
    return nil
}

func (r *mysqlUserRepository) GetByUsername(username string) (*domain.User, error) {
    query := `
        SELECT id, username, password, created_at, updated_at
        FROM users
        WHERE username = ?
    `

    user := &domain.User{}
    err := r.db.QueryRow(query, username).Scan(
        &user.ID,
        &user.Username,
        &user.Password,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("error getting user by username: %v", err)
    }

    return user, nil
}

func (r *mysqlUserRepository) GetByID(id int64) (*domain.User, error) {
    query := `
        SELECT id, username, password, created_at, updated_at
        FROM users
        WHERE id = ?
    `

    user := &domain.User{}
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Password,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("error getting user by id: %v", err)
    }

    return user, nil
}