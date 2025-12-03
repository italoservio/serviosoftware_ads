# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based HTTP API service for managing ad cloakers. The service uses MongoDB for data persistence, JWT for authentication, and is containerized with Docker. The codebase is written in Portuguese (comments, error messages, logs).

## Development Commands

### Build and Run
```bash
# Build the application
go build -o bin ./cmd/main.go

# Run locally
go run cmd/main.go

# Build Docker image
docker build -t serviosoftware_ads .

# Run Docker container
docker run -p 8080:8080 serviosoftware_ads
```

### Dependencies
```bash
# Install/update dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### Environment Setup
The application requires a `.env` file with:
- `MONGODB_URI` - MongoDB connection string
- `AUTH_SECRET` - JWT secret key for authentication

## Architecture

### Project Structure

The codebase follows a modular, clean architecture pattern:

```
cmd/               - Application entry point
internal/          - Private application code
  api/            - HTTP route registration and handlers
  deps/           - Dependency injection container
  modules/        - Business modules (organized by domain)
    {module}/
      commands/   - Use cases/application logic
      models/     - Domain models
      repos/      - Repository interfaces and implementations
      clients/    - External service clients
pkg/              - Shared packages (can be exported)
  db/            - Database connection
  jwt/           - JWT parsing and middleware
  exception/     - Custom error types
  validation/    - Request validation
  env/           - Environment variable loading
  cast/          - Type conversion utilities
```

### Key Architectural Patterns

#### Dependency Injection Container
The `internal/deps/container.go` file creates all dependencies and wires them together. This is the central place where:
- Database connections are established
- Validators are initialized
- Repositories are instantiated
- Commands are created with their dependencies
- HTTP API handlers are assembled

When adding new modules, update the Container struct and NewContainer function.

#### Command Pattern
Each business operation is implemented as a "Command" in `internal/modules/{module}/commands/{command}/`:
- `{command}.go` - Contains the business logic and `Exec()` method
- `http_api.go` - HTTP handler that validates input and calls the command

#### Repository Pattern
Data access is abstracted through repository interfaces defined in `internal/modules/{module}/repos/`:
- Interface defines the contract (e.g., `CloakerRepository`)
- MongoDB implementation (e.g., `MongoCloakerRepository`)

#### Exception Handling
Custom exception types in `pkg/exception/` implement the `AppException` interface with `WriteJSON()` method:
- All exceptions have `Messages`, `StatusText`, and `StatusCode`
- Use `exception.ToAppException(err)` to convert errors to HTTP responses
- Specific exceptions: `PayloadException`, `ValidatorException`, `UnauthorizedException`, `NotFoundException`, `RepositoryException`, etc.

#### Validation
Uses `go-playground/validator` with custom validators:
- Struct tags define validation rules (e.g., `validate:"required,url"`)
- Custom validator `oneofsortorder` in `pkg/validation/` for sort order validation
- Validation errors are automatically converted to `ValidatorException`

#### JWT Middleware
Protected routes use `jwt.Middleware(env)` which:
- Extracts Bearer token from Authorization header
- Validates token signature using `AUTH_SECRET`
- Checks token expiration
- Adds session to request context
- Returns 401 for invalid/expired tokens

### Module Structure (Cloakers Example)

The `cloakers` module demonstrates the pattern to follow:

1. **Models** ([internal/modules/cloakers/models/cloaker.go](internal/modules/cloakers/models/cloaker.go))
   - Domain entities with BSON and JSON tags
   - Validation tags for struct validation

2. **Repository** ([internal/modules/cloakers/repos/](internal/modules/cloakers/repos/))
   - Interface defining data operations
   - MongoDB implementation

3. **Commands** ([internal/modules/cloakers/commands/](internal/modules/cloakers/commands/))
   - Each command in its own subdirectory
   - Business logic separated from HTTP concerns
   - Input/Output structs with validation tags

4. **HTTP API** ([internal/modules/cloakers/commands/api.go](internal/modules/cloakers/commands/api.go))
   - Aggregates all command HTTP APIs
   - Injected into container

5. **Route Registration** ([internal/api/cloakers.go](internal/api/cloakers.go))
   - Defines HTTP routes and methods
   - Applies middleware (JWT authentication)
   - Maps routes to handler methods

### Adding a New Module

1. Create module directory structure under `internal/modules/{module}/`
2. Define models in `models/`
3. Create repository interface and implementation in `repos/`
4. Implement commands in `commands/{command}/` with business logic and HTTP API
5. Create aggregator in `commands/api.go`
6. Update `internal/deps/container.go` to wire dependencies
7. Add route registration in `internal/api/{module}.go`
8. Register routes in `cmd/main.go`

### HTTP Server Configuration

The server ([cmd/main.go](cmd/main.go)):
- Listens on port 8080
- Uses Gorilla Mux for routing
- CORS enabled for localhost:5173 and users.serviosoftware.com
- Graceful shutdown with 30-second timeout on SIGINT/SIGTERM
- Disconnects from MongoDB on shutdown

## CI/CD

GitHub Actions workflow ([.github/workflows/cicd.yaml](.github/workflows/cicd.yaml)):
- Triggers on push to main
- Builds and pushes Docker image to GHCR
- Generates artifact attestation
- Cleans up old images (keeps last 3 tagged versions)
- Railway deployment is disabled (`if: false`)

## Language & Conventions

- Code comments, error messages, and logs are in Portuguese
- Use `bson` tags for MongoDB field mapping
- Use `json` tags for HTTP response serialization
- Use `validate` tags for input validation
- Error messages start with lowercase (Portuguese convention)
