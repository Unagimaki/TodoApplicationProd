package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"todo-backend/model"
	"todo-backend/store"
)

func findTodoById(id uint64, todos []model.Todo) (int, error) {
	for i, todo := range todos {
		if todo.ID == id {
			return i, nil
		}
	}
	return 0, errors.New("todo not found")
}

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandleGetTodos(w)
	case http.MethodPost:
		HandleCreateTodo(w, r)
	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}

func TodoByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		HandleToggleTodo(w, r)
	case http.MethodDelete:
		HandleDeleteTodo(w, r)
	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idNum, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/todos/"), 10, 64)
	if err != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	i, err := findTodoById(idNum, store.Todos)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	store.Todos = slices.Delete(store.Todos, i, i+1)
	w.WriteHeader(http.StatusNoContent)
}

func HandleToggleTodo(w http.ResponseWriter, r *http.Request) {
	idNum, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/todos/"), 10, 64)
	if err != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	i, err := findTodoById(idNum, store.Todos)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	AddHeader(w)
	store.Todos[i].Completed = !store.Todos[i].Completed
	err = json.NewEncoder(w).Encode(store.Todos[i])
	if err != nil {
		http.Error(w, "failed to toggle todo", http.StatusInternalServerError)
		return
	}
}

func HandleGetTodos(w http.ResponseWriter) {
	response := model.TodosResponse{
		Todos: store.Todos,
	}
	w.Header().Set("Content-Type", "application/json")
	AddHeader(w)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func HandleCreateTodo(w http.ResponseWriter, request *http.Request) {
	var createRequest model.CreateTodoRequest

	err := json.NewDecoder(request.Body).Decode(&createRequest)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if createRequest.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	newTodo := model.Todo{
		Title:     createRequest.Title,
		Completed: false,
		ID:        store.NextID,
	}
	store.NextID++

	store.Todos = append(store.Todos, newTodo)
	AddHeader(w)
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newTodo)
	if err != nil {
		http.Error(w, "failed to encode respone", http.StatusInternalServerError)
		return
	}
}
