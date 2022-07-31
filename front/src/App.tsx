import { Routes, Route } from 'react-router-dom'
import Home from './pages/Home'
import Signin from './pages/Signin'
import Todo from './pages/Todo'

function App() {
  return (
    <div>
      <Routes>
        <Route index element={<Home />} />
        <Route path="signin" element={<Signin />} />
        <Route path="todo" element={<Todo />} />
      </Routes>
    </div>
  )
}

export default App
