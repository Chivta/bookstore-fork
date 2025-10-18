# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A production-ready distributed bookstore system built with microservices architecture, following industry best practices and principles from "Distributed Systems" by Tanenbaum & van Steen.

### Tech Stack

**Backend Microservices**:
- Go 1.21+ with Fiber (web framework)
- GORM (ORM)
- JWT authentication
- PostgreSQL (primary database)
- gRPC for inter-service communication
- REST APIs for client communication

**Frontend**:
- TypeScript + React 18+
- shadcn/ui component library
- TanStack Query (data fetching)
- Zustand (state management)
- React Router v6

**Infrastructure**:
- Docker & Docker Compose (local development)
- Kubernetes (orchestration)
- PostgreSQL (relational data)
- Redis (caching & sessions)

**Observability**:
- Structured logging (zerolog)
- Distributed tracing (OpenTelemetry)
- Metrics (Prometheus)

## Project Structure

```
bookstore/
├── services/
│   ├── books-service/          # Book catalog management
│   ├── orders-service/         # Order processing
│   ├── users-service/          # Authentication & user management
│   ├── logging-service/        # Centralized logging microservice
│   └── shared/                 # Shared libraries & protobuf definitions
├── frontend/
│   ├── customer-app/           # Customer-facing React app
│   └── admin-app/              # Admin dashboard (optional)
├── k8s/                        # Kubernetes manifests
│   ├── base/                   # Base configurations
│   └── overlays/               # Environment-specific configs
├── docker/                     # Dockerfiles and compose files
├── docs/
│   ├── architecture.md         # System architecture documentation
│   ├── aws-deployment.md       # AWS deployment guide
│   └── api-specs/              # OpenAPI/Swagger specs
└── scripts/                    # Development & deployment scripts
```

## Core Microservices

### 1. Books Service
**Responsibilities**: Book catalog, inventory, search, categories
**Database**: PostgreSQL (books, categories, publishers, authors)
**Port**: 8081 (HTTP), 9091 (gRPC)

**Key Entities**:
- `books` (id, isbn, title, description, price, stock_quantity, publisher_id)
- `authors` (id, name, bio)
- `book_authors` (book_id, author_id) - many-to-many
- `categories` (id, name, parent_id) - hierarchical
- `book_categories` (book_id, category_id)
- `publishers` (id, name, country)

**Endpoints**:
- `GET /api/v1/books` - List books (with pagination, filters)
- `GET /api/v1/books/:id` - Get book details
- `POST /api/v1/books` - Create book (admin)
- `PUT /api/v1/books/:id` - Update book (admin)
- `GET /api/v1/categories` - List categories

### 2. Users Service
**Responsibilities**: Authentication, authorization, user profiles
**Database**: PostgreSQL (users, roles, sessions)
**Port**: 8082 (HTTP), 9092 (gRPC)

**Key Entities**:
- `users` (id, email, password_hash, full_name, created_at, updated_at)
- `roles` (id, name, permissions)
- `user_roles` (user_id, role_id)
- `sessions` (id, user_id, token_hash, expires_at)
- `addresses` (id, user_id, street, city, state, postal_code, country, is_default)

**Endpoints**:
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - Login (returns JWT)
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update profile
- `GET /api/v1/users/me/addresses` - List addresses

### 3. Orders Service
**Responsibilities**: Shopping cart, order processing, payment coordination
**Database**: PostgreSQL (orders, order_items, payments)
**Port**: 8083 (HTTP), 9093 (gRPC)

**Key Entities**:
- `orders` (id, user_id, status, total_amount, shipping_address_id, created_at)
- `order_items` (id, order_id, book_id, quantity, unit_price, subtotal)
- `order_status_history` (id, order_id, status, changed_at, notes)
- `payments` (id, order_id, amount, method, status, transaction_id, processed_at)
- `shopping_carts` (id, user_id, created_at, updated_at)
- `cart_items` (id, cart_id, book_id, quantity)

**Endpoints**:
- `POST /api/v1/cart/items` - Add to cart
- `GET /api/v1/cart` - Get cart
- `DELETE /api/v1/cart/items/:id` - Remove from cart
- `POST /api/v1/orders` - Create order from cart
- `GET /api/v1/orders` - List user orders
- `GET /api/v1/orders/:id` - Get order details

