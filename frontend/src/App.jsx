import { Navigate, Route, Routes } from 'react-router-dom';
import AppShell from './layout/AppShell';
import HomePage from './pages/HomePage';
import LabLayout from './lab/LabLayout';
import LabHome from './lab/LabHome';
import TrafficLab from './lab/TrafficLab';
import LatencyLab from './lab/LatencyLab';
import ErrorsLab from './lab/ErrorsLab';
import SaturationLab from './lab/SaturationLab';

export default function App() {
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/lab" element={<LabLayout />}>
          <Route index element={<LabHome />} />
          <Route path="traffic" element={<TrafficLab />} />
          <Route path="latency" element={<LatencyLab />} />
          <Route path="errors" element={<ErrorsLab />} />
          <Route path="saturation" element={<SaturationLab />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  );
}
