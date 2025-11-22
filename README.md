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
* Service designed with extensibility in mind: easily integrates with Postgres, Kafka, and external services.
* JWT-based authentication implemented for protected endpoints.
* Event-driven architecture: Kafka used to publish events on all mutating operations (create, update, delete).
* Repository layer abstracts database interactions, making the service testable and maintainable.
* Business logic encapsulated in the service layer, separate from HTTP handlers.
* Dockerized environment for local development and testing with Postgres and Kafka.

## Helper Scripts and Makefile

The project includes a `Makefile` and several helper scripts to simplify common tasks:

### Makefile

| Command        | Description |
|----------------|-------------|
| `make up`      | Builds and starts the Dockerized environment (Postgres, Kafka, the service) in detached mode. |
| `make down`    | Stops and removes all Docker containers defined in the compose file. |
| `make kafka-topic` | Runs the `create_kafka_topic.sh` script to create the required Kafka topic (`companies-events`). |
| `make kafka-consume` | Runs the `kafka_consume.sh` script to consume and display events from the Kafka topic for debugging or testing. |
| `make test`    | Runs the `tests.sh` script which executes integration and unit tests against the service. |

### Scripts

* `create_kafka_topic.sh`  
  Automatically finds the Kafka container, enters it, and creates the `companies-events` topic if it does not exist. Prints a success message once the topic is created.

* `kafka_consume.sh`  
  Enters the Kafka container and starts a console consumer on the `companies-events` topic, showing published events in real-time.

* `tests.sh`  
  Runs all the integration and unit tests for the service. It can also be adapted to inject JWT tokens, setup test data, or run tests against the Kafka and Postgres containers.

These scripts and Makefile targets simplify setup, testing, and development workflows, allowing you to manage Docker services and Kafka topics without manually entering containers.

## Requirements

To run this project locally, you need the following installed on your system:

* **Docker** – to run the application, Postgres, and Kafka containers.  
  [Get Docker](https://docs.docker.com/get-docker/)

* **Docker Compose** – to orchestrate multi-container services.  
  Usually included with Docker Desktop; verify with `docker-compose --version`.

* **Make** – to run Makefile targets for building, starting, and stopping services.  
  [Install Make](https://www.gnu.org/software/make/)

* **cURL** (or any HTTP client) – to test the REST API endpoints.  
  Available on most Unix systems by default.

* **Go** (for building the service locally without Docker, optional) – minimum version required: 1.20+.  
  [Install Go](https://golang.org/doc/install)

> Note: The Makefile automates most tasks like building containers, creating Kafka topics, and running tests. Using `make up` is sufficient if you have Docker and Docker Compose installed.

## How to Run

This project uses Docker Compose and a Makefile to simplify setup, testing, and running the service locally.

---

### 1. Start the Environment

Build and start all services (Postgres, Kafka, and the application) in detached mode:

```bash
make up
```

This will:
* Build the Docker images.
* Start Postgres, Kafka, and the service container.
* Set up any initial configuration defined in the Docker Compose file.

2. Create Kafka Topic

You need to create the Kafka topic manually, run:

```bash
make kafka-topic
```

This will create the `companies-events` topic inside the Kafka container.

3. Obtain a JWT Token
The service uses JWT authentication for protected endpoints. Login as the admin user:

```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{
  "username": "admin",
  "password": "admin123"
}'
```

You will receive a JSON response with a token:

```json
{"token": "<JWT_TOKEN>"}
```
Use this token for requests that require authentication.

4. Create a Company
```bash
curl -X POST http://localhost:8080/companies \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <JWT_TOKEN>" \
-d '{
  "name": "Acme Corp",
  "description": "A sample company",
  "amount_of_employees": 100,
  "registered": true,
  "type": "Corporations"
}'
```

5. Get a Company

```bash
curl -X GET http://localhost:8080/companies/<COMPANY_ID> \
-H "Content-Type: application/json"
6. Update a Company
bash
Copy code
curl -X PATCH http://localhost:8080/companies/<COMPANY_ID> \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <JWT_TOKEN>" \
-d '{
  "description": "Updated description",
  "amount_of_employees": 120
}'
```

7. Delete a Company
```bash
curl -X DELETE http://localhost:8080/companies/<COMPANY_ID> \
-H "Authorization: Bearer <JWT_TOKEN>"
```

8. Run Tests
Run the integration and unit tests:
```bash
make test
```

9. Stop the Environment
Once finished, stop all services:

```bash
make down
```
