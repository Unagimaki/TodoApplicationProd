import { useEffect, useRef, useState } from 'react';
import TodoItem from './components/TodoItem';
import './App.css';

function App() {
  const [todos, setTodos] = useState([]);
  const [isLoading, setLoading] = useState(false);
  const [newTitle, setNewTitle] = useState('');
  const [error, setError] = useState('');
  const createTodoControllerRef = useRef(null);
  const isMountedRef = useRef(true);
  const URL = 'http://localhost:8080';

  useEffect(() => {
    fetchTodo();

    return () => {
      isMountedRef.current = false;
      createTodoControllerRef.current?.abort();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const removeTodo = async (id) => {
    if (isLoading) return;
    setLoading(true);
    setError('');

    try {
      const response = await fetch(`${URL}/todos/${id}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || 'Failed to delete todo');
      }

      setTodos(prev => prev.filter(todo => todo.id !== id));
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const toggleTodo = async (id) => {
    if (isLoading) return;
    setLoading(true);
    setError('');

    try {
      const response = await fetch(`${URL}/todos/${id}/toggle`, {
        method: 'PATCH',
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || 'Failed to update todo status');
      }

      const data = await response.json();
      setTodos(prev => prev.map(todo => todo.id === data.id ? data : todo));
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchTodo = async () => {
    if (isLoading) return;
    setLoading(true);
    setError('');

    try {
      const response = await fetch(`${URL}/todos`);

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || 'Failed to fetch todos');
      }

      const data = await response.json();
      setTodos(data.todos || []);
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const createTodo = async (event) => {
    event.preventDefault();
    if (isLoading) return;

    const title = newTitle.trim();
    if (!title) {
      setError('Введите название задачи');
      return;
    }

    setLoading(true);
    setError('');
    createTodoControllerRef.current?.abort();
    const controller = new AbortController();
    createTodoControllerRef.current = controller;

    try {
      const response = await fetch(`${URL}/todos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title }),
        signal: controller.signal,
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || 'Failed to create todo');
      }

      const data = await response.json();
      if (!isMountedRef.current) return;

      setTodos(prev => [...prev, data]);
      setNewTitle('');
    } catch (error) {
      if (error.name === 'AbortError') {
        return;
      }

      setError(error.message);
    } finally {
      if (createTodoControllerRef.current === controller) {
        createTodoControllerRef.current = null;
      }

      if (isMountedRef.current) {
        setLoading(false);
      }
    }
  };

  const updateTodoTitle = async (id, title) => {
    const nextTitle = title.trim();
    if (!nextTitle) {
      setError('Название задачи не может быть пустым');
      return false;
    }
    if (isLoading) return false;

    setLoading(true);
    setError('');

    try {
      const response = await fetch(`${URL}/todos/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: nextTitle }),
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || 'Failed to update todo title');
      }

      const data = await response.json();
      setTodos(prev => prev.map(todo => todo.id === data.id ? data : todo));
      return true;
    } catch (error) {
      setError(error.message);
      return false;
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <main className="todo-shell">
        <section className="todo-header">
          <p className="eyebrow">Todo app</p>
          <h1>Мои задачи</h1>
          <p className="subtitle">Добавляй задачи, редактируй названия и держи список под рукой.</p>
        </section>

        <form className="todo-form" onSubmit={createTodo}>
          <input
            value={newTitle}
            onChange={(event) => setNewTitle(event.target.value)}
            placeholder="Например: купить продукты"
            disabled={isLoading}
          />
          <button type="submit" disabled={isLoading}>
            Создать
          </button>
          <button type="button" className="secondary-button" onClick={fetchTodo} disabled={isLoading}>
            Обновить
          </button>
        </form>

        {error && <div className="error-message">{error}</div>}

        <section className="todo-list-panel">
          <div className="list-title">
            <h2>Список</h2>
            <span>{todos.length} шт.</span>
          </div>

          {!todos.length && <div className="empty-state">Пока задач нет</div>}

          {todos.length > 0 && (
            <ul className="todo-list">
              {todos.map((todo, index) => (
                <TodoItem
                  key={todo.id}
                  todo={todo}
                  index={index + 1}
                  isLoading={isLoading}
                  toggleTodo={toggleTodo}
                  removeTodo={removeTodo}
                  updateTodoTitle={updateTodoTitle}
                />
              ))}
            </ul>
          )}
        </section>
      </main>
    </div>
  );
}

export default App;
