# httpprotocol

Collection of small experiments around HTTP-like protocol handling and simple network demos. The goal is to understand how requests, headers, responses, and simple servers work by implementing them from scratch.

## Layout

- `cmd/httpserver/` — minimal HTTP server demo (entry: `main.go`).
- `cmd/tcplistener/` — raw TCP listener example (entry: `main.go`).
- `cmd/udpsender/` — simple UDP sender for testing responses (entry: `main.go`).
- `internal/headers/` — header parsing and utilities.
- `internal/request/` — request parsing logic and tests.
- `internal/response/` — building responses.
- `internal/server/` — core server handling and helpers.
- `tmp/` — scratch files used by demos (e.g. `tcp.txt`).

## Running demos

From the `httpprotocol` directory you can run each demo directly with `go run`:

```bash
cd httpprotocol
go run ./cmd/httpserver
go run ./cmd/tcplistener
go run ./cmd/udpsender
```

Adjust the `main.go` in each `cmd/*` folder for different ports or behaviors.

## Tests

Run package tests with:

```bash
cd httpprotocol
go test ./...
```
