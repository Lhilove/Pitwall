# Pitwall 🏎️

A real-time F1 telemetry backend built with Go. Ingests, processes, and streams live driver telemetry data via WebSockets.

## Features

- REST API for telemetry ingestion and querying
- Real-time WebSocket feed for live telemetry streaming
- Background worker pool for async processing
- Buffered queue for handling telemetry bursts
- Driver leaderboard and aggregated stats
- Paginated telemetry with sorting support
- Built-in race simulator for testing

## Tech Stack

- **Go** + Gin (HTTP framework)
- **PostgreSQL** (storage)
- **Gorilla WebSocket** (real-time feed)
- **Docker** (database)

## Getting Started

### Prerequisites

- Go 1.21+
- Docker

### Setup

1. Clone the repo

```bash
git clone https://github.com/lhilove/pitwall.git
cd pitwall
```

2. Start the database

```bash
docker compose up -d
```

3. Create a `.env` file

```env
db=postgres://pitwalluser:pitwallpass@localhost:5432/f1telemetry
PORT=8080
```

4. Run the server

```bash
go run cmd/main.go
```

## API Reference

### Telemetry

| Method | Endpoint                 | Description                                            |
|--------|----------                |-------------                                           |
| `POST` | `/telemetry`             | Submit a telemetry record                              |
| `GET`  | `/telemetry`             | Get all telemetry (filter by `?driver=`)               |
| `GET`  | `/telemetry/all`         | Paginated telemetry (`?page=&limit=&sort=&order=`)     |
| `GET`  | `/telemetry/stats`       | Aggregated stats (avg speed, max speed, total records) |
| `GET`  | `/telemetry/leaderboard` | Driver leaderboard by average speed                    |

### System

| Method | Endpoint       | Description                                            |
|--------|----------      |-------------                                           |
| `POST` | `/simulate`    | Start a telemetry simulation (500 records, Verstappen) |
| `GET`  | `/queue/stats` | Queue length and capacity                              |
| `GET`  | `/metrics`     | Total processed records by workers                     |
| `GET`  | `/ws`          | WebSocket endpoint for live telemetry feed             |

### Telemetry Payload

```json
{
  "driver": "Verstappen",
  "lap": 12,
  "speed": 312,
  "throttle": 87,
  "brake": 4,
  "gear": 6
}
```

## WebSocket

Connect to `ws://localhost:8080/ws` to receive live telemetry as JSON the moment each record is processed and saved.


## Project Structure

```
pitwall/
├── cmd/
│   └── main.go
├── internal/
│   ├── db/
│   ├── handlers/
│   ├── models/
│   ├── queue/
│   ├── repository/
│   ├── service/
│   ├── websocket/
│   └── workers/
├── Simulator/
├── docker-compose.yml
└── .env
```

## Roadmap

- [ ] OpenF1 API integration for live race data
- [ ] Per-lap fastest sector times
- [ ] Driver head-to-head comparison endpoint
- [ ] Redis queue for persistence across restarts
- [ ] Live telemetry dashboard frontend