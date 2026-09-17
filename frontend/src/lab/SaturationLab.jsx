import { useState } from 'react'
import Button from '@mui/material/Button'
import Divider from '@mui/material/Divider'
import Stack from '@mui/material/Stack'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'
import { apiFetch, jsonPost } from '../api'
import { runParallel } from './runParallel'
import LabStation from './LabStation'

export default function SaturationLab() {
  const [cpuMs, setCpuMs] = useState(5000)
  const [burstN, setBurstN] = useState(20)
  const [mib, setMib] = useState(64)
  const [busy, setBusy] = useState(false)
  const [log, setLog] = useState('')
  const [error, setError] = useState(false)

  const run = async (label, fn) => {
    setBusy(true)
    setError(false)
    setLog(label)
    try {
      await fn()
    } catch (err) {
      setError(true)
      setLog(err.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <LabStation
      title="Saturación"
      hint="cAdvisor, ventana [2m]. CPU larga también sube RPS y p95. La memoria se retiene hasta Liberar."
      log={log}
      error={error}
    >
      <Stack spacing={3} divider={<Divider />}>
        <Stack spacing={1.5}>
          <Typography variant="subtitle1">CPU sostenida</Typography>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
            <TextField
              label="Duración (ms)"
              type="number"
              slotProps={{ htmlInput: { min: 100, max: 10000 } }}
              value={cpuMs}
              onChange={(e) => setCpuMs(Number(e.target.value))}
              sx={{ maxWidth: 200 }}
            />
            <Button
              variant="contained"
              disabled={busy}
              onClick={() =>
                run('Quemando CPU...', async () => {
                  const data = await apiFetch('/api/lab/cpu', jsonPost({ duration_ms: cpuMs }))
                  setLog(`CPU spin ${data.spun_ms} ms — mirá uso/límite`)
                })
              }
            >
              Quemar CPU
            </Button>
          </Stack>
        </Stack>

        <Stack spacing={1.5}>
          <Typography variant="subtitle1">Throttle CFS</Typography>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
            <TextField
              label="Ráfagas en paralelo"
              type="number"
              slotProps={{ htmlInput: { min: 1, max: 40 } }}
              value={burstN}
              onChange={(e) => setBurstN(Number(e.target.value))}
              sx={{ maxWidth: 200 }}
            />
            <Button
              variant="contained"
              color="secondary"
              disabled={busy}
              onClick={() =>
                run('Ráfagas CFS...', async () => {
                  const { ok, total } = await runParallel(burstN, () =>
                    apiFetch('/api/lab/burst', jsonPost({ spin_ms: 80 })),
                  )
                  setLog(`${ok}/${total} bursts de 80 ms — mirá throttle CFS`)
                })
              }
            >
              Disparar picos
            </Button>
          </Stack>
        </Stack>

        <Stack spacing={1.5}>
          <Typography variant="subtitle1">Memoria</Typography>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }} flexWrap="wrap">
            <TextField
              label="MiB a retener"
              type="number"
              slotProps={{ htmlInput: { min: 1, max: 128 } }}
              value={mib}
              onChange={(e) => setMib(Number(e.target.value))}
              sx={{ maxWidth: 200 }}
            />
            <Button
              variant="contained"
              disabled={busy}
              onClick={() =>
                run('Reteniendo RAM...', async () => {
                  const data = await apiFetch('/api/lab/memory', jsonPost({ megabytes: mib }))
                  setLog(`Hold ${data.held_mib} MiB — working set no baja hasta Liberar`)
                })
              }
            >
              Retener
            </Button>
            <Button
              variant="outlined"
              disabled={busy}
              onClick={() =>
                run('Liberando...', async () => {
                  await apiFetch('/api/lab/memory/release', jsonPost({}))
                  setLog('Hold en 0 (el working set puede tardar en bajar)')
                })
              }
            >
              Liberar
            </Button>
          </Stack>
        </Stack>
      </Stack>
    </LabStation>
  )
}
