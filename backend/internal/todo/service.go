package todo

import (
	"errors"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTodo(title string) (Todo, error) {
	if title == "" {
		return Todo{}, errors.New("title is required")
	}
	todo, err := s.repo.CreateTodo(title)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (s *Service) GetAllTodos() ([]Todo, error) {
	todos, err := s.repo.GetAllTodos()
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *Service) DeleteTodo(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	err := s.repo.DeleteTodo(id)
	if err != nil {
		return fmt.Errorf("service error: %w", err)
	}
	return nil
}

func (s *Service) UpdateTodoTitle(id int, title string) (Todo, error) {
	if id <= 0 {
		return Todo{}, ErrInvalidID
	}
	if title == "" {
		return Todo{}, ErrTitleRequired
	}
	todo, err := s.repo.UpdateTodoTitle(id, title)
	if err != nil {
		return Todo{}, fmt.Errorf("service update todo title: %w", err)
	}
	return todo, nil
}
func (s *Service) ToggleTodo(id int) (Todo, error) {
	if id <= 0 {
		return Todo{}, ErrInvalidID
	}
	todo, err := s.repo.ToggleTodo(id)
	if err != nil {
		return Todo{}, fmt.Errorf("service update todo status %w", err)
	}
	return todo, nil
}
