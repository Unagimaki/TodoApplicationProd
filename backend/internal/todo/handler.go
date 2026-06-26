package todo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}
type CreateTodoRequest struct {
	Title string `json:"title"`
}
type UpdateTitleTodoRequest struct {
	Title string `json:"title"`
}
type GetAllTodosResponse struct {
	Todos []Todo `json:"todos"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, context.Canceled):
		log.Println("request canseled by client")
	case errors.Is(err, context.DeadlineExceeded):
		log.Println("request deadline exceeded")
		writeJSONError(w, "request timeout", http.StatusGatewayTimeout)
	case errors.Is(err, ErrInvalidID):
		writeJSONError(w, "invalid id", http.StatusBadRequest)
	case errors.Is(err, ErrTitleRequired):
		writeJSONError(w, "title is required", http.StatusBadRequest)
	case errors.Is(err, ErrTodoNotFound):
		writeJSONError(w, "todo not found", http.StatusNotFound)
	default:
		log.Println(err)
		writeJSONError(w, "internal server error", http.StatusInternalServerError)
	}
}
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Error: message,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to marshal error response", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(responseJSON)
}
func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.service.CreateTodo(ctx, req.Title)
	if err != nil {
		handleError(w, err)
		return
	}

	response, err := json.Marshal(todo)
	if err != nil {
		writeJSONError(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
func (h *Handler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	todos, err := h.service.GetAllTodos(ctx)

	if err != nil {
		handleError(w, err)
		return
	}

	response := GetAllTodosResponse{
		Todos: todos,
	}

	responseJSON, err := json.Marshal(response)
	fmt.Println(responseJSON)
	if err != nil {
		writeJSONError(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		writeJSONError(w, "invalid todo id", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteTodo(ctx, id)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) UpdateTodoTitle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		writeJSONError(w, "invalid todo id", http.StatusBadRequest)
		return
	}
	var req UpdateTitleTodoRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	todo, err := h.service.UpdateTodoTitle(ctx, id, req.Title)
	if err != nil {
		handleError(w, err)
		return
	}
	response, err := json.Marshal(todo)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func (h *Handler) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		writeJSONError(w, "invalid todo id", http.StatusBadRequest)
		return
	}
	todo, err := h.service.ToggleTodo(ctx, id)
	if err != nil {
		handleError(w, err)
		return
	}
	response, err := json.Marshal(todo)
	if err != nil {
		log.Println(err)
		writeJSONError(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
