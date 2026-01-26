---
name: lets-go-standard
description: Enforce idiomatic Go code style and best practices from the "Let's Go" series.
---
# Go Development Standards

When generating Go backend code, strictly adhere to these principles:
- **Standard Library First**: Prioritize `net/http` for routing and handlers unless a specific framework (like Chi or Gin) is explicitly requested.
- **Dependency Injection**: Use an `application` struct to hold shared dependencies (loggers, database models, template caches). Avoid global variables at all costs.
- **Structured Error Handling**: Use custom error helpers (e.g., `serverError`, `clientError`, `notFound`) to ensure consistent responses and logging.
- **Internal Model Pattern**: Encapsulate SQL logic within the `internal/models` package. Use receiver functions on model types for DB operations.
- **Environment Parity**: Always read configuration (DSN, Port, etc.) from environment variables or command-line flags, never hardcode them.