import React, { useEffect, useState } from 'react';
import RoomCard from './RoomCard';
import { RoomData } from '../types';
import './Dashboard.css';

const Dashboard: React.FC = () => {
    const [rooms, setRooms] = useState<string[]>([]);
    const [roomData, setRoomData] = useState<Record<string, RoomData[]>>({});

    const API_BASE = 'http://${window.location.hostname}:8080';

    useEffect(() => {
        const fetchRooms = async () => {
            try {
                const response = await fetch(`${API_BASE}/api/room`);
                const data = await response.json();
                setRooms(data);
            } catch (error) {
                console.error('Error fetching rooms:', error);
            }
        };

        fetchRooms();
    }, []);

    return (
        <div className="dashboard">
            {rooms.map(room => (
                <RoomCard key={room} room={room} />
            ))}
        </div>
    );
};

export default Dashboard;