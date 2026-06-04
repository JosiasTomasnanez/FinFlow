import { useState } from 'react'

const apiFetch = async (path, options = {}) => {
  const response = await fetch(path, options)
  const data = await response.json().catch(() => null)
  if (!response.ok) {
    throw new Error(data?.error || response.statusText)
  }
  return data
}

function Login({ onLogin, featureEnabled }) {
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('password')
  const [message, setMessage] = useState('')

  const handleLogin = async (event) => {
    event.preventDefault()
    setMessage('Iniciando sesión...')

    try {
      const result = await apiFetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      })
      setMessage(`Bienvenido ${result.username}`)
      onLogin(result)
    } catch (error) {
      setMessage(`Error: ${error.message}`)
    }
  }

  if (!featureEnabled) {
    return null
  }

  return (
    <section>
      <h2>Login</h2>
      <form onSubmit={handleLogin}>
        <label>
          Usuario
          <input value={username} onChange={(e) => setUsername(e.target.value)} required />
        </label>
        <label>
          Contraseña
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
        </label>
        <button type="submit">Iniciar sesión</button>
      </form>
      <pre>{message}</pre>
    </section>
  )
}

export default Login
