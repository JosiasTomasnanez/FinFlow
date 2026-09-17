import { useState } from 'react'
import Button from '@mui/material/Button'
import Stack from '@mui/material/Stack'
import TextField from '@mui/material/TextField'
import { apiFetch } from '../api'
import { runParallel } from './runParallel'
import LabStation from './LabStation'

export default function TrafficLab() {
  const [count, setCount] = useState(40)
  const [busy, setBusy] = useState(false)
  const [log, setLog] = useState('')
  const [error, setError] = useState(false)

  const fire = async () => {
    setBusy(true)
    setError(false)
    setLog('Disparando...')
    try {
      const { ok, fail, total } = await runParallel(count, () => apiFetch('/api/lab/ok'))
      setLog(`${ok}/${total} OK` + (fail ? `, ${fail} fallaron` : ''))
      setError(Boolean(fail))
    } catch (err) {
      setError(true)
      setLog(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <LabStation
      title="Tráfico"
      hint="GET /api/lab/ok en paralelo. En Grafana: RPS y series por código."
      log={log}
      error={error}
    >
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
        <TextField
          label="Requests"
          type="number"
          slotProps={{ htmlInput: { min: 1, max: 100 } }}
          value={count}
          onChange={(e) => setCount(Number(e.target.value))}
          sx={{ maxWidth: 180 }}
        />
        <Button variant="contained" disabled={busy} onClick={fire}>
          Disparar ráfaga 2xx
        </Button>
      </Stack>
    </LabStation>
  )
}
