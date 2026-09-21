import { useState } from 'react';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import { apiRequest, jsonPost } from '../api';
import { runParallel } from './runParallel';
import LabStation from './LabStation';
import { logEvent } from '../Logger';

export default function ErrorsLab() {
  const [count, setCount] = useState(10);
  const [busy, setBusy] = useState(false);
  const [log, setLog] = useState('');
  const [error, setError] = useState(false);

  const fire = async () => {
    setBusy(true);
    setError(false);
    setLog('Inyectando 500...');
    try {
      const { ok, total } = await runParallel(count, async () => {
        const { response } = await apiRequest('/api/lab/error', jsonPost({ code: 500 }));
        if (response.status < 500) {
          throw new Error(`status ${response.status}`);
        }
      });
      setLog(`${ok}/${total} HTTP 5xx`);
    } catch (err) {
      setError(true);
      setLog(err.message);
    } finally {
      setBusy(false);
    }
  };

  return (
    <LabStation
      title="Errores"
      hint="POST /api/lab/error → 500. Grafana: tasa 5xx y barras por código. El p95 de exitosas casi no cambia."
      log={log}
      error={error}
    >
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
        <TextField
          label="Cantidad"
          type="number"
          slotProps={{ htmlInput: { min: 1, max: 50 } }}
          value={count}
          onChange={(e) => setCount(Number(e.target.value))}
          sx={{ maxWidth: 180 }}
        />
        <Button variant="contained" color="error" disabled={busy} onClick={fire}>
          Tirar 5xx
        </Button>
      </Stack>
    </LabStation>
  );
}
