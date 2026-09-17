import Alert from '@mui/material/Alert'
import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import Stack from '@mui/material/Stack'
import Typography from '@mui/material/Typography'

export default function LabStation({ title, hint, children, log, error }) {
  return (
    <Card>
      <CardContent sx={{ p: { xs: 2.5, md: 3.5 } }}>
        <Typography variant="h5" component="h2" gutterBottom>
          {title}
        </Typography>
        <Typography color="text.secondary" sx={{ mb: 3, maxWidth: 640, lineHeight: 1.65 }}>
          {hint}
        </Typography>
        <Box>{children}</Box>
        {log ? (
          <Stack sx={{ mt: 2.5 }}>
            <Alert severity={error ? 'error' : 'success'}>
              {log}
            </Alert>
          </Stack>
        ) : null}
      </CardContent>
    </Card>
  )
}