### 4. Logging Service
**Responsibilities**: Centralized log aggregation, distributed tracing, audit trails
**Database**: PostgreSQL (structured logs) + optional time-series DB
**Port**: 8084 (HTTP), 9094 (gRPC)

**Key Entities**:
- `logs` (id, service_name, level, message, trace_id, span_id, user_id, metadata, timestamp)
- `traces` (trace_id, service_name, operation, duration_ms, status, parent_span_id)
- `audit_logs` (id, user_id, action, resource_type, resource_id, changes, ip_address, timestamp)

**Features**:
- Receive logs from all services via gRPC
- Structured logging with correlation IDs
- Query interface for debugging
- Real-time log streaming
- Retention policies

## Database Design Principles

### Normalization
- All tables in 3NF (Third Normal Form)
- Proper foreign key constraints with ON DELETE CASCADE/SET NULL
- Composite indexes on frequently queried columns
- JSONB columns for flexible metadata (e.g., book metadata, order notes)

### Example Schema (Books Service)

```sql
CREATE TABLE publishers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100),
    website VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    birth_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    parent_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_parent ON categories(parent_id);

CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    isbn VARCHAR(13) UNIQUE NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    publisher_id UUID REFERENCES publishers(id) ON DELETE SET NULL,
    publication_date DATE,
    language VARCHAR(10) DEFAULT 'en',
    pages INTEGER,
    format VARCHAR(50), -- hardcover, paperback, ebook
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    stock_quantity INTEGER NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    cover_image_url TEXT,
    metadata JSONB, -- flexible additional data
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_books_isbn ON books(isbn);
CREATE INDEX idx_books_publisher ON books(publisher_id);
CREATE INDEX idx_books_price ON books(price);

CREATE TABLE book_authors (
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    author_id UUID REFERENCES authors(id) ON DELETE CASCADE,
    author_order INTEGER DEFAULT 1,
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE book_categories (
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, category_id)
);
```

## Go Service Architecture

Each service follows Clean Architecture / Hexagonal Architecture:

```
books-service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/                  # Domain entities
│   │   ├── book.go
│   │   ├── author.go
│   │   └── category.go
│   ├── repository/              # Data access layer (GORM)
│   │   ├── book_repository.go
│   │   └── postgres/
│   │       └── book_repo_impl.go
│   ├── service/                 # Business logic
│   │   └── book_service.go
│   ├── handler/                 # HTTP handlers (Fiber)
│   │   └── book_handler.go
│   ├── middleware/              # Auth, logging, CORS
│   │   ├── auth.go
│   │   └── logger.go
│   └── config/                  # Configuration
│       └── config.go
├── pkg/                         # Public libraries
│   ├── jwt/
│   └── validator/
├── proto/                       # gRPC definitions
│   └── books.proto
├── migrations/                  # SQL migrations
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── tests/
│   ├── integration/
│   └── unit/
├── Dockerfile
├── go.mod
└── go.sum
```

### Design Patterns Used

1. **Repository Pattern**: Abstract data access
2. **Dependency Injection**: Pass dependencies via constructors
3. **Factory Pattern**: Create complex objects (e.g., service initialization)
4. **Strategy Pattern**: Multiple payment/shipping strategies
5. **Observer Pattern**: Event-driven communication between services
6. **Circuit Breaker**: Fault tolerance for inter-service calls
7. **Saga Pattern**: Distributed transactions (order creation)

### Code Quality Standards

- Use `golangci-lint` with strict rules
- 80%+ test coverage minimum
- Context propagation for all operations
- Structured logging with zerolog
- Error wrapping with `fmt.Errorf` or `errors.Wrap`
- Graceful shutdown handling
- Health check endpoints: `GET /health` and `GET /ready`

## Frontend Architecture

### Customer App Structure

