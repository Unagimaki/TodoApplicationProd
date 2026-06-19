package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"todo-app/internal/todo"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func corsMiddlewareWithLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		fmt.Println("request:", r.Method, r.URL.Path)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := sql.Open("pgx", "postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	fmt.Println("Connected to postgres")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", HealthHandler)

	repo := todo.NewRepo(db)
	service := todo.NewService(repo)
	handler := todo.NewHandler(service)

	mux.HandleFunc("GET /todos", handler.GetAllTodos)
	mux.HandleFunc("POST /todos", handler.CreateTodo)
	mux.HandleFunc("DELETE /todos/{id}", handler.DeleteTodo)
	mux.HandleFunc("PATCH /todos/{id}", handler.UpdateTodoTitle)
	mux.HandleFunc("PATCH /todos/{id}/toggle", handler.ToggleTodo)

	fmt.Println("server started on :8080")
	err = http.ListenAndServe(":8080", corsMiddlewareWithLog(mux))

	if err != nil {
		fmt.Println(err)
	}
}
