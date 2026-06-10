# Jackfruit Orchard Monitor

A full-stack IoT system for monitoring jackfruit orchard conditions. ESP32 firmware polls environmental data from industrial RS485 sensors via Modbus RTU and publishes telemetry to an MQTT broker. A Go web application provides the dashboard and REST API backend.

---

## Repository Structure

```
Jackfruit_Orchard_Monitor/
├── cmd/                            # Go application entry point
├── internal/                       # HTTP handlers, middleware, routes, views, database
├── frontend/                       # SSR static assets
├── nginx/                          # Nginx configuration
├── src/                            # ESP32 firmware source
├── include/
├── lib/
├── test/
├── main/
├── secrets/
├── docker-compose.yml
├── Dockerfile
├── go.mod / go.sum
├── sqlc.yaml
├── Makefile
└── platformio.ini
```

---

## Firmware (ESP32)

Built with PlatformIO and the Arduino framework. Designed for the ESP32 WROOM 32 (38-pin NodeMCU variant).

### Features

- Modbus RTU polling over RS485 (non-blocking loop)
- Secure MQTT client with TLS encryption and JSON payloads
- Connection supervisor for automated WiFi and MQTT reconnection
- Static analysis integration via Clang-Format and Clang-Tidy

### Hardware Wiring

#### ESP32 WROOM 32 (38-pin) → RS485 TTL V3 Module

| ESP32 Pin    | RS485 Module Pin | Wire Color |
|--------------|------------------|------------|
| GPIO17 (TX2) | TXD              | Orange     |
| GPIO16 (RX2) | RXD              | Yellow     |
| 3.3V         | VCC              | Red        |
| GND          | GND              | Black      |

#### RS485 Module → SHT40 RS485 Sensor

| RS485 Module | SHT40 Pin       | Wire Color |
|--------------|-----------------|------------|
| V+           | External 5V–28V | Brown      |
| V-           | System GND      | Black      |
| RS485-A      | A+              | Yellow     |
| RS485-B      | B-              | Blue       |

> [!WARNING]
> The SHT40 RS485 sensor requires at least 5V on V+. Do not power it from the ESP32's 3.3V pin. Swapping A+ and B- will cause a communication timeout error.

### Firmware Commands

```bash
# Build
platformio run --environment esp32dev

# Upload
platformio run --target upload --environment esp32dev

# Lint (Clang-Tidy)
platformio check --environment esp32dev
```

---

## Web Application

A server-side rendered dashboard built with Go, Templ, HTMX, and Tailwind CSS.

### Tech Stack

| Layer       | Technology                    |
|-------------|-------------------------------|
| Backend     | Go 1.25.x + Gin framework     |
| Templates   | Templ (type-safe SSR)         |
| Styling     | Tailwind CSS v4               |
| Interactivity | HTMX                        |
| Database    | PostgreSQL                    |
| Hot Reload  | Air                           |

### Prerequisites

- Go 1.25.x or higher
- Node.js 24.x or higher
- Docker & Docker Compose
- Make

### Quick Start

```bash
# Install dependencies
go mod download

# Setup environment variables
cp .env.example .env
# Edit .env with your configuration

# Run with Docker (recommended)
make docker-watch

# Or run locally with hot reload
make watch
```

Visit http://localhost:8080

### Environment Variables

```env
PORT=8080
APP_ENV=dev
GIN_MODE=debug

POSTGRES_VERSION=18
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DATABASE=goapp
POSTGRES_USERNAME=your_username
POSTGRES_PASSWORD=your_password
POSTGRES_SCHEMA=public
```

### Makefile Commands

#### Build

| Command               | Description                  |
|-----------------------|------------------------------|
| `make all`            | Build and test               |
| `make build`          | Build application binary     |
| `make templ-generate` | Generate templ files         |
| `make sqlc-generate`  | Generate sqlc database files |
| `make tailwind-build` | Build Tailwind CSS           |

#### Development

| Command      | Description                       |
|--------------|-----------------------------------|
| `make watch` | Hot reload with Air (recommended) |
| `make run`   | Run SSR server + SPA frontend     |

#### Docker

| Command                  | Description                           |
|--------------------------|---------------------------------------|
| `make docker-watch`      | Start dev environment with hot reload |
| `make docker-watch-down` | Stop dev environment                  |
| `make docker-prod`       | Start production environment          |
| `make docker-prod-down`  | Stop production environment           |

#### Database Migrations

| Command             | Description              |
|---------------------|--------------------------|
| `make migrate-up`   | Run pending migrations   |
| `make migrate-down` | Roll back last migration |

#### Testing & Code Quality

| Command         | Description              |
|-----------------|--------------------------|
| `make test`     | Run all tests            |
| `make itest`    | Run integration tests    |
| `make lint`     | Run linter               |
| `make lint-fix` | Auto-fix lint issues     |
| `make vet`      | Run static analysis      |
| `make fmt`      | Format code              |
| `make clean`    | Remove binary and generated files |

### API Endpoints

#### Pages (SSR)

| Method | Path            | Description                        |
|--------|-----------------|------------------------------------|
| GET    | `/`             | Home page                          |
| GET    | `/contact`      | Contact page                       |
| POST   | `/contact`      | Submit contact form                |
| GET    | `/login`        | Login page                         |
| POST   | `/login`        | Authenticate user                  |
| GET    | `/register`     | Register page                      |
| POST   | `/register`     | Create new account                 |
| GET    | `/logout`       | Logout and clear session           |
| GET    | `/dashboard`    | Dashboard (requires auth)          |
| GET    | `/sitemap.xml`  | Sitemap                            |
| GET    | `/robots.txt`   | Robots file                        |

#### API

| Method | Path              | Description           |
|--------|-------------------|-----------------------|
| GET    | `/api/`           | API information       |
| GET    | `/api/health`     | Health check          |
| GET    | `/api/websocket`  | WebSocket connection  |

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/name`)
3. Commit changes (`git commit -m 'Add feature'`)
4. Push to branch (`git push origin feature/name`)
5. Open a Pull Request

---

## License

[GNU Affero General Public License v3.0 (AGPL-3.0)](https://www.gnu.org/licenses/agpl-3.0.html) — strong copyleft. Derivative works and network deployments must be open-sourced.