```
customer-app/
├── src/
│   ├── components/
│   │   ├── ui/                  # shadcn/ui components
│   │   ├── layout/              # Layout components
│   │   ├── books/               # Book-related components
│   │   │   ├── BookCard.tsx
│   │   │   ├── BookDetails.tsx
│   │   │   └── BookList.tsx
│   │   ├── cart/
│   │   │   ├── CartItem.tsx
│   │   │   └── CartSummary.tsx
│   │   └── orders/
│   ├── pages/
│   │   ├── HomePage.tsx
│   │   ├── BookDetailsPage.tsx
│   │   ├── CartPage.tsx
│   │   ├── CheckoutPage.tsx
│   │   └── OrderHistoryPage.tsx
│   ├── api/                     # API client with axios
│   │   ├── books.ts
│   │   ├── orders.ts
│   │   └── auth.ts
│   ├── hooks/                   # Custom React hooks
│   │   ├── useAuth.ts
│   │   ├── useCart.ts
│   │   └── useBooks.ts
│   ├── store/                   # Zustand stores
│   │   ├── authStore.ts
│   │   └── cartStore.ts
│   ├── types/                   # TypeScript types
│   │   ├── book.ts
│   │   ├── order.ts
│   │   └── user.ts
│   ├── utils/
│   └── App.tsx
├── public/
├── package.json
└── tsconfig.json
```

### Key Features

- **Authentication**: JWT stored in httpOnly cookies
- **Shopping Cart**: Persistent cart (synced with backend)
- **Search & Filters**: Debounced search, category filters, price range
- **Pagination**: Infinite scroll or traditional pagination
- **Responsive Design**: Mobile-first approach with Tailwind CSS
- **Error Handling**: Toast notifications with shadcn/ui
- **Loading States**: Skeleton screens and spinners

## Docker Setup

### Development (docker-compose.yml)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: bookstore
      POSTGRES_PASSWORD: dev_password
      POSTGRES_DB: bookstore
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  books-service:
    build:
      context: ./services/books-service
      dockerfile: Dockerfile
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: bookstore
      DB_PASSWORD: dev_password
      DB_NAME: bookstore_books
      REDIS_URL: redis:6379
    ports:
      - "8081:8081"
      - "9091:9091"
    depends_on:
      - postgres
      - redis

  users-service:
    build:
      context: ./services/users-service
    environment:
      DB_HOST: postgres
      JWT_SECRET: dev_jwt_secret_change_in_production
    ports:
      - "8082:8082"
    depends_on:
      - postgres

  orders-service:
    build:
      context: ./services/orders-service
    ports:
      - "8083:8083"
    depends_on:
      - postgres
      - books-service

  logging-service:
    build:
      context: ./services/logging-service
    ports:
      - "8084:8084"
      - "9094:9094"
    depends_on:
      - postgres

  frontend:
    build:
      context: ./frontend/customer-app
    ports:
      - "3000:3000"
    environment:
      REACT_APP_API_GATEWAY_URL: http://localhost:8080

volumes:
  postgres_data:
```

### Production Dockerfile Example (Go Service)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8081 9091

CMD ["./main"]
```

## Kubernetes Configuration

### Service Deployment Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: books-service
  namespace: bookstore
spec:
  replicas: 3
  selector:
    matchLabels:
      app: books-service
  template:
    metadata:
      labels:
        app: books-service
    spec:
      containers:
      - name: books-service
        image: bookstore/books-service:latest
        ports:
        - containerPort: 8081
          name: http
        - containerPort: 9091
          name: grpc
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: bookstore-config
              key: db_host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: bookstore-secrets
              key: db_password
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: books-service
  namespace: bookstore
spec:
  type: ClusterIP
  selector:
    app: books-service
  ports:
  - name: http
    port: 8081
    targetPort: 8081
  - name: grpc
    port: 9091
    targetPort: 9091
