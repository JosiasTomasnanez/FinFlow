import { useEffect, useMemo, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import AppBar from '@mui/material/AppBar'
import Box from '@mui/material/Box'
import Container from '@mui/material/Container'
import IconButton from '@mui/material/IconButton'
import Tab from '@mui/material/Tab'
import Tabs from '@mui/material/Tabs'
import Toolbar from '@mui/material/Toolbar'
import Tooltip from '@mui/material/Tooltip'
import Typography from '@mui/material/Typography'
import DarkModeOutlinedIcon from '@mui/icons-material/DarkModeOutlined'
import LightModeOutlinedIcon from '@mui/icons-material/LightModeOutlined'
import { apiFetch } from '../api'
import { useColorMode } from '../ColorMode'

export default function AppShell() {
  const { mode, toggleColorMode } = useColorMode()
  const { pathname } = useLocation()
  const navigate = useNavigate()
  const [featureEnabled, setFeatureEnabled] = useState(false)
  const [labEnabled, setLabEnabled] = useState(true)
  const [user, setUser] = useState(null)

  useEffect(() => {
    let cancelled = false
    apiFetch('/api/flags')
      .then((data) => {
        if (cancelled) return
        setFeatureEnabled(Boolean(data.feature_login))
        setLabEnabled(data.feature_lab !== false)
      })
      .catch(() => {
        if (cancelled) return
        setFeatureEnabled(false)
        setLabEnabled(true)
      })
    return () => {
      cancelled = true
    }
  }, [])

  const outletContext = useMemo(
    () => ({ featureEnabled, labEnabled, user, setUser }),
    [featureEnabled, labEnabled, user],
  )

  const section = pathname.startsWith('/lab') ? '/lab' : '/'

  return (
    <Box sx={{ minHeight: '100vh' }}>
      <AppBar position="sticky">
        <Container maxWidth="lg">
          <Toolbar disableGutters>
            <Typography variant="h6" component="p" sx={{ mr: 2 }}>
              FinFlow
            </Typography>
            <Tabs
              value={section}
              onChange={(_, value) => navigate(value)}
              textColor="inherit"
              slotProps={{ indicator: { sx: { bgcolor: 'common.white' } } }}
            >
              <Tab value="/" label="Cuentas" />
              {labEnabled ? <Tab value="/lab" label="Laboratorio" /> : null}
            </Tabs>
            <Box sx={{ flexGrow: 1 }} />
            <Tooltip title={mode === 'light' ? 'Modo oscuro' : 'Modo claro'}>
              <IconButton color="inherit" onClick={toggleColorMode} aria-label="Cambiar tema">
                {mode === 'light' ? <DarkModeOutlinedIcon /> : <LightModeOutlinedIcon />}
              </IconButton>
            </Tooltip>
          </Toolbar>
        </Container>
      </AppBar>
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Outlet context={outletContext} />
      </Container>
    </Box>
  )
}
