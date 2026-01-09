# Speed Test Server

A lightweight Golang service for self-hosted internet speed testing with MySQL database.

## Features

- **Speed Testing**: Latency, jitter, download, and upload speed measurements
- **IP Geolocation**: IP information with multiple provider fallbacks
- **MySQL Database**: Persistent storage for metrics and test results
- **REST API**: Simple HTTP endpoints for all operations
- **Health Monitoring**: Built-in health checks

## Quick Start

### Prerequisites

- Docker
- Docker Compose

### Setup

1. **Copy environment configuration**:
   ```bash
   cp .env.example .env
   ```

2. **Configure environment variables** (optional):
   Edit `.env` file with your settings. Default values work out of the box.

3. **Start the application**:
   ```bash
   docker-compose up -d
   ```

4. **Check status**:
   ```bash
   docker-compose ps
   docker-compose logs -f
   ```

### Access the Application

- **API Base URL**: http://localhost:8080
- **Health Check**: http://localhost:8080/healthz
- **Version Info**: http://localhost:8080/version

## API Endpoints

- `GET /ping` - Latency measurement
- `GET /download?size=BYTES` - Download speed test
- `POST /upload` - Upload speed test
- `GET /ws` - WebSocket for jitter measurement
- `GET /ip` - IP geolocation information
- `GET /healthz` - Health check
- `GET /version` - Application version
- `GET /config` - Server configuration

## Management Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
docker-compose logs app
docker-compose logs mysql

# Restart services
docker-compose restart

# Rebuild and restart
docker-compose up -d --build

# Check status
docker-compose ps
```

## Database

MySQL runs on port 3306 and data is persisted in a Docker volume.

To access MySQL directly:
```bash
docker-compose exec mysql mysql -u speedtest -pspeedtest speedtest
```

## Configuration

All configuration is done through environment variables in the `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Application port |
| MYSQL_DATABASE | speedtest | Database name |
| MYSQL_USER | speedtest | Database user |
| MYSQL_PASSWORD | speedtest | Database password |
| IPINFO_TOKEN | provided | IP info API token |
| RATE_LIMIT_ENABLED | false | Enable rate limiting (disabled for speed tests) |
| RATE_LIMIT_REQUESTS_PER_MINUTE | 60 | Requests per minute when enabled |

## Development

Build the application locally:
```bash
go build -o speed-test-server ./cmd/speed-test-server
./speed-test-server
```

Run tests:
```bash
go test ./...
```

## License

MIT License © 2025
