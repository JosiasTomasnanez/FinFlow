import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'

export default function PageHeader({ title, children }) {
  return (
    <Box sx={{ mb: 2 }}>
      <Typography variant="h4" component="h1">
        {title}
      </Typography>
      {children ? (
        <Typography color="text.secondary" sx={{ mt: 1 }}>
          {children}
        </Typography>
      ) : null}
    </Box>
  )
}
