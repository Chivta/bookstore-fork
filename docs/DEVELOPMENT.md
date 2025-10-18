# Development Guide

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Node.js 18+ (for frontend)
- PostgreSQL 15+ (if running services locally without Docker)

## Project Structure

```
My-Distributed-Bookstore/
├── services/
│   ├── books-service/      # Book catalog management
│   ├── users-service/      # Authentication & user management
│   └── logging-service/    # Centralized logging
├── frontend/
│   └── customer-app/       # React frontend (to be implemented)
├── docker/                 # Docker configuration
├── scripts/                # Utility scripts
├── docs/                   # Documentation
└── docker-compose.yml      # Docker Compose configuration
```

## Quick Start with Docker Compose

### 1. Build and start all services

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build
```

### 2. Verify services are running

```bash
# Check service health
curl http://localhost:8081/health  # Books Service
curl http://localhost:8082/health  # Users Service
curl http://localhost:8084/health  # Logging Service

# Check readiness
curl http://localhost:8081/ready
curl http://localhost:8082/ready
curl http://localhost:8084/ready
```

### 3. View logs

```bash
# View logs for all services
docker-compose logs -f

# View logs for specific service
docker-compose logs -f books-service
docker-compose logs -f users-service
docker-compose logs -f logging-service
```

### 4. Stop services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

## Running Services Locally (without Docker)

### 1. Start PostgreSQL and Redis

```bash
# Using Docker for just the databases
docker-compose up -d postgres redis
```

### 2. Run Books Service

```bash
cd services/books-service

# Download dependencies
go mod download

# Run the service
go run cmd/server/main.go

# Or build and run
go build -o books-service cmd/server/main.go
./books-service
```

### 3. Run Users Service

```bash
cd services/users-service

# Download dependencies
go mod download

# Set JWT secret (required)
export JWT_SECRET="your-secret-key-here"

# Run the service
go run cmd/server/main.go
```

### 4. Run Logging Service

```bash
cd services/logging-service

# Download dependencies
go mod download

# Run the service
go run cmd/server/main.go
```

## API Testing

### Users Service - Authentication

#### Register a new user

```bash
curl -X POST http://localhost:8082/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8082/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

Response includes a JWT token. Save it for authenticated requests.

#### Get user profile (requires authentication)

```bash
curl -X GET http://localhost:8082/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### Books Service

#### Create a book

```bash
curl -X POST http://localhost:8081/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "isbn": "9780134190440",
    "title": "The Go Programming Language",
    "description": "The authoritative resource to writing clear and idiomatic Go",
    "price": 44.99,
    "stock_quantity": 50,
    "language": "en",
    "pages": 400,
    "format": "paperback"
  }'
```

#### List all books

```bash
curl -X GET "http://localhost:8081/api/v1/books?limit=10&offset=0"
```

#### Get a specific book

```bash
curl -X GET http://localhost:8081/api/v1/books/{book-id}
```

#### Update book stock

```bash
curl -X PATCH http://localhost:8081/api/v1/books/{book-id}/stock \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 10
  }'
```

### Logging Service

#### Create a log entry

```bash
curl -X POST http://localhost:8084/api/v1/logs \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "books-service",
    "level": "INFO",
    "message": "Book created successfully",
    "trace_id": "trace-123-456",
    "metadata": "{\"book_id\": \"some-uuid\"}"
  }'
```

#### Query logs

```bash
# Get logs for a specific service
curl "http://localhost:8084/api/v1/logs?service_name=books-service&limit=50"

# Get error logs
curl "http://localhost:8084/api/v1/logs?level=ERROR&limit=50"

# Get logs by trace ID
curl "http://localhost:8084/api/v1/logs?trace_id=trace-123-456"
```

## Development Workflow

### Running Tests

```bash
# Run tests for Books Service
cd services/books-service
go test ./... -v -cover

# Run tests for Users Service
cd services/users-service
go test ./... -v -cover
```

### Code Formatting

```bash
# Format code
cd services/books-service
gofmt -s -w .

cd services/users-service
gofmt -s -w .
```

### Linting (requires golangci-lint)

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
cd services/books-service
golangci-lint run

cd services/users-service
golangci-lint run
```

## Database Management

### Access PostgreSQL

```bash
# Connect to PostgreSQL container
docker exec -it bookstore-postgres psql -U bookstore -d bookstore_books

# Or connect to specific database
docker exec -it bookstore-postgres psql -U bookstore -d bookstore_users
docker exec -it bookstore-postgres psql -U bookstore -d bookstore_logs
```

### Common SQL Commands

```sql
-- List all books
SELECT * FROM books LIMIT 10;

-- List all users
\c bookstore_users
SELECT id, email, full_name, created_at FROM users;

-- List all logs
\c bookstore_logs
SELECT * FROM logs ORDER BY timestamp DESC LIMIT 20;
```

## Troubleshooting

### Service won't start

1. Check if ports are already in use:
   ```bash
   lsof -i :8081  # Books Service
   lsof -i :8082  # Users Service
   lsof -i :8084  # Logging Service
   lsof -i :5432  # PostgreSQL
   ```

2. Check service logs:
   ```bash
   docker-compose logs books-service
   docker-compose logs users-service
   ```

3. Rebuild containers:
   ```bash
   docker-compose down
   docker-compose up --build
   ```

### Database connection issues

1. Ensure PostgreSQL is running:
   ```bash
   docker-compose ps postgres
   ```

2. Check database health:
   ```bash
   docker exec bookstore-postgres pg_isready -U bookstore
   ```

3. Verify databases were created:
   ```bash
   docker exec bookstore-postgres psql -U bookstore -c "\l"
   ```

### Reset everything

```bash
# Stop all services and remove volumes
docker-compose down -v

# Rebuild from scratch
docker-compose up --build
```

## Environment Variables

### Books Service

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: bookstore)
- `DB_PASSWORD` - Database password (default: dev_password)
- `DB_NAME` - Database name (default: bookstore_books)
- `REDIS_URL` - Redis connection URL (default: localhost:6379)
- `PORT` - HTTP port (default: 8081)
- `GRPC_PORT` - gRPC port (default: 9091)

### Users Service

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: bookstore)
- `DB_PASSWORD` - Database password (default: dev_password)
- `DB_NAME` - Database name (default: bookstore_users)
- `JWT_SECRET` - JWT signing secret (required)
- `JWT_EXPIRATION_HOURS` - JWT expiration in hours (default: 24)
- `PORT` - HTTP port (default: 8082)
- `GRPC_PORT` - gRPC port (default: 9092)

### Logging Service

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: bookstore)
- `DB_PASSWORD` - Database password (default: dev_password)
- `DB_NAME` - Database name (default: bookstore_logs)
- `PORT` - HTTP port (default: 8084)
- `GRPC_PORT` - gRPC port (default: 9094)

## Next Steps

1. Implement React frontend (customer-app)
2. Add API Gateway for unified entry point
3. Implement Orders Service
4. Add service-to-service communication via gRPC
5. Implement caching with Redis
6. Add Prometheus metrics
7. Add distributed tracing with OpenTelemetry
8. Deploy to Kubernetes