```

## Distributed Systems Principles (Tanenbaum & van Steen)

### Communication Patterns

1. **Synchronous Communication**: REST APIs for client-server
2. **Asynchronous Communication**: Message queues for event-driven (future: RabbitMQ/Kafka)
3. **RPC**: gRPC for inter-service communication
4. **Request-Reply**: HTTP/REST with timeout handling

### Consistency & Replication

- **Database Replication**: PostgreSQL primary-replica setup
- **Eventual Consistency**: Cart synchronization uses eventual consistency
- **Strong Consistency**: Orders and payments use strong consistency (transactions)

### Fault Tolerance

- **Circuit Breaker**: Prevent cascading failures (using `gobreaker`)
- **Retry Logic**: Exponential backoff for transient failures
- **Timeouts**: All external calls have timeouts
- **Health Checks**: Kubernetes probes for automatic recovery
- **Graceful Degradation**: Services can operate with reduced functionality

### Naming & Discovery

- **Kubernetes DNS**: Service discovery via DNS (e.g., `books-service.bookstore.svc.cluster.local`)
- **Environment Variables**: Service endpoints configured via ConfigMaps

### Security

- **Authentication**: JWT tokens with refresh mechanism
- **Authorization**: Role-based access control (RBAC)
- **Secrets Management**: Kubernetes Secrets (migrate to AWS Secrets Manager)
- **TLS**: HTTPS for all external communication
- **Input Validation**: Strict validation on all endpoints

### Scalability

- **Horizontal Scaling**: Kubernetes HPA (Horizontal Pod Autoscaler)
- **Database Connection Pooling**: GORM with connection pooling
- **Caching**: Redis for frequently accessed data
- **CDN**: Static assets served via CDN (future)

## Development Workflow

### Initial Setup

```bash
# Clone repository
git clone <repo-url>
cd bookstore

# Start infrastructure
docker-compose up -d postgres redis

# Run migrations for each service
cd services/books-service
make migrate-up

# Start services
make run

# Start frontend
cd frontend/customer-app
npm install
npm start
```

### Common Commands

**Backend (Go)**:
```bash
# Run service
go run cmd/server/main.go

# Run tests
go test ./... -v -cover

# Generate mocks
mockgen -source=internal/repository/book_repository.go -destination=internal/repository/mocks/mock_book_repository.go

# Lint
golangci-lint run

# Format
gofmt -s -w .

# Migrations
migrate -path migrations -database "postgresql://user:pass@localhost:5432/db?sslmode=disable" up
```

**Frontend (React)**:
```bash
# Install dependencies
npm install

# Development server
npm start

# Build for production
npm run build

# Run tests
npm test

# Type check
npm run type-check

# Lint
npm run lint
```

**Docker**:
```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f books-service

# Stop all services
docker-compose down

# Rebuild specific service
docker-compose up -d --build books-service
```

**Kubernetes**:
```bash
# Apply manifests
kubectl apply -f k8s/base/

# Check pods
kubectl get pods -n bookstore

# View logs
kubectl logs -f deployment/books-service -n bookstore

# Port forward for local testing
kubectl port-forward svc/books-service 8081:8081 -n bookstore

# Scale deployment
kubectl scale deployment/books-service --replicas=5 -n bookstore

