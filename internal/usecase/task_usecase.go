package usecase

import (
	"fmt"
	"todo-app/internal/domain"
)

type taskUsecase struct {
    taskRepo domain.TaskRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository) domain.TaskUsecase {
    return &taskUsecase{
        taskRepo: taskRepo,
    }
}

func (u *taskUsecase) Create(userID int64, title, description string) error {
    if title == "" {
        return fmt.Errorf("title is required")
    }

    task := &domain.Task{
        UserID:      userID,
        Title:       title,
        Description: description,
        Done:        false,
    }

    if err := u.taskRepo.Create(task); err != nil {
        return fmt.Errorf("error creating task: %v", err)
    }

    return nil
}

func (u *taskUsecase) Update(id, userID int64, title, description string, done bool) error {
    // Check if task exists and belongs to user
    existingTask, err := u.taskRepo.GetByID(id, userID)
    if err != nil {
        return fmt.Errorf("error getting task: %v", err)
    }
    if existingTask == nil {
        return fmt.Errorf("task not found or unauthorized")
    }

    if title == "" {
        title = existingTask.Title
    }

    task := &domain.Task{
        ID:          id,
        UserID:      userID,
        Title:       title,
        Description: description,
        Done:        done,
    }

    if err := u.taskRepo.Update(task); err != nil {
        return fmt.Errorf("error updating task: %v", err)
    }

    return nil
}

func (u *taskUsecase) Delete(id, userID int64) error {
    // Check if task exists and belongs to user
    existingTask, err := u.taskRepo.GetByID(id, userID)
    if err != nil {
        return fmt.Errorf("error getting task: %v", err)
    }
    if existingTask == nil {
        return fmt.Errorf("task not found or unauthorized")
    }

    if err := u.taskRepo.Delete(id, userID); err != nil {
        return fmt.Errorf("error deleting task: %v", err)
    }

    return nil
}

func (u *taskUsecase) GetByID(id, userID int64) (*domain.Task, error) {
    task, err := u.taskRepo.GetByID(id, userID)
    if err != nil {
        return nil, fmt.Errorf("error getting task: %v", err)
    }
    if task == nil {
        return nil, fmt.Errorf("task not found or unauthorized")
    }

    return task, nil
}

func (u *taskUsecase) GetAllByUserID(userID int64) ([]domain.Task, error) {
    tasks, err := u.taskRepo.GetAllByUserID(userID)
    if err != nil {
        return nil, fmt.Errorf("error getting tasks: %v", err)
    }

    return tasks, nil
}