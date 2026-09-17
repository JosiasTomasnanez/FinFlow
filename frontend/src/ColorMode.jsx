import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { ThemeProvider } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import { createAppTheme } from './theme'

const STORAGE_KEY = 'finflow-color-mode'
const ColorModeContext = createContext({
  mode: 'light',
  toggleColorMode: () => {},
})

export function useColorMode() {
  return useContext(ColorModeContext)
}

function readStoredMode() {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored === 'dark' || stored === 'light') return stored
  } catch {
    /* ignore */
  }
  return 'light'
}

export function ColorModeProvider({ children }) {
  const [mode, setMode] = useState(readStoredMode)
  const theme = useMemo(() => createAppTheme(mode), [mode])

  useEffect(() => {
    document.documentElement.style.colorScheme = mode
  }, [mode])

  const value = useMemo(
    () => ({
      mode,
      toggleColorMode: () => {
        setMode((current) => {
          const next = current === 'light' ? 'dark' : 'light'
          try {
            localStorage.setItem(STORAGE_KEY, next)
          } catch {
            /* ignore */
          }
          document.documentElement.style.colorScheme = next
          return next
        })
      },
    }),
    [mode],
  )

  return (
    <ColorModeContext.Provider value={value}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </ThemeProvider>
    </ColorModeContext.Provider>
  )
}
