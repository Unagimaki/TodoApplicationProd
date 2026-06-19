import { useState } from 'react';

const TodoItem = ({ todo, index, isLoading, toggleTodo, removeTodo, updateTodoTitle }) => {
  const [isEditing, setEditing] = useState(false);
  const [title, setTitle] = useState(todo.title);

  const startEditing = () => {
    setTitle(todo.title);
    setEditing(true);
  };

  const cancelEditing = () => {
    setTitle(todo.title);
    setEditing(false);
  };

  const saveTitle = async (event) => {
    event.preventDefault();
    const isSaved = await updateTodoTitle(todo.id, title);
    if (isSaved) {
      setEditing(false);
    }
  };

  return (
    <li className={`todo-item ${todo.completed ? 'completed' : ''}`}>
      <div className="todo-index">{index}</div>

      {isEditing ? (
        <form className="edit-form" onSubmit={saveTitle}>
          <input
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            disabled={isLoading}
            autoFocus
          />
          <button type="submit" disabled={isLoading}>Сохранить</button>
          <button type="button" className="ghost-button" onClick={cancelEditing} disabled={isLoading}>
            Отмена
          </button>
        </form>
      ) : (
        <>
          <button
            type="button"
            className="status-button"
            onClick={() => toggleTodo(todo.id)}
            disabled={isLoading}
            aria-label={todo.completed ? 'Вернуть задачу в работу' : 'Отметить задачу выполненной'}
            title={todo.completed ? 'Вернуть в работу' : 'Отметить выполненной'}
          >
            {todo.completed ? '✓' : ''}
          </button>
          <div className="todo-title">{todo.title}</div>
          <div className="todo-actions">
            <button type="button" className="ghost-button" onClick={startEditing} disabled={isLoading}>
              Изменить
            </button>
            <button type="button" className="danger-button" onClick={() => removeTodo(todo.id)} disabled={isLoading}>
              Удалить
            </button>
          </div>
        </>
      )}
    </li>
  );
};

export default TodoItem;
