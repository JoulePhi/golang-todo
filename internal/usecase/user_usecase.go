package usecase

import (
	"fmt"
	"todo-app/internal/domain"
	"todo-app/internal/pkg/auth"
	"todo-app/internal/pkg/security"
)

type userUsecase struct {
    userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
    return &userUsecase{
        userRepo: userRepo,
    }
}

func (u *userUsecase) Register(username, password string) error {
    // Check if username already exists
    existingUser, err := u.userRepo.GetByUsername(username)
    if err != nil {
        return fmt.Errorf("error checking username: %v", err)
    }
    if existingUser != nil {
        return fmt.Errorf("username already exists")
    }

    // Hash password
    hashedPassword, err := security.HashPassword(password)
    if err != nil {
        return fmt.Errorf("error hashing password: %v", err)
    }

    // Create user
    user := &domain.User{
        Username: username,
        Password: hashedPassword,
    }

    if err := u.userRepo.Create(user); err != nil {
        return fmt.Errorf("error creating user: %v", err)
    }

    return nil
}

func (u *userUsecase) Login(username, password string) (string, error) {
    // Get user by username
    user, err := u.userRepo.GetByUsername(username)
    if err != nil {
        return "", fmt.Errorf("error getting user: %v", err)
    }
    if user == nil {
        return "", fmt.Errorf("invalid username or password")
    }

    // Verify password
    valid, err := security.VerifyPassword(user.Password, password)
    if err != nil {
        return "", fmt.Errorf("error verifying password: %v", err)
    }
    if !valid {
        return "", fmt.Errorf("invalid username or password")
    }

    // Generate JWT token
    token, err := auth.GenerateToken(user.ID, user.Username)
    if err != nil {
        return "", fmt.Errorf("error generating token: %v", err)
    }

    return token, nil
}