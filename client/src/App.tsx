import React, { useState, useEffect } from 'react';
import './App.css';

import { getStatus, IStatus } from './serialService/SerialService';
import Ports from "./components/ports/Ports";
import Mode from "./components/mode/Mode";

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
        {status?.error ? `Error: ${status.error}` : 'Everything is in order'}
        <br />
      </p>
      <Ports />
      <Mode name="Random Walk" value="randomWalk"/>
    </div>
  );
}

export default App;
