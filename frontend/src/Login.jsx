import { useState } from 'react'
import Alert from '@mui/material/Alert'
import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import Stack from '@mui/material/Stack'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'
import { apiFetch, jsonPost } from './api'

export default function Login({ onLogin, featureEnabled }) {
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('password')
  const [message, setMessage] = useState('')
  const [error, setError] = useState(false)

  const handleLogin = async (event) => {
    event.preventDefault()
    setError(false)
    setMessage('Iniciando sesión...')
    try {
      const result = await apiFetch('/api/login', jsonPost({ username, password }))
      setMessage(`Bienvenido ${result.username}`)
      onLogin(result)
    } catch (err) {
      setError(true)
      setMessage(err.message)
    }
  }

  if (!featureEnabled) {
    return null
  }

  return (
    <Card>
      <CardContent sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Iniciar sesión
        </Typography>
        <Box component="form" onSubmit={handleLogin}>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'flex-end' }}>
            <TextField label="Usuario" value={username} onChange={(e) => setUsername(e.target.value)} required />
            <TextField
              label="Contraseña"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            <Button type="submit" variant="contained">
              Entrar
            </Button>
          </Stack>
        </Box>
        {message ? (
          <Alert severity={error ? 'error' : 'success'} sx={{ mt: 2 }}>
            {message}
          </Alert>
        ) : null}
      </CardContent>
    </Card>
  )
}
