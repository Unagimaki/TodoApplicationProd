import { useEffect, useState } from 'react';
import TodoItem from './components/TodoItem';

function App() {
  const [todos, setTodos] = useState([])
  const [isLoading, setLoading] = useState(false)
  const URL = "http://localhost:8080"

  useEffect(() => {
    fetchTodo()
  }, [])

  const removeTodo = async (id) => {
    if (isLoading) return
    setLoading(true)
    try {
      const response = await fetch(`${URL}/todos/${id}`, {
        method: "DELETE",
      })      
      
      if (!response.ok) {
        const text = await response.text()
        throw new Error(text || 'Failed to delete todos')
      }
      setTodos(prev => prev.filter(todo => todo.id !== id))
      
    } catch (error) {
      console.log(error)      
    }
    finally {
      setLoading(false)
    }
  }
  const toggleTodo = async (id) => {
    if (isLoading) return
    setLoading(true)
    try {
      const response = await fetch(`${URL}/todos/${id}`, {
        method: "PUT",
      })      
      if (!response.ok) {
        const text = await response.text()
        throw new Error(text || 'Failed to toggle todos')
      }
      const data = await response.json()
      setTodos(prev => prev.map(todo => todo.id === data.id ? data : todo))
      
    } catch (error) {
      console.log(error)      
    }
    finally {
      setLoading(false)
    }
  }

  const fetchTodo = async () => {
    if (isLoading) return
    setLoading(true)
    try {
      const response = await fetch(`${URL}/todos`)  
      console.log(response)       
      if (!response.ok) {
        const text = await response.text()
        throw new Error(text || 'Failed to fetch todos')
      } 
      const data = await response.json()
      setTodos(data.todos)  
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false)
    }
  }

  const createTodo = async () => {
    if (isLoading) return
    setLoading(true)
    const todo = {
      title: prompt('Title todo!')
    }
    try {
      const response = await fetch(`${URL}/todos`, {
        method: 'POST',
        body: JSON.stringify(todo),
      })
      
      if (!response.ok) {
        const text = await response.text()
        throw new Error(text || 'Failed to create todo')
      }
      const data = await response.json()
      console.log(data);
      
      setTodos(prev => [...prev, data])
    } catch (error) {
      console.log(error)      
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="App">

      <button onClick={fetchTodo}>Получить Todo</button>
      <button onClick={createTodo}>Создать Todo</button>
      {  !todos.length && <h2>No todos</h2>}
      {
        todos.length > 0 &&
        <ul>
          {
            todos.map((todo, index) => {
              return <TodoItem key={todo.id} todo={todo} index={++index} toggleTodo={toggleTodo} removeTodo={removeTodo}/>
            })
          }
        </ul>
      }
    </div>
  );
}

export default App;
