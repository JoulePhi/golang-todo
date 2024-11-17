package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, code int, response Response) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, code int, message string) {
    response := Response{
        Status: "error",
        Error:  message,
    }
    JSON(w, code, response)
}

func Success(w http.ResponseWriter, code int, message string, data interface{}) {
    response := Response{
        Status:  "success",
        Message: message,
        Data:    data,
    }
    JSON(w, code, response)
}