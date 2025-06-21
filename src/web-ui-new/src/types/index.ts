export interface RoomData {
    name: string;
    sensorData: SensorData[];
}

export interface SensorData {
    period: string;
    max_temp: number;
    max_humidity: number;
}