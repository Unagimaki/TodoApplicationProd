package todo

import (
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

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidID):
		http.Error(w, "invalid id", http.StatusBadRequest)
	case errors.Is(err, ErrTitleRequired):
		http.Error(w, "title is required", http.StatusBadRequest)
	case errors.Is(err, ErrTodoNotFound):
		http.Error(w, "todo not found", http.StatusNotFound)
	default:
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.service.CreateTodo(req.Title)
	if err != nil {
		handleError(w, err)
		return
	}

	response, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *Handler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	todos, err := h.service.GetAllTodos()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetAllTodosResponse{
		Todos: todos,
	}

	responseJSON, err := json.Marshal(response)
	fmt.Println(responseJSON)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteTodo(id)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTodoTitle(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}
	var req UpdateTitleTodoRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	todo, err := h.service.UpdateTodoTitle(id, req.Title)
	if err != nil {
		handleError(w, err)
		return
	}
	response, err := json.Marshal(todo)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *Handler) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalit todo id", http.StatusBadRequest)
		return
	}
	todo, err := h.service.ToggleTodo(id)
	if err != nil {
		handleError(w, err)
		return
	}
	response, err := json.Marshal(todo)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