# Delete resources
kubectl delete -f k8s/base/
```

## Testing Strategy

### Unit Tests
- Repository layer: Mock database with testcontainers
- Service layer: Mock repositories
- Handler layer: Mock services with httptest

### Integration Tests
- Docker Compose for test environment
- Test entire request flow
- Database cleanup between tests

### E2E Tests
- Cypress for frontend
- Postman/Newman for API tests
- Test critical user journeys

### Load Tests
- k6 or Apache JMeter
- Test under realistic load
- Identify bottlenecks

## Monitoring & Observability

### Metrics (Prometheus)
- HTTP request duration histograms
- Active database connections
- gRPC call success/failure rates
- Custom business metrics (orders/minute, revenue)

### Logging (Structured)
- JSON format with correlation IDs
- Log levels: DEBUG, INFO, WARN, ERROR
- Centralized in logging service
- Query interface for debugging

### Tracing (OpenTelemetry)
- Distributed tracing across services
- Trace propagation via headers
- Visualize with Jaeger (future)

### Dashboards
- Grafana for metrics visualization
- Service health overview
- Business KPIs dashboard

## AWS Deployment Architecture

See `docs/aws-deployment.md` for detailed architecture including:

- **VPC Design**: Multi-AZ setup with public/private subnets
- **EKS Cluster**: Kubernetes on AWS
- **RDS PostgreSQL**: Managed database with read replicas
- **ElastiCache Redis**: Managed caching layer
- **Application Load Balancer**: Traffic distribution
- **S3 Buckets**: Static assets and backups
- **CloudWatch**: Centralized logging and monitoring
- **Secrets Manager**: Secure credential storage
- **Route53**: DNS management
- **CloudFront**: CDN for frontend assets
- **ECR**: Docker image registry

**Cost Optimization**:
- Auto Scaling Groups for compute
- Spot Instances for non-critical workloads
- S3 lifecycle policies
- Reserved Instances for stable workloads

## Future Enhancements

### Phase 2 (Post-MVP)
- [ ] Payment gateway integration (Stripe)
- [ ] Email service (SendGrid/SES)
- [ ] Notification service (push notifications)
- [ ] Review & rating system
- [ ] Recommendation engine (ML-based)
- [ ] Wishlist functionality

### Phase 3 (Advanced)
- [ ] GraphQL API gateway
- [ ] Event sourcing with Kafka
- [ ] CQRS pattern for reads/writes
- [ ] Multi-region deployment
- [ ] A/B testing framework
- [ ] Advanced analytics service

## Security Best Practices

- [ ] OWASP Top 10 mitigation
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS protection (CSP headers)
- [ ] CSRF tokens for state-changing operations
- [ ] Rate limiting on all endpoints
- [ ] DDoS protection (CloudFlare/AWS Shield)
- [ ] Regular dependency updates
- [ ] Secrets rotation policy
- [ ] Security audit logs
- [ ] Penetration testing before production

## Documentation Requirements

- [ ] API documentation (OpenAPI/Swagger)
- [ ] Architecture Decision Records (ADRs)
- [ ] Runbooks for common operations
- [ ] Incident response procedures
- [ ] Database schema diagrams
- [ ] Sequence diagrams for critical flows
- [ ] Developer onboarding guide
- [ ] Code style guide

## Performance Targets

- **API Response Time**: p95 < 200ms, p99 < 500ms
- **Database Queries**: < 50ms for simple queries
- **Frontend Load Time**: First Contentful Paint < 1.5s
- **Uptime**: 99.9% (three nines)
- **Concurrent Users**: Support 10,000+ simultaneous users
- **Order Processing**: < 2s end-to-end

## Contribution Guidelines

- Follow Go Code Review Comments
- Use conventional commits (feat:, fix:, docs:, etc.)
- PR template with checklist
- Code review required before merge
- CI/CD checks must pass
- Update documentation with code changes
- Write tests for new features

## Troubleshooting

### Common Issues

**Service can't connect to database**:
- Check `DB_HOST` environment variable
- Verify PostgreSQL is running: `docker ps`
- Check network connectivity: `docker network inspect bookstore_default`

**JWT token invalid**:
- Verify `JWT_SECRET` matches across services
- Check token expiration
- Ensure clock sync across services (NTP)

**Service crash on startup**:
- Check logs: `kubectl logs <pod-name>`
- Verify database migrations ran successfully
- Check required environment variables are set

**High response times**:
- Check database connection pool exhaustion
- Look for N+1 query problems
- Verify indexes on frequently queried columns
- Check Redis cache hit rates

## Notes for Claude Code

- **Go Version**: Use Go 1.21 or later for all services
- **Node Version**: Use Node 18+ LTS for frontend
- **Package Management**: Go modules (`go.mod`), npm for frontend
- **Environment Variables**: Never commit `.env` files
- **Database Migrations**: Use `golang-migrate` for version control
- **Code Generation**: Use `protoc` for gRPC, `sqlc` for type-safe SQL (optional)
- **Naming Conventions**: 
  - Go: camelCase for private, PascalCase for public
  - TypeScript: camelCase for variables/functions, PascalCase for types/components
  - Database: snake_case for tables and columns
  - URLs: kebab-case for routes

## Quick Reference

**Project Ports**:
- Books Service: 8081 (HTTP), 9091 (gRPC)
- Users Service: 8082 (HTTP), 9092 (gRPC)
- Orders Service: 8083 (HTTP), 9093 (gRPC)
- Logging Service: 8084 (HTTP), 9094 (gRPC)
- Frontend: 3000
- PostgreSQL: 5432
- Redis: 6379

**Essential URLs**:
- Frontend: http://localhost:3000
- API Gateway: http://localhost:8080 (future)
- Books API: http://localhost:8081/api/v1
- Grafana: http://localhost:3001 (future)
- Prometheus: http://localhost:9090 (future)
