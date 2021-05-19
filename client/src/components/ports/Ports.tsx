import React, { useState, useEffect } from 'react';
import './Ports.css';

import {getPorts, IPorts } from '../../serialService/SerialService';

const Ports: React.FC = () => {
    let [ports, setPorts] = useState<IPorts[]>([]);

    useEffect(() => {
        getPorts().then((ports) => setPorts(ports));
    }, []);

    let portsRows = [];
    for (let port of ports) {
        portsRows.push(<option value={port.Name}>{port.Name}</option>);
    }

    return (
        <form>
            <label htmlFor="ports">Ports:</label>
            <select name="ports" id="ports">
                {portsRows}
            </select>
            <input type="submit" value="Submit"/>
        </form>
    );
};

export default Ports;
