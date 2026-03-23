# Running on Windows

## Prerequisites

- Go 1.16+ (`go version` to verify)
- GCC (MinGW-w64 recommended: `choco install mingw`)
  - Required by `go-sqlite3` (CGO)
- protoc (`choco install protoc`)
  - Required for proto code generation (`go generate ./...`)
- golangci-lint
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```
  - Run with `golangci-lint run ./...` in each service directory

## Build & Run

```bash
# Terminal 1 — racing service
cd racing
go build -o racing.exe .
./racing.exe
# gRPC server listening on: localhost:9000

# Terminal 2 — API gateway
cd api
go build -o api.exe .
./api.exe
# API server listening on: localhost:8000
```

## Verify

```bash
curl -X POST http://localhost:8000/v1/list-races \
  -H "Content-Type: application/json" \
  -d '{"filter": {}}'
```
