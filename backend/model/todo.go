package model

type Todo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	ID        uint64 `json:"id"`
}
type TodosResponse struct {
	Todos []Todo `json:"todos"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}
