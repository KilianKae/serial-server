import React, { useState, useEffect } from 'react';
import logo from './logo.svg';
import './App.css';

import { getStatus, IStatus } from './SerialService/serialService';
function App() {
  let [status, setStatus] = useState<IStatus | undefined>(undefined);

  useEffect(() => {
    getStatus().then((status) => setStatus(status));
  }, []);

  return (
    <div className='App'>
      <p>
        Name: {status ? status.name : 'Loading...'}
        <br />
        Baud: {status ? status.baud : 'Loading...'}
        <br />
        {status?.error ? `Error: ${status.error}` : 'Loading...'}
        <br />
      </p>
    </div>
  );
}

export default App;
