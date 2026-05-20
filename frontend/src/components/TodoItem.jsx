const TodoItem = ({todo, index, toggleTodo, removeTodo}) => {
    return(
        <div style={{display: 'flex', alignItems: 'center', gap: '10px', cursor: 'pointer'}} >
            <div>{index}.</div>
            <h1 onClick={() => toggleTodo(todo.id)} style={{textDecoration: todo.completed ? 'line-through' : '', fontSize: '20px'}} >{todo.title}</h1>
            <button onClick={() => removeTodo(todo.id)}>Удалить</button>
        </div>
    )
}
export default TodoItem