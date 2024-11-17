package handler

import (
	"net/http"
	"todo-app/internal/delivery/http/middleware"
	"todo-app/internal/delivery/http/request"
	"todo-app/internal/delivery/http/response"
	"todo-app/internal/domain"
)

type TaskHandler struct {
    taskUsecase domain.TaskUsecase
}

func NewTaskHandler(taskUsecase domain.TaskUsecase) *TaskHandler {
    return &TaskHandler{
        taskUsecase: taskUsecase,
    }
}

type createTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
}

type updateTaskRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        response.Error(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    var req createTaskRequest
    if err := request.ParseJSON(r, &req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    err := h.taskUsecase.Create(claims.UserID, req.Title, req.Description)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(w, http.StatusCreated, "Task created successfully", nil)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        response.Error(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    taskID, err := request.GetIDParam(r)
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    var req updateTaskRequest
    if err := request.ParseJSON(r, &req); err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    err = h.taskUsecase.Update(taskID, claims.UserID, req.Title, req.Description, req.Done)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(w, http.StatusOK, "Task updated successfully", nil)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        response.Error(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    taskID, err := request.GetIDParam(r)
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    err = h.taskUsecase.Delete(taskID, claims.UserID)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(w, http.StatusOK, "Task deleted successfully", nil)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        response.Error(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    taskID, err := request.GetIDParam(r)
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid task ID")
        return
    }

    task, err := h.taskUsecase.GetByID(taskID, claims.UserID)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(w, http.StatusOK, "Task retrieved successfully", task)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    claims, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        response.Error(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    tasks, err := h.taskUsecase.GetAllByUserID(claims.UserID)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(w, http.StatusOK, "Tasks retrieved successfully", tasks)
}