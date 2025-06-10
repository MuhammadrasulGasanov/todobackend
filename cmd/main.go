package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MuhammadrasulGasanov/go-tasks/internal/config"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/handler"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/middleware"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/repository"
	"github.com/MuhammadrasulGasanov/go-tasks/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	db, err := sql.Open("postgres", cfg.GetDBConnString())
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}

	//Dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)
	taskService := service.NewTaskService(db)
	taskHandler := handler.NewTaskHandler(taskService)
	categoryService := service.NewCategoryService(db)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	//Router
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "DB not reachable", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "OK")
	})

	r.Post("/register", userHandler.Register)
	r.Post("/login", authHandler.Login)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(cfg.JWTSecret))

		r.Get("/me", handler.ProtectedEndpoint)

		// Task routes
		r.Post("/tasks", taskHandler.CreateTask)
		r.Get("/tasks", taskHandler.GetTasks)
		r.Put("/tasks/{id}", taskHandler.UpdateTask)
		r.Delete("/tasks/{id}", taskHandler.DeleteTask)
		r.Patch("/tasks/{id}", taskHandler.MarkTaskCompletion)
		// Category routes
		r.Post("/categories", categoryHandler.CreateCategory)
		r.Get("/categories", categoryHandler.GetCategories)
		r.Get("/categories/{id}", categoryHandler.GetCategoryById)
		r.Delete("/categories/{id}", categoryHandler.DeleteCategory)
	})

	log.Printf("Server is running on %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(""+cfg.ServerPort, r))
}
