package store

import (
	"todo-backend/model"
)

var Todos = []model.Todo{
	{Title: "Learn Go", Completed: true, ID: 1},
	{Title: "Build a REST API", Completed: false, ID: 2},
	{Title: "Write tests", Completed: false, ID: 3},
}

var NextID uint64 = 4
