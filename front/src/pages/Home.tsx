import { Link } from 'react-router-dom'

function Home() {
  return (
    <div>
      <h1>Home</h1>
      <div>
        <ul>
          <li>
            <Link to="/todo">Todos</Link>
          </li>
          <li>
            <Link to="/signin">Sign in</Link>
          </li>
        </ul>
      </div>
    </div>
  )
}

export default Home
