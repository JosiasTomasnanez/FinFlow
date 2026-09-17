import { useEffect, useState } from 'react'
import { useOutletContext } from 'react-router-dom'
import Alert from '@mui/material/Alert'
import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import Grid from '@mui/material/Grid'
import Stack from '@mui/material/Stack'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'
import Login from '../Login'
import PageHeader from '../layout/PageHeader'
import { apiFetch, jsonPost } from '../api'

export default function HomePage() {
  const { featureEnabled, user, setUser } = useOutletContext()
  const [wallets, setWallets] = useState([])
  const [owner, setOwner] = useState('alice')
  const [initialBalance, setInitialBalance] = useState(1000)
  const [createResult, setCreateResult] = useState('')
  const [payment, setPayment] = useState({ from_wallet_id: '', to_wallet_id: '', amount: 100 })
  const [paymentResult, setPaymentResult] = useState('')
  const [listError, setListError] = useState('')

  const loadWallets = async () => {
    try {
      const data = await apiFetch('/api/wallets')
      setWallets(Array.isArray(data) ? data : [])
      setListError('')
    } catch (error) {
      setWallets([])
      setListError(`No se pudieron cargar las wallets: ${error.message}`)
    }
  }

  useEffect(() => {
    loadWallets()
  }, [])

  const handleCreateWallet = async (event) => {
    event.preventDefault()
    setCreateResult('Creando wallet...')
    try {
      const wallet = await apiFetch('/api/wallets', jsonPost({ owner, initial_balance: Number(initialBalance) }))
      setCreateResult(`Wallet ${wallet.id} para ${wallet.owner}`)
      setOwner('')
      setInitialBalance(1000)
      loadWallets()
    } catch (error) {
      setCreateResult(`Error: ${error.message}`)
    }
  }

  const handlePayment = async (event) => {
    event.preventDefault()
    setPaymentResult('Procesando pago...')
    try {
      const result = await apiFetch('/api/payments', jsonPost(payment))
      setPaymentResult(`Pago ${result.status ?? 'ok'} · ${result.amount}`)
      setPayment({ ...payment, amount: 100 })
      loadWallets()
    } catch (error) {
      setPaymentResult(`Error: ${error.message}`)
    }
  }

  return (
    <Stack spacing={3}>
      <PageHeader title="Cuentas">Wallets y transferencias.</PageHeader>

      <Login onLogin={setUser} featureEnabled={featureEnabled} />

      <Grid container spacing={2.5}>
        <Grid size={{ xs: 12, md: 5 }}>
          <Card>
            <CardContent sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Nueva wallet
              </Typography>
              <Box component="form" onSubmit={handleCreateWallet}>
                <Stack spacing={2}>
                  <TextField label="Dueño" value={owner} onChange={(e) => setOwner(e.target.value)} required fullWidth />
                  <TextField
                    label="Saldo inicial"
                    type="number"
                    value={initialBalance}
                    onChange={(e) => setInitialBalance(Number(e.target.value))}
                    required
                    fullWidth
                  />
                  <Button type="submit" variant="contained">
                    Crear
                  </Button>
                  {createResult ? <Alert severity="info">{createResult}</Alert> : null}
                </Stack>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid size={{ xs: 12, md: 7 }}>
          <Card>
            <CardContent sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Transferir
              </Typography>
              <Box component="form" onSubmit={handlePayment}>
                <Stack spacing={2}>
                  <TextField
                    label="Desde wallet ID"
                    value={payment.from_wallet_id}
                    onChange={(e) => setPayment({ ...payment, from_wallet_id: e.target.value })}
                    required
                    fullWidth
                  />
                  <TextField
                    label="Hacia wallet ID"
                    value={payment.to_wallet_id}
                    onChange={(e) => setPayment({ ...payment, to_wallet_id: e.target.value })}
                    required
                    fullWidth
                  />
                  <TextField
                    label="Monto"
                    type="number"
                    value={payment.amount}
                    onChange={(e) => setPayment({ ...payment, amount: Number(e.target.value) })}
                    required
                    fullWidth
                  />
                  <Button type="submit" variant="contained" color="secondary">
                    Enviar
                  </Button>
                  {paymentResult ? <Alert severity="info">{paymentResult}</Alert> : null}
                </Stack>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Card>
        <CardContent sx={{ p: 3 }}>
          <Stack
            direction={{ xs: 'column', sm: 'row' }}
            spacing={1}
            justifyContent="space-between"
            alignItems={{ sm: 'center' }}
            sx={{ mb: 2 }}
          >
            <Box>
              <Typography variant="h6">Wallets</Typography>
              <Typography variant="body2" color="text.secondary">
                {user ? user.username : 'anónimo'}
              </Typography>
            </Box>
            <Button variant="outlined" onClick={loadWallets}>
              Actualizar
            </Button>
          </Stack>
          {listError ? (
            <Alert severity="error">{listError}</Alert>
          ) : wallets.length === 0 ? (
            <Alert severity="info">Todavía no hay wallets. Creá la primera arriba.</Alert>
          ) : (
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Dueño</TableCell>
                  <TableCell align="right">Saldo</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {wallets.map((wallet) => (
                  <TableRow key={wallet.id} hover>
                    <TableCell>{wallet.id}</TableCell>
                    <TableCell>{wallet.owner}</TableCell>
                    <TableCell align="right">{wallet.balance}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </Stack>
  )
}
