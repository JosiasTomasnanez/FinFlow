import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './style.css';
import { initTelemetry } from './Telemetry';

try {
  initTelemetry();
} catch (err) {
  console.error('Error inicializando telemetría:', err);
}

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
