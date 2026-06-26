package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateTodo(ctx context.Context, title string) (Todo, error) {
	query := `
		INSERT INTO todos (title)
		VALUES ($1)
		RETURNING id, title, completed
	`

	var todo Todo

	err := r.db.QueryRowContext(ctx, query, title).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
	)
	if err != nil {
		return Todo{}, fmt.Errorf("repository create todo: %w", err)
	}
	return todo, nil

}
func (r *Repository) GetAllTodos(ctx context.Context) ([]Todo, error) {
	todos := []Todo{}

	query := `
		SELECT id, title, completed
		FROM todos
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return todos, nil
}
func (r *Repository) DeleteTodo(ctx context.Context, id int) error {
	query := `
		DELETE FROM todos
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository delete todo exec: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository delete todo rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("repository delete todo id: %d: %w", id, ErrTodoNotFound)
	}

	return nil
}
func (r *Repository) UpdateTodoTitle(ctx context.Context, id int, title string) (Todo, error) {
	query := `
		UPDATE todos
		SET title = $1
		WHERE id = $2
		RETURNING id, title, completed
	`

	var todo Todo

	err := r.db.QueryRowContext(ctx, query, title, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, fmt.Errorf("repository update todo title id %d: %w", id, ErrTodoNotFound)
		}
		return Todo{}, fmt.Errorf("repository update todo title: %w", err)
	}
	return todo, nil
}
func (r *Repository) ToggleTodo(ctx context.Context, id int) (Todo, error) {
	query := `
		UPDATE todos
		SET completed = NOT completed
		WHERE id = $1
		RETURNING id, title, completed
	`
	var todo Todo
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, fmt.Errorf("repository toggle todo id %d: %w", id, ErrTodoNotFound)
		}

		return Todo{}, fmt.Errorf("repository update todo status: %w", err)
	}
	return todo, nil
}
