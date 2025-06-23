// filepath: web-ui/src/components/RoomCard.tsx
import React, { useEffect, useState } from 'react';
import { RoomData, SensorData } from '../types';
import { Line } from 'react-chartjs-2';
import { Chart, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend } from 'chart.js';
import './RoomCard.css';

Chart.register(LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Legend);

interface RoomCardProps {
    room: string;
}

const RoomCard: React.FC<RoomCardProps> = ({ room }) => {
    const [data, setData] = useState<SensorData[]>([]);
    const [period, setPeriod] = useState<string>('hour');

    const fetchRoomData = async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/room/${room}?period=${period}`);
            const result = await response.json();
            if( result != null){
                setData(result);
            }
        } catch (error) {
            console.error(`Error fetching data for ${room}:`, error);
        }
    };

    useEffect(() => {
        fetchRoomData();
    }, [room, period]);

    return (
        <div className="card">
            <div className="card-header">
                <h2 className="card-title">{room}</h2>
                <span className="card-period">
                    <label htmlFor={`period-${room}`}>Period:</label>
                    <select
                        id={`period-${room}`}
                        value={period}
                        onChange={(e) => setPeriod(e.target.value)}
                    >
                        <option value="hour">Hour</option>
                        <option value="day">Day</option>
                        <option value="week">Week</option>
                        <option value="month">Month</option>
                        <option value="year">Year</option>
                    </select>
                </span>
            </div>
            <Line
                data={{
                    labels: data.map(d => {
                        const date = new Date(d.period + 'Z');
                        return date.toLocaleTimeString();
                }),
                    datasets: [
                        {
                            label: 'Temperature',
                            data: data.map(d => d.max_temp),
                            borderColor: 'red',
                            fill: false,
                        },
                        {
                            label: 'Humidity',
                            data: data.map(d => d.max_humidity),
                            borderColor: 'blue',
                            fill: false,
                        }
                    ]
                }}
            />
        </div>
    );
};

export default RoomCard;