import { useState } from 'react'
import { BrowserRouter, Route, Link } from 'react-router-dom'

function App() {

  return (
    <BrowserRouter basename="/todo-go">
      <div>
        <h1>
          Hello world!
        </h1>
      </div>
    </BrowserRouter>
  )
}

export default App
