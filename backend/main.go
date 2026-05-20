package main

import (
	"log"
	"net/http"
	"todo-backend/handlers"
)

func main() {
	http.HandleFunc("/todos", handlers.CorsHandler(handlers.TodosHandler))
	http.HandleFunc("/todos/", handlers.CorsHandler(handlers.TodoByIDHandler))

	log.Println("server started on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
