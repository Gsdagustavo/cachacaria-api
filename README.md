# Cachacaria Wilbert API

A Go-based REST API for managing users, authentication (JWT), and products (with photo uploads) for Cachacaria Wilbert (liquor store). It exposes HTTP endpoints for user registration/login, browsing products, and serving uploaded product images. A MySQL database is used for persistence. Docker and Docker Compose are provided for local development.

## Tech Stack
- Language: Go (module: `cachacariaapi`)
- Frameworks/Libraries:
  - HTTP router: `gorilla/mux`
  - JWT: `github.com/dgrijalva/jwt-go`
  - MySQL driver: `github.com/go-sql-driver/mysql`
  - Crypto: `golang.org/x/crypto`
- Database: MySQL 8
- Containerization: Docker, Docker Compose

## Project Layout
- `cmd/api/main.go` — API service entry point
- `cmd/client/main.go` — Client entry (currently unused in Compose; see TODO)
- `internal/` — Clean architecture layers
  - `domain/entities` — Domain entities and request/response models
  - `domain/repositories` — Repository interfaces
  - `domain/usecases` — Application use cases
  - `infrastructure/persistence` — MySQL implementations and errors
  - `interfaces/http` — HTTP layer (router, handlers, middleware)
- `database.sql` — Initial schema/data used by MySQL container
- `Dockerfile` — Multi-stage build for API
- `docker-compose.yml` — MySQL + API services and volumes
- `index.html` — Served as simple docs page at `/docs`

## API Overview
Base URL: `http://localhost:8080`

Public endpoints:
- `POST /auth/register` — Register user; returns a JWT (Authorization header + JSON body)
- `POST /auth/login` — Login; returns a JWT (Authorization header + JSON body)
- `GET /products` — List products
- `GET /products/{id}` — Get a product by ID
- `GET /images/{filename}` — Serves uploaded product images
- `GET /docs` — Serves `index.html`; protected by JWT middleware

Product management:
- `POST /products/add` — Multipart form upload to add a product with images
  - Form fields: `name` (string), `description` (string), `price` (float), `stock` (int)
  - Files: `photos` (one or more), allowed types: `image/jpeg` or `image/png`

User management (additional):
- `GET /users` — List users
- `GET /users/{id}` — Get user by ID
- `POST /users/update/{id}` — Update user (method enforced by handler)
- `POST /users/delete/{id}` — Delete user (method enforced by handler)

Note: Some endpoints enforce HTTP methods via middleware; ensure you use the correct method (see code in `interfaces/http/core/utils.go`).

## Requirements
- Go 1.24.x (module uses `go 1.24.4`)
- Docker 24+ and Docker Compose v2
- Make sure ports 8080 (API) and 3307 (host mapped to MySQL 3306) are available

## Configuration (Environment Variables)
The API reads the following variables (see `cmd/api/main.go`, handlers):
- `DB_USER` — MySQL username
- `DB_PASSWORD` — MySQL password
- `DB_HOST` — MySQL host (e.g., `localhost` or service name `mysql` in Compose)
- `DB_PORT` — MySQL port (e.g., `3306`)
- `DB_NAME` — Database name
- `SERVER_PORT` — Port API listens on (default used in Compose: `8080`)
- `JWT_SECRET` — Secret for signing/validating JWTs (used by auth and middleware)
- `BASE_URL` — Base URL to prefix product image URLs in responses (e.g., `http://localhost:8080`)



## Database
When using Docker Compose, `database.sql` is mounted into the MySQL container and executed on first initialization. This file should contain your schema and any seed data.

## Running the Project

### Option A: Docker Compose (recommended)
- Start services:
  - `docker compose up --build`
- API available at: `http://localhost:8080`
- MySQL available on host: `localhost:3307` (container internal port 3306)
- Product images are persisted in the named volume `products_images` and served from `/images/*`.


### Building a container image
- `docker build -t cachacaria-api:local .`
- `docker run --rm -p 8080:8080 --env-file <(printenv | grep -E 'DB_|SERVER_PORT|JWT_SECRET|BASE_URL') cachacaria-api:local`  
  Note: On Windows PowerShell, use `--env` flags or an env file instead of process substitution.



## Known Issues / TODOs
- `cmd/client/main.go` exists but is not documented/integrated. TODO: Document intended use or remove if obsolete.
- Some request method enforcement may not align with conventional REST (e.g., update/delete via POST paths). TODO: Review and standardize HTTP verbs and routes.
- Product delete/update endpoints for products are not implemented; only add/get endpoints exist. TODO: Complete product CRUD as needed.

## Project Structure (abridged)
- `cmd/`
  - `api/main.go`
  - `client/main.go`
- `internal/`
  - `domain/`
    - `entities/` (auth.go, order.go, product.go, review.go, user.go)
    - `repositories/` (product_repository.go, user_repository.go)
    - `usecases/` (product/, user/)
  - `infrastructure/persistence/` (mysql_*_repository.go, errors.go)
  - `interfaces/http/`
    - `core/` (errors.go, utils.go)
    - `handlers/` (authhandler, producthandler, userhandler, router/middleware)
- `database.sql`
- `Dockerfile`
- `docker-compose.yml`
- `index.html`

## License
This project, Cachacaria Wilbert API, is licensed under the MIT License — see the LICENSE file for details.
