package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-app/internal/todo"

	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	serverPort := os.Getenv("SERVER_PORT")
	address := ":" + serverPort
	db, err := sql.Open("pgx", databaseURL)

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

	fmt.Println("server started on", address)
	server := &http.Server{
		Addr:    address,
		Handler: corsMiddlewareWithLog(mux),
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", address, err)
		}
	}()
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
	<-shutdownSignal
	fmt.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server shutdown failed: %v\n", err)
	}
	fmt.Println("server stopped")

}
