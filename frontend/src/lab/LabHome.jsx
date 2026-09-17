import { Link as RouterLink } from 'react-router-dom'
import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import CardActionArea from '@mui/material/CardActionArea'
import CardContent from '@mui/material/CardContent'
import Grid from '@mui/material/Grid'
import Typography from '@mui/material/Typography'
import { useTheme } from '@mui/material/styles'

export default function LabHome() {
  const theme = useTheme()
  const signals = [
    {
      to: '/lab/traffic',
      title: 'Tráfico',
      body: 'Ráfaga de GET 2xx. Sube el RPS; 5xx y saturación no deberían moverse.',
      accent: theme.palette.primary.main,
    },
    {
      to: '/lab/latency',
      title: 'Latencia',
      body: 'Una request lenta con HTTP 200. Entra al p95 de respuestas exitosas.',
      accent: theme.palette.info.main,
    },
    {
      to: '/lab/errors',
      title: 'Errores',
      body: 'Inyectá 500. Tasa 5xx y barras por código. El p95 de 2xx casi no cambia.',
      accent: theme.palette.error.main,
    },
    {
      to: '/lab/saturation',
      title: 'Saturación',
      body: 'CPU sostenida, picos CFS y retención de RAM contra los límites del contenedor.',
      accent: theme.palette.success.main,
    },
  ]

  return (
    <Grid container spacing={2}>
      {signals.map((signal) => (
        <Grid key={signal.to} size={{ xs: 12, sm: 6 }}>
          <Card
            sx={{
              '&:hover': { transform: 'translateY(-3px)', borderColor: signal.accent },
            }}
          >
            <CardActionArea component={RouterLink} to={signal.to} sx={{ height: '100%' }}>
              <CardContent sx={{ minHeight: 176, p: 3 }}>
                <Box sx={{ width: 28, height: 3, borderRadius: 99, bgcolor: signal.accent, mb: 2 }} />
                <Typography variant="h5" component="h2" gutterBottom>
                  {signal.title}
                </Typography>
                <Typography color="text.secondary">{signal.body}</Typography>
              </CardContent>
            </CardActionArea>
          </Card>
        </Grid>
      ))}
    </Grid>
  )
}
