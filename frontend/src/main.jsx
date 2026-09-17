import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import { ColorModeProvider } from './ColorMode';
import { initTelemetry } from './Telemetry';

try {
  initTelemetry();
} catch (err) {
  console.error('Error inicializando telemetría:', err);
}

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <ColorModeProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </ColorModeProvider>
  </React.StrictMode>,
);
