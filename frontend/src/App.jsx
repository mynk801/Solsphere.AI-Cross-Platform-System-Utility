
import { useState, useEffect } from 'react';

const API_URL = 'http://127.0.0.1:5000/api/devices';

function StatusIndicator({ label, isOk }) {
  const statusText = isOk ? 'OK' : 'Issue';
  const statusColor = isOk ? 'green' : 'red';

  return (
    <div>
      {label}: <span style={{ color: statusColor, fontWeight: 'bold' }}>{statusText}</span>
    </div>
  );
}

function App() {
  const [devices, setDevices] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchDevices = async () => {
      try {
        const response = await fetch(API_URL);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setDevices(data); 
      } catch (e) {
        setError('Failed to fetch data from the backend. Is the server running?');
        console.error(e);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDevices();
  }, []); 


  if (isLoading) {
    return <div>Loading device data...</div>;
  }

  if (error) {
    return <div style={{ color: 'red' }}>Error: {error}</div>;
  }

  return (
    <div className="dashboard-container">
      <h1>System Health Dashboard</h1>
      {devices.length === 0 ? (
        <p>No devices are reporting yet. Run the utility on a machine to see its data here.</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>Machine ID</th>
              <th>Operating System</th>
              <th>Last Check-in</th>
              <th>Health Status</th>
            </tr>
          </thead>
          <tbody>
            {devices.map((device) => (
              <tr key={device.machineId}>
                <td>{device.machineId}</td>
                <td>{device.os}</td>
                <td>{new Date(device.lastCheckIn).toLocaleString()}</td>
                <td>
                  <StatusIndicator label="Disk Encrypted" isOk={device.latestData.diskEncrypted} />
                  <StatusIndicator label="OS Up-to-date" isOk={device.latestData.osUpToDate} />
                  <StatusIndicator label="Antivirus Active" isOk={device.latestData.antivirusActive} />
                  <StatusIndicator label="Sleep Settings" isOk={device.latestData.sleepSettingsOk} />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default App;