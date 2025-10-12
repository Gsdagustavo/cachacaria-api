# Cachacaria Wilbert API

A Go-based REST API for managing a liquor store (Cachacaria Wilbert) with user authentication, product management, and image uploads.

## Prerequisites

- **Docker** and **Docker Compose** (recommended)
- **Go 1.24.6** (if building locally without Docker)
- **MySQL 8** (if running locally without Docker)

## Tech Stack

- **Language:** Go 1.24.6
- **HTTP Router:** gorilla/mux
- **Database:** MySQL 8.0
- **Authentication:** PASETO tokens
- **Configuration:** TOML files
- **Containerization:** Docker & Docker Compose

## Quick Start with Docker Compose

### Development Mode

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd cachacaria-api
   ```

2. **Start the services**
   ```bash
   docker compose --profile dev up --build
   ```

3. **Access the API**
    - API: http://localhost:8080
    - MySQL: localhost:3307 (from host machine)

### Production Mode
```
bash
docker compose --profile prod up --build
```
The database schema and initial data from `database.sql` will be automatically loaded on first run.

## Configuration

The API uses TOML configuration files located in `build/config/`:

- **`dev.toml`** - Development configuration
- **`prod.toml`** - Production configuration

### Configuration Structure
```
toml
symmetric_key = "your-32-character-secret-key"

[Server]
port = 8080
address = "0.0.0.0"
base_url = "http://localhost:8080"

[Database]
port = 3306
driver = "mysql"
host = "mysql"
user = "appuser"
password = "apppass"
name = "cachacadb"
```
The configuration file is selected via the `CONFIG_PATH` environment variable (set in docker-compose.yml).

## Building the API

### Using Docker

Build the Docker image:
```
bash
docker build -t cachacaria-api:latest .
```
Run the container:
```
bash
docker run -p 8080:8080 \
  -e CONFIG_PATH=/app/config/dev.toml \
  -v $(pwd)/build/config:/app/config \
  cachacaria-api:latest
```
### Building Locally (without Docker)

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Build the binary**
   ```bash
   go build -o api-server main.go
   ```

3. **Set the configuration path**
   ```bash
   export CONFIG_PATH=./build/config/dev.toml
   ```

4. **Run the API**
   ```bash
   ./api-server
   ```

## Docker Architecture

The Dockerfile uses a multi-stage build:

1. **Builder stage:** Uses `golang:1.24.6` to compile the Go application
2. **Runtime stage:** Uses `debian:bookworm-slim` for a minimal production image
    - Runs as non-root user (`appuser`)
    - Includes config files from `build/config/`
    - Exposes port 8080

### Volumes

- **`mysql_data`** - Persists MySQL database
- **`products_images`** - Stores uploaded product images
- **`logs`** - Application logs

## API Endpoints

Base URL: `http://localhost:8080/api`

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and receive PASETO token

### Products
- `GET /api/products` - List all products
- `GET /api/products/{id}` - Get product by ID
- `POST /api/products` - Add product with images (multipart/form-data)

### Users
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID

### Health Check
- `GET /health` - Health check endpoint (outside `/api` prefix)

## Database Setup

The MySQL container is configured with:
- **Root Password:** `root`
- **Database:** `cachacadb`
- **User:** `appuser`
- **Password:** `apppass`
- **Port:** 3306 (internal), 3307 (host)

The `database.sql` file is automatically executed on first initialization.

## Stopping the Services
```
bash
docker compose --profile dev down
```
To remove volumes (including database data):
```bash
docker compose --profile dev down -v
```


## Project Structure

```
.
├── main.go                          # Application entry point
├── infrastructure/
│   ├── infrastructure.go            # Infrastructure initialization
│   ├── config/                      # Configuration management
│   ├── datastore/repositories/      # Database repositories
│   ├── modules/                     # HTTP modules/routes
│   └── util/                        # Utilities (crypto, validation)
├── domain/
│   ├── entities/                    # Domain models
│   └── usecases/                    # Business logic
├── build/
│   └── config/
│       ├── dev.toml                 # Development config
│       └── prod.toml                # Production config
├── database.sql                     # Database schema and seeds
├── docker-compose.yml               # Docker Compose configuration
└── Dockerfile                       # Multi-stage Docker build
```

## License

This project, Cachacaria Wilbert API, is licensed under the MIT License — see the LICENSE file for details.