import { createTheme } from '@mui/material/styles'

const navy = '#0D3B66'
const ink = '#12202A'
const wash = '#F3F5F7'
const muteLight = '#3D4D57'
const forest = '#1E4A42'
const errorBank = '#8F1D2C'

const night = '#0B1218'
const nightPaper = '#121C26'
const steel = '#9BB4C8'
const muteDark = '#8A9AA6'

const font = '"Figtree", "Roboto", "Helvetica", "Arial", sans-serif'

export function createAppTheme(mode) {
  const dark = mode === 'dark'
  const hairline = dark ? 'rgba(232, 238, 242, 0.12)' : 'rgba(18, 32, 42, 0.1)'

  return createTheme({
    palette: {
      mode,
      primary: {
        main: navy,
        contrastText: '#FFFFFF',
      },
      secondary: {
        main: dark ? '#7FA89F' : forest,
      },
      error: {
        main: dark ? '#E07A86' : errorBank,
      },
      info: {
        main: dark ? steel : '#3A6A88',
      },
      success: {
        main: dark ? '#7FA89F' : forest,
      },
      background: {
        default: dark ? night : wash,
        paper: dark ? nightPaper : '#FFFFFF',
      },
      text: {
        primary: dark ? '#E8EEF2' : ink,
        secondary: dark ? muteDark : muteLight,
      },
      divider: hairline,
    },
    shape: {
      borderRadius: 12,
    },
    typography: {
      fontFamily: font,
      h4: { fontWeight: 700, letterSpacing: '-0.03em' },
      h5: { fontWeight: 600, letterSpacing: '-0.02em' },
      h6: { fontWeight: 600, letterSpacing: '-0.02em' },
      button: { fontWeight: 600, letterSpacing: 0 },
    },
    components: {
      MuiAppBar: {
        defaultProps: {
          color: 'primary',
          enableColorOnDark: true,
          elevation: 0,
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            textTransform: 'none',
          },
        },
      },
      MuiCard: {
        defaultProps: {
          variant: 'outlined',
          elevation: 0,
        },
        styleOverrides: {
          root: {
            borderColor: hairline,
          },
        },
      },
      MuiOutlinedInput: {
        styleOverrides: {
          notchedOutline: {
            borderColor: hairline,
          },
        },
      },
      MuiTableCell: {
        styleOverrides: {
          root: {
            borderColor: hairline,
          },
        },
      },
    },
  })
}
