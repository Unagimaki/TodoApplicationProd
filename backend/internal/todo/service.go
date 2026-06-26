package todo

import (
	"context"
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

func (s *Service) CreateTodo(ctx context.Context, title string) (Todo, error) {
	if title == "" {
		return Todo{}, errors.New("title is required")
	}
	todo, err := s.repo.CreateTodo(ctx, title)
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (s *Service) GetAllTodos(ctx context.Context) ([]Todo, error) {
	todos, err := s.repo.GetAllTodos(ctx)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *Service) DeleteTodo(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	err := s.repo.DeleteTodo(ctx, id)
	if err != nil {
		return fmt.Errorf("service error: %w", err)
	}
	return nil
}

func (s *Service) UpdateTodoTitle(ctx context.Context, id int, title string) (Todo, error) {
	if id <= 0 {
		return Todo{}, ErrInvalidID
	}
	if title == "" {
		return Todo{}, ErrTitleRequired
	}
	todo, err := s.repo.UpdateTodoTitle(ctx, id, title)
	if err != nil {
		return Todo{}, fmt.Errorf("service update todo title: %w", err)
	}
	return todo, nil
}
func (s *Service) ToggleTodo(ctx context.Context, id int) (Todo, error) {
	if id <= 0 {
		return Todo{}, ErrInvalidID
	}
	todo, err := s.repo.ToggleTodo(ctx, id)
	if err != nil {
		return Todo{}, fmt.Errorf("service update todo status %w", err)
	}
	return todo, nil
}
