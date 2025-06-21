# React Sensor Dashboard

This project is a React-based web application that serves as a dashboard for monitoring sensor data from various rooms in a house. It fetches data from a backend API and visualizes temperature and humidity information.

## Project Structure

```
web-ui
├── public
│   └── index.html          # Main HTML file for the React application
├── src
│   ├── components          # Contains React components
│   │   ├── Dashboard.tsx   # Dashboard component that fetches and displays room data
│   │   └── RoomCard.tsx    # Component that displays individual room data
│   ├── App.tsx             # Main application component
│   ├── index.tsx           # Entry point of the React application
│   └── types               # TypeScript interfaces
│       └── index.ts
├── package.json            # npm configuration file
├── tsconfig.json           # TypeScript configuration file
└── README.md               # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd web-ui
   ```

2. **Install dependencies:**
   ```
   npm install
   ```

3. **Run the application:**
   ```
   npm start
   ```

   The application will be available at `http://localhost:3000`.

## Usage

- The dashboard will display a list of rooms with their respective temperature and humidity data.
- Users can select different time periods to visualize the data for each room.

## API

The application fetches data from the following API endpoints:

- `GET /api/room`: Retrieves a list of available rooms.
- `GET /api/room/{roomname}?period={period}`: Retrieves temperature and humidity data for a specific room over a specified period.

## License

This project is licensed under the MIT License.