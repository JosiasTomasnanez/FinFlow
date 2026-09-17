import { useState } from 'react'
import Button from '@mui/material/Button'
import Stack from '@mui/material/Stack'
import TextField from '@mui/material/TextField'
import { apiFetch, jsonPost } from '../api'
import LabStation from './LabStation'

export default function LatencyLab() {
  const [delayMs, setDelayMs] = useState(400)
  const [busy, setBusy] = useState(false)
  const [log, setLog] = useState('')
  const [error, setError] = useState(false)

  const fire = async () => {
    setBusy(true)
    setError(false)
    setLog('Esperando respuesta lenta...')
    try {
      const data = await apiFetch('/api/lab/slow', jsonPost({ delay_ms: delayMs }))
      setLog(`200 tras ${data.slept_ms} ms (p95 de exitosas)`)
    } catch (err) {
      setError(true)
      setLog(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <LabStation
      title="Latencia"
      hint="Sleep + HTTP 200. El p95 del dashboard filtra códigos 2xx/3xx; un 500 no contaría."
      log={log}
      error={error}
    >
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
        <TextField
          label="Delay (ms)"
          type="number"
          slotProps={{ htmlInput: { min: 1, max: 5000 } }}
          value={delayMs}
          onChange={(e) => setDelayMs(Number(e.target.value))}
          sx={{ maxWidth: 180 }}
        />
        <Button variant="contained" disabled={busy} onClick={fire}>
          Una request lenta
        </Button>
      </Stack>
    </LabStation>
  )
}
