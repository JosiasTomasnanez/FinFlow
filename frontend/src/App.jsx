import { useEffect, useState } from 'react'

const apiFetch = async (path, options = {}) => {
  const response = await fetch(path, options)
  const data = await response.json().catch(() => null)
  if (!response.ok) {
    throw new Error(data?.error || response.statusText)
  }
  return data
}

function App() {
  const [wallets, setWallets] = useState([])
  const [owner, setOwner] = useState('alice')
  const [initialBalance, setInitialBalance] = useState(1000)
  const [createResult, setCreateResult] = useState('')
  const [payment, setPayment] = useState({ from_wallet_id: '', to_wallet_id: '', amount: 100 })
  const [paymentResult, setPaymentResult] = useState('')

  const loadWallets = async () => {
    try {
      const data = await apiFetch('/api/wallets')
      setWallets(data)
    } catch (error) {
      setWallets([])
      setCreateResult(`Error al cargar wallets: ${error.message}`)
    }
  }

  useEffect(() => {
    loadWallets()
  }, [])

  const handleCreateWallet = async (event) => {
    event.preventDefault()
    setCreateResult('Creando wallet...')

    try {
      const wallet = await apiFetch('/api/wallets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ owner, initial_balance: Number(initialBalance) }),
      })
      setCreateResult(JSON.stringify(wallet, null, 2))
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
      const result = await apiFetch('/api/payments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payment),
      })
      setPaymentResult(JSON.stringify(result, null, 2))
      setPayment({ ...payment, amount: 100 })
      loadWallets()
    } catch (error) {
      setPaymentResult(`Error: ${error.message}`)
    }
  }

  return (
    <div className="page">
      <header>
        <h1>FinFlow</h1>
        <p>Demo sencilla con backend Gin y frontend React.</p>
      </header>

      <main>
        <section>
          <h2>Crear Wallet</h2>
          <form onSubmit={handleCreateWallet}>
            <label>
              Dueño
              <input value={owner} onChange={(e) => setOwner(e.target.value)} required />
            </label>
            <label>
              Saldo inicial
              <input
                type="number"
                value={initialBalance}
                onChange={(e) => setInitialBalance(Number(e.target.value))}
                required
              />
            </label>
            <button type="submit">Crear wallet</button>
          </form>
          <pre>{createResult}</pre>
        </section>

        <section>
          <h2>Wallets</h2>
          <button onClick={loadWallets}>Actualizar lista</button>
          <pre>{JSON.stringify(wallets, null, 2)}</pre>
        </section>

        <section>
          <h2>Transferencia</h2>
          <form onSubmit={handlePayment}>
            <label>
              Desde wallet ID
              <input
                value={payment.from_wallet_id}
                onChange={(e) => setPayment({ ...payment, from_wallet_id: e.target.value })}
                required
              />
            </label>
            <label>
              Hacia wallet ID
              <input
                value={payment.to_wallet_id}
                onChange={(e) => setPayment({ ...payment, to_wallet_id: e.target.value })}
                required
              />
            </label>
            <label>
              Monto
              <input
                type="number"
                value={payment.amount}
                onChange={(e) => setPayment({ ...payment, amount: Number(e.target.value) })}
                required
              />
            </label>
            <button type="submit">Enviar pago</button>
          </form>
          <pre>{paymentResult}</pre>
        </section>
      </main>
    </div>
  )
}

export default App
