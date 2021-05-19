import React, { useState, useEffect } from 'react';
import Slider from '@material-ui/core/Slider';

import {write} from '../../serialService/SerialService';

type ModeProps = {
    name: string,
    value: string,
}

const Mode: React.FC<ModeProps> = ({ name, value }: ModeProps) => {
    const [speed, setSpeed] = React.useState(20);

    const handleChange = (event: React.ChangeEvent<{}>, newSpeed: number | number[]) => {
        write(`speed:${newSpeed}`)
    };

    return (
    <div>
        <button onClick={() => write(value)}>{name}</button>
        <Slider
            defaultValue={20}
            aria-labelledby="discrete-slider-custom"
            step={10}
            valueLabelDisplay="auto"
            onChange={handleChange}
        />
    </div>
    );
};

export default Mode;
