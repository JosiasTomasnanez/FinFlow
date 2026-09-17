import { Navigate, Outlet, useLocation, useNavigate, useOutletContext } from 'react-router-dom'
import Stack from '@mui/material/Stack'
import Tab from '@mui/material/Tab'
import Tabs from '@mui/material/Tabs'
import PageHeader from '../layout/PageHeader'

const TABS = [
  { to: '/lab', label: 'Overview', end: true },
  { to: '/lab/traffic', label: 'Tráfico' },
  { to: '/lab/latency', label: 'Latencia' },
  { to: '/lab/errors', label: 'Errores' },
  { to: '/lab/saturation', label: 'Saturación' },
]

export default function LabLayout() {
  const { labEnabled } = useOutletContext()
  const { pathname } = useLocation()
  const navigate = useNavigate()

  if (!labEnabled) {
    return <Navigate to="/" replace />
  }

  const current = TABS.find((tab) => (tab.end ? pathname === tab.to : pathname.startsWith(tab.to))) ?? TABS[0]

  return (
    <Stack spacing={3}>
      <PageHeader title="Golden Signals">
        Una pantalla por señal. Grafana: tráfico y errores en 1m, saturación en 2m.
      </PageHeader>
      <Tabs
        value={current.to}
        onChange={(_, value) => navigate(value)}
        variant="scrollable"
        scrollButtons="auto"
        sx={{ borderBottom: 1, borderColor: 'divider' }}
      >
        {TABS.map((tab) => (
          <Tab key={tab.to} value={tab.to} label={tab.label} />
        ))}
      </Tabs>
      <Outlet />
    </Stack>
  )
}
