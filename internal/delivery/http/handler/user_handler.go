package handler

import (
	"net/http"
	"todo-app/internal/delivery/http/request"
	"todo-app/internal/delivery/http/response"
	"todo-app/internal/domain"
)

type UserHandler struct {
	userUseCase domain.UserUsecase
}


func NewUserHandler(userUseCase domain.UserUsecase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}


type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type loginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request){

	var req registerRequest

	if err := request.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {	
		response.Error(w, http.StatusBadRequest, "Username and Password field are required")
		return
	}

	err := h.userUseCase.Register(req.Username, req.Password)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "User Registered Succesfully", nil)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request){

	var req loginRequest

	if err := request.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	
	if req.Username == "" || req.Password == "" {	
		response.Error(w, http.StatusBadRequest, "Username and Password field are required")
		return
	}

	token, err := h.userUseCase.Login(req.Username, req.Password)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "Login Success", loginResponse{Token: token})
}