# Pitwall 🏎️

A real-time F1 telemetry backend built with Go. Ingests, processes, and streams live driver telemetry data via WebSockets — fed by real session data from FastF1.

## Features

- REST API for telemetry ingestion and querying
- Real-time WebSocket feed for live telemetry streaming
- Background worker pool for async processing
- Buffered queue for handling telemetry bursts
- Driver leaderboard and aggregated stats
- Driver-to-driver comparison endpoint
- Paginated telemetry with sorting support
- Built-in race simulator for testing
- Python feeder script pulling real session data from FastF1

## Tech Stack

- **Go** + Gin (HTTP framework)
- **PostgreSQL** (storage)
- **Gorilla WebSocket** (real-time feed)
- **Docker** (database)
- **Python + FastF1** (live data feeder)

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

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/telemetry` | Submit a telemetry record |
| `GET` | `/telemetry` | Get all telemetry (filter by `?driver=`) |
| `GET` | `/telemetry/all` | Paginated telemetry (`?page=&limit=&sort=&order=`) |
| `GET` | `/telemetry/stats` | Aggregated stats (avg speed, max speed, total records) |
| `GET` | `/telemetry/leaderboard` | Driver leaderboard by average speed |
| `GET` | `/telemetry/compare?a=&b=` | Compare aggregated stats between two drivers |

### System

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/simulate` | Start a telemetry simulation (500 records, Verstappen) |
| `GET` | `/queue/stats` | Queue length and capacity |
| `GET` | `/metrics` | Total processed records by workers |
| `GET` | `/ws` | WebSocket endpoint for live telemetry feed |

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

```js
const ws = new WebSocket("ws://localhost:8080/ws")
ws.onmessage = (e) => console.log(JSON.parse(e.data))
```

## Live data feed (FastF1)

A small Python sidecar (`feeder.py`) pulls real session data using [FastF1](https://github.com/theOehrly/Fast-F1) and streams it into the Go backend via the `/telemetry` endpoint — no changes needed on the Go side.

```bash
pip install fastf1 requests
mkdir cache
python feeder.py
```

The feeder loads a real race session, walks through each driver's lap-by-lap car data, and POSTs sampled telemetry points with a small delay so it streams in like a live broadcast rather than dumping all at once.

```
FastF1 (Python) → POST /telemetry → Go queue → worker → PostgreSQL + WebSocket broadcast
```

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
├── feeder.py
├── docker-compose.yml
└── .env
```

## Roadmap

- [x] OpenF1 / FastF1 integration for real race data
- [x] Driver head-to-head comparison endpoint
- [ ] Per-lap fastest sector times
- [ ] Redis queue for persistence across restarts
- [ ] Live telemetry dashboard frontend