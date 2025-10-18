# My Distributed Bookstore

A production-ready distributed bookstore system built with microservices architecture, following industry best practices and distributed systems principles.

## Overview

This project implements a complete distributed bookstore with:
- **3 Microservices**: Books, Users (Auth), and Logging services
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Modern Tech Stack**: Go, Fiber, PostgreSQL, Redis, Docker
- **Security**: JWT authentication, bcrypt password hashing, RBAC
- **Observability**: Centralized logging, health checks, structured logging
- **Production-ready**: Docker Compose orchestration, graceful shutdown, auto-migrations

Built following principles from "Distributed Systems" by Tanenbaum & van Steen.

## Features

### Books Service (Port 8081)
- Complete book catalog management (CRUD)
- Multi-table relationships (books, authors, publishers, categories)
- Advanced filtering and search
- Stock management
- Pagination support

### Users Service (Port 8082)
- User registration and authentication
- JWT-based auth with token refresh
- Role-based access control (RBAC)
- Secure password hashing (bcrypt)
- User profiles and addresses

### Logging Service (Port 8084)
- Centralized log aggregation
- Distributed tracing support (trace_id, span_id)
- Advanced querying and filtering
- Time-series log storage

## Tech Stack

**Backend:**
- Go 1.21+ with Fiber web framework
- GORM (ORM)
- PostgreSQL 15 (database-per-service pattern)
- Redis 7 (caching)
- JWT authentication
- zerolog (structured logging)

**Infrastructure:**
- Docker & Docker Compose
- PostgreSQL (3 separate databases)
- Redis
- Custom bridge networking

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### Start All Services

```bash
# Using Make
make up-build

# Or using docker compose directly
docker compose up -d --build

# Check service health
make health

# View logs
make logs
```

### Test the APIs

**Register a User:**
```bash
curl -X POST http://localhost:8082/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8082/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Create a Book:**
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

**List Books:**
```bash
curl "http://localhost:8081/api/v1/books?limit=10&offset=0"
```

## Project Structure

```
My-Distributed-Bookstore/
├── services/
│   ├── books-service/          # Book catalog microservice
│   ├── users-service/          # Authentication & user management
│   └── logging-service/        # Centralized logging
├── docker-compose.yml          # Service orchestration
├── Makefile                    # Development commands
├── scripts/                    # Database initialization
└── docs/                       # Documentation
```

## Documentation

- [DEVELOPMENT.md](docs/DEVELOPMENT.md) - Detailed development guide
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Complete project overview
- [CLAUDE.md](CLAUDE.md) - AI assistant guidance

## Available Commands

```bash
make help          # Show all available commands
make up-build      # Build and start all services
make down          # Stop all services
make logs          # View all service logs
make health        # Check service health
make clean         # Clean up everything
```

## Architecture Highlights

### Microservices Pattern
- Each service has its own database (database-per-service)
- Services communicate via REST APIs (gRPC ready)
- Stateless design for horizontal scaling

### Clean Architecture
Each service follows:
1. **Domain Layer** - Business entities
2. **Repository Layer** - Data access abstraction
3. **Service Layer** - Business logic
4. **Handler Layer** - HTTP endpoints
5. **Middleware Layer** - Auth, logging, CORS

### Security
- JWT authentication with expiration
- bcrypt password hashing
- Role-based access control
- Input validation
- SQL injection prevention
- CORS configuration

### Observability
- Centralized logging service
- Structured logging with zerolog
- Health check endpoints
- Distributed tracing support

## Service Ports

- **Books Service**: HTTP 8081, gRPC 9091
- **Users Service**: HTTP 8082, gRPC 9092
- **Logging Service**: HTTP 8084, gRPC 9094
- **PostgreSQL**: 5432
- **Redis**: 6379

## Environment Variables

Each service can be configured via environment variables. See [DEVELOPMENT.md](docs/DEVELOPMENT.md) for details.

## Distributed Systems Principles

This project implements key concepts from Tanenbaum & van Steen:

- **Transparency**: Location-independent service access
- **Scalability**: Stateless services, database per service
- **Fault Tolerance**: Health checks, graceful shutdown
- **Consistency**: Strong consistency for critical operations
- **Security**: Authentication, authorization, encryption
- **Communication**: REST APIs, structured messaging

## What's Next

- [ ] React frontend (TypeScript + shadcn/ui)
- [ ] API Gateway (unified entry point)
- [ ] Orders Service (shopping cart & checkout)
- [ ] gRPC inter-service communication
- [ ] Kubernetes deployment
- [ ] Prometheus + Grafana monitoring
- [ ] Jaeger distributed tracing

## License

Apache License 2.0

## Contributing

Contributions are welcome! Please read the development guide first.

## Author

Built with best practices in distributed systems architecture.
