package domain

import "time"

type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"` 
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
    Create(user *User) error
    GetByUsername(username string) (*User, error)
    GetByID(id int64) (*User, error)
}

type UserUsecase interface {
    Register(username, password string) error
    Login(username, password string) (string, error) 
}