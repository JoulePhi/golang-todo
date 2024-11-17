package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"todo-app/internal/delivery/http/handler"
	"todo-app/internal/delivery/http/middleware"
	"todo-app/internal/pkg/mysql"
	repository "todo-app/internal/repository/mysql"
	"todo-app/internal/usecase"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func main() {
	config := parseConfig()

	db, err := mysql.NewConnection(
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName)
	
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	defer db.Close()

	userRepo := repository.NewMysqlUserRepository(db)
	taskRepo := repository.NewMysqlTaskRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo)
	taskUseCase := usecase.NewTaskUsecase(taskRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	taskHandler := handler.NewTaskHandler(taskUseCase)

	authMiddleware := middleware.NewAuthMiddleware()

	router := http.NewServeMux()

	router.HandleFunc("POST /api/register", middleware.Chain(
		userHandler.Register,
		middleware.CORS,
		middleware.Logger))
	
	router.HandleFunc("POST /api/login", middleware.Chain(
		userHandler.Login,
		middleware.CORS,
		middleware.Logger))

	router.HandleFunc("POST /api/tasks", middleware.Chain(
        taskHandler.CreateTask,
        authMiddleware.Authenticate,
        middleware.Logger,
        middleware.CORS,
    ))

	router.HandleFunc("/api/tasks/", middleware.Chain(
        func(w http.ResponseWriter, r *http.Request) {
            switch r.Method {
            case http.MethodGet:
                if r.URL.Path == "/api/tasks/" {
                    taskHandler.GetAllTasks(w, r)
                } else {
                    taskHandler.GetTask(w, r)
                }
            case http.MethodPut:
                taskHandler.UpdateTask(w, r)
            case http.MethodDelete:
                taskHandler.DeleteTask(w, r)
            default:
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            }
        },
        authMiddleware.Authenticate,
        middleware.Logger,
        middleware.CORS,
    ))

	serverAddr := fmt.Sprintf(":%s", config.ServerPort)
    log.Printf("Server starting on %s", serverAddr)
    log.Fatal(http.ListenAndServe(serverAddr, router))
}

func parseConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.DBHost, "db-host", "localhost", "Database host")
	flag.StringVar(&config.DBPort, "db-port", "3306", "Database port")
	flag.StringVar(&config.DBUser, "db-user", "root", "Database user")
	flag.StringVar(&config.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&config.DBName, "db-name", "go_todo", "Database name")

	flag.StringVar(&config.ServerPort, "port", "8080", "Server port")

	flag.Parse()
	return config
}