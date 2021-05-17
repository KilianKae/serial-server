import React, { useState, useEffect } from 'react';
import './App.css';

import { getStatus, getPorts, IStatus, IPorts } from './SerialService/serialService';

function App() {
  let [status, setStatus] = useState<IStatus | undefined>(undefined);
  let [ports, setPorts] = useState<IPorts[]>([]);

  useEffect(() => {
    getStatus().then((status) => setStatus(status));
    getPorts().then((ports) => setPorts(ports));
  }, []);

  let portsRows = [];

  for (let port of ports) {
    portsRows.push(<option value={port.Name}>{port.Name}</option>);
  }

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
      <form>
        <label htmlFor="ports">Ports:</label>
        <select name="ports" id="ports">
          {portsRows}
        </select>
          <input type="submit" value="Submit"/>
      </form>
    </div>
  );
}

export default App;
