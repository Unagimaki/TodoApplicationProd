package todo

import "errors"

var (
	ErrInvalidID     = errors.New("invalid ID")
	ErrTitleRequired = errors.New("title is required")
	ErrTodoNotFound  = errors.New("todo not found")
)
