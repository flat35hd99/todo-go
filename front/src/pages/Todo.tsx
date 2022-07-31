import React, { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import axios from 'axios'
import Dialog from '@mui/material/Dialog'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import DialogTitle from '@mui/material/DialogTitle'
import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'
import { Checkbox } from '@mui/material'

interface Todo {
  id: number
  title: string
  body: string
  done: boolean
}

interface TodoForm {
  title: string
  body: string
  done: boolean
}

const Todo: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([])
  const [open, setOpen] = useState(false)

  const { register, handleSubmit } = useForm<TodoForm>()
  const onSubmit = handleSubmit(async (data) => {
    const { title, body, done } = data
    const response = await axios.post<Todo>('http://localhost:8080/todos', {
      title,
      body,
      done,
    })
    setOpen(false)
    const todo = response.data
    setTodos([...todos, todo])
  })

  useEffect(() => {
    const f = async () => {
      const { data } = await axios.get<{ todos: Todo[] }>(
        `http://localhost:8080/todos`
      )
      setTodos(data.todos)
    }
    f()
  })
  return (
    <div>
      <h1>Todo</h1>
      <ul>
        {todos.map((t) => {
          return <li key={t.id}>{t.title}</li>
        })}
      </ul>
      <Button onClick={() => setOpen(true)}>Add todo</Button>
      <Dialog
        open={open}
        onClose={() => {
          setOpen(false)
        }}
      >
        <DialogTitle>Add todo</DialogTitle>
        <form onSubmit={onSubmit}>
          <DialogContent>
            <TextField {...register('title')} fullWidth />
            <TextField {...register('body')} fullWidth />
            <Checkbox {...register('done')} />
          </DialogContent>
          <DialogActions>
            <Button type="submit">Add</Button>
          </DialogActions>
        </form>
      </Dialog>
    </div>
  )
}

export default Todo
