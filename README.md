# Censys Take-Home Assessment: Overview of deliverable / testing instructions

### Overview of two main services

**API Gateway (REST/JSON on :8080):** This is the public-facing HTTP interface that clients interact with directly. It translates HTTP requests into gRPC calls and handles JSON serialization/deserialization.

**KV Store (gRPC on :8081):** The internal service that actually stores the data in an in-memory map with thread-safe operations using `sync.RWMutex`.

### A few key implementation details / choices that were made

**Protobuf as the Contract:** I defined the gRPC service contract in `proto/kv/v1/kv.proto` with three operations (Put, Get, Delete). This generated the client/server interfaces for both services, ensuring type safety across the network boundary.

**Test-Driven Development:** I wrote unit tests for the in-memory store with a race detector to verify thread safety. 

**Docker Setup:** Created Dockerfiles for both services using multi-stage builds to keep image sizes small. Used docker-compose to run both services on the same network so they can communicate via service names.

## How to run 

### Quick Start with Docker

```bash
# Start both services
docker compose up -d

# Test it works
curl -X PUT http://localhost:8080/kv/test -H "Content-Type: application/json" -d '{"value": "hello"}'
curl http://localhost:8080/kv/test
curl -X DELETE http://localhost:8080/kv/test

# Stop services
docker compose down
```

### Local Development

If you want to run it locally without Docker:

1. Generate protobuf code:
```bash
export PATH=$PATH:$HOME/go/bin
bash scripts/gen-proto.sh
```

2. Terminal 1 - Start KV Store:
```bash
cd services/kv-store && go run ./cmd/server
```

3. Terminal 2 - Start API Gateway:
```bash
cd services/api-gateway && go run ./cmd/server
```

Then test with the same curl commands.

## API Endpoints

All endpoints follow the pattern `/kv/{key}`:

- **PUT** `/kv/{key}` with JSON body `{"value": "data"}` → 201 Created
- **GET** `/kv/{key}` → 200 OK with `{"value": "data"}` or 404 Not Found
- **DELETE** `/kv/{key}` → 204 No Content or 404 Not Found

## Some things I learned / some key technical decisions that were made

**Why did I use Protobuf?** The requirements specifically asked for gRPC communication between services. gRPC uses Protocol Buffers as its serialization format, so I had to define the service contract in a `.proto` file and generate Go code from it. This created 42 lines of proto that generated 585 lines of Go code automatically.

**Thread safety:** I used `sync.RWMutex` to allow multiple concurrent reads while making sure of exclusive access for writes.

## testing
- Unit tests for the store (includes race detector):
```bash
cd services/kv-store && go test -race ./internal/store
```
  covers put/get/delete, overwrite behavior, and concurrent read/write
- Manual end‑to‑end (shown above): PUT → GET → DELETE → GET (404)

## What could be added / increased scope

This is an MVP built within the ~4 hour recommended time constraint. Here's what I would add if I had more time:

**Persistence:** Currently all data is in-memory and lost on restart. I could add disk-backed storage options (like PostgreSQL or Redis) with configurable backends.

**Logging / Metrics:** Could add structured logging with `log/slog`, some type of metrics software, and distributed tracing for debugging production issues.

**Testing:** Integration tests using testcontainers, chaos testing for resilience, and load testing to understand certain bottlenecks / edge cases.

**Configuration:** Environment-specific configs, feature flags, and shutdown with request draining.

## My project structure

```
.
├── proto/kv/v1/kv.proto           # service contract definition
├── services/
│   ├── kv-store/                  # gRPC service
│   │   ├── internal/proto/kv/v1/  # generated proto stubs (ignored)
│   │   ├── internal/server.go     # gRPC handler
│   │   ├── internal/store/        # storage logic with tests
│   │   └── cmd/server/main.go     # entrypoint
│   └── api-gateway/               # HTTP gateway
│       ├── internal/
│       │   ├── proto/kv/v1/       # generated proto stubs (ignored)
│       │   ├── grpc_client.go     # gRPC client wrapper
│       │   └── handlers.go        # HTTP handlers
│       └── cmd/server/main.go     # entrypoint
├── scripts/gen-proto.sh           # proto → go code generation
└── docker-compose.yml             # service orchestration
```
