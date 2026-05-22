# Goneon

A clean, modular, and production-ready boilerplate for Go (Golang) applications. It implements a standard layout structure, env-based configurations, modern structured logging, and graceful HTTP shutdown with 100% Go Standard Library dependencies.

## Key Features

- **Standard Go Layout**: Clean separation between application startup (`cmd/`) and core logic (`internal/`).
- **Standard Library Routing**: Uses Go's native, highly-performant routing capabilities (method-based matching, path parameters) introduced in modern Go. No bulky routing packages needed.
- **Graceful Shutdown**: Intercepts terminal interrupt signals (`SIGINT`, `SIGTERM`) to let ongoing requests wrap up before server termination.
- **Structured Logging**: Built-in structured logs with `log/slog`. Switches automatically to JSON format in production and color/text format in development.
- **Environment Configuration**: Robust and simple config parsing using environment variables with smart defaults.
- **Developer Workflow**: Includes a full `Makefile` and testing routines out-of-the-box.

---

## Directory Structure

```text
├── .gitignore          # Git exclusion rules
├── Makefile            # Developer task targets
├── README.md           # Project documentation
├── go.mod              # Go modules declaration
├── cmd/
│   └── goneon/
│       └── main.go     # Application entrypoint & signal capturer
└── internal/
    ├── config/
    │   └── config.go   # Environment variable parser
    ├── logger/
    │   └── logger.go   # slog logger wrapper
    └── server/
        ├── handlers.go # API endpoint controllers
        ├── handlers_test.go # Handler unit tests
        └── server.go   # HTTP server configuration and middleware
```

---

## Getting Started

### Prerequisites

- Go 1.22 or higher installed. (Currently developed on Go 1.24.5)

### Running Locally

To run the web server:

```bash
make run
```

The server will start by default on `http://localhost:8080`.

### Building

To compile the application binary into `bin/`:

```bash
make build
```

This generates the executable `./bin/goneon`.

### Running Tests

To run the suite of unit tests with Go's race detector active:

```bash
make test
```

### Auto-formatting Code

To clean up all formatting imports and syntax spacing:

```bash
make fmt
```

---

## Configuration

Configurations can be customized using environment variables:

| Variable    | Description                                             | Default       |
|-------------|---------------------------------------------------------|---------------|
| `PORT`      | The port the web server listens on.                    | `8080`        |
| `ENV`       | Execution stage (`development`, `production`, `test`).  | `development` |
| `LOG_LEVEL` | Logging threshold (`debug`, `info`, `warn`, `error`).   | `info`        |

Example running on a custom port:

```bash
PORT=9000 make run
```

---

## API Endpoints

### 1. Health Check
Checks if the server is healthy.
- **URL**: `/health`
- **Method**: `GET`
- **Response**: `200 OK`
  ```json
  {
    "status": "ok",
    "time": "2026-05-21T09:15:24+07:00"
  }
  ```

### 2. Greet
Greets the caller.
- **URL**: `/greet`
- **Method**: `GET`
- **Params**:
  - `name` (optional): Name to greet.
- **Response**: `200 OK`
  ```json
  {
    "message": "Hello, World!"
  }
  ```
  Or using `?name=Gopher`:
  ```json
  {
    "message": "Hello, Gopher!"
  }
  ```
