# companies

This microservice is built in Go and implements the foundations of a `Companies` Service. The design follows a clean, maintainable architecture with clear separation of concerns and dependency injection. Below is an overview of what is implemented so far.

## Service Architecture

### Service Initialization

A dedicated `service` package is responsible for:
* Validating and loading configuration from environment variables.
* Initializing the global structured logger (Zap).
* Constructing the core `Service` object that holds application-wide dependencies.
* Preparing the HTTP server configuration.

This ensures consistent setup logic and centralizes all app-level initialization.

### Graceful Lifecycle Management

#### Graceful Startup

A fully configured HTTP server is created using:
* `gin` for routing and middleware.
* Configurable timeouts for production safety.
* Route registration via the `internal/http/routes` package.

#### Graceful Shutdown

The service listens for termination signals (`SIGINT`, `SIGTERM`) and ensures:
* Clean termination of the HTTP server.
* Proper logging of shutdown events.
* Controlled timeout on shutdown to allow ongoing requests to finish.

All shutdown logic is encapsulated inside the `service.Run` method, keeping `main.go` clean.

## HTTP Layer

### Routing

Routes are defined under `internal/http/routes`.  
Currently, the `/companies` endpoints are registered here. This modular approach allows future expansion with additional routes or services.

### Server Infrastructure

A reusable HTTP server wrapper exists in `internal/http/server.go` that handles:
* Starting the server.
* Observing context cancellation.
* Clean shutdown and logging.

## Logging

A dedicated logger package (`internal/logger`) initializes a structured Zap logger with:
* Production-ready configuration.
* Centralized initialization for all components.
* Injection into the main `Service` object for consistent logging across layers.

## Config Management

Configuration is validated and loaded using `internal/http/config.go`, which provides:
* Strict configuration validation.
* Detailed error propagation.
* Environment variable overrides for flexibility.

## Main Application Entrypoint

The `cmd/main.go` file:
* Creates a cancellable context that responds to OS signals.
* Initializes the service using `service.New`.
* Starts the HTTP server using `service.Run`.
* Logs lifecycle events cleanly.

This keeps the entrypoint simple and follows Go best practices.

## Current Implementation Status

### Core Infrastructure
* Service lifecycle: initialization → run → graceful shutdown.
* Configuration validation and loading.
* Logger setup with structured logging.
* HTTP server with configurable timeouts.
* Signal handling for clean shutdown.

### HTTP Layer
* Route registration for `/companies`.
* Gin router initialization.
* Graceful server shutdown with context handling.

### Architecture Decisions
* Clear separation of concerns: `service`, `http`, `logger`, `config`, `routes`.
* Service ready for future integration with database, Kafka, and business logic.

## Next Steps (Planned / TODO)

* Implement repository layer with Postgres.
* Add JWT authentication middleware for protected endpoints.
* Implement Kafka event publishing on mutating operations.
* Dockerize the application and external services.
* Add health and readiness endpoints for observability.
