# Gophkeeper

> Secure CLI-based secret manager with a REST backend, built in Go with clean architecture, test coverage, and production-oriented tooling.

##  Summary

**Gophkeeper** is a secure secret management system that combines a RESTful backend with a CLI client. It allows users (or applications) to store, retrieve, and manage sensitive data such as passwords, bank card details, and arbitrary binary data.

This project demonstrates:
- layered architecture (domain-driven separation)
- REST API design with proper error handling
- CLI application design using Cobra
- configuration management with Viper
- structured logging with Zap
- PostgreSQL integration with migrations
- testable codebase with mocks and ~70% coverage
- containerized development via Docker Compose

---

##  Features

- Secure storage of multiple secret types:
    - Passwords (with validation)
    - Bank card data (with validation)
    - Text notes
    - Images
    - Arbitrary binary data

- CLI interface for end users:
    - `add [type]` — store a new secret
    - `get --id` — retrieve a secret
    - `delete --id` — delete a secret
    - `list` — list all stored secrets

-  REST API backend:
    - Follows REST conventions
    - Structured JSON responses
    - Proper HTTP status codes
    - Authentication support

-  Developer-friendly:
    - Can be used as a backend service for other applications
    - Clean separation of concerns
    - Extensible architecture

---

##  Architecture

The project follows a **layered architecture** with clear separation of responsibilities:

- **Presentation layer** — HTTP handlers & CLI interface
- **Service layer** — business logic
- **Data layer** — database interaction via repositories
- **Domain layer** — core entities and contracts

This structure ensures:
- testability
- maintainability
- easy extension (e.g., adding caching or new storage backends)

---

##  Tech Stack

- **Language:** Go
- **Router:** Gin
- **Database:** PostgreSQL
- **Logging:** Zap
- **CLI:** Cobra
- **Configuration:** Viper
- **Testing:** Go testing + mocks + test suites
- **Containerization:** Docker Compose

---

##  API Overview

Example endpoint:
```GET /secrets```

## CLI usage
Examples:
```
# Add a password
gophkeeper add password

# Get a secret
gophkeeper get --id <id>

# Delete a secret
gophkeeper delete --id <id>

# List all secrets
gophkeeper list
```

## Project structure
```
/cmd
/agent        # CLI application
/server       # REST server entrypoint

/database
/fixtures
/migrations

/examples

/internal
/agent        # CLI logic
/app          # application setup
/config       # configuration management
/core         # core application wiring
/datasources  # database integrations
/domain       # business entities & interfaces
/logger       # logging setup
/mocks        # generated mocks for testing
/presentation # HTTP handlers
```

## Testing
- ~70% test coverage
- Unit tests for core logic
- Use of mocks for isolation
- Test suites for structured setup

This ensures reliability and makes the system safe to extend.

## Running the project
Using docker-compose

```docker-compose up --build```

This will start:
- PostgreSQL database
- Gophkeeper server


## Configuration

The project uses:

- Viper for configuration
- Environment variables for runtime flexibility

Example:
```DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
# yp_gophkeeper
