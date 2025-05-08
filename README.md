# Sample MCP Server

A sample MCP (Model Context Protocol) server that connects to a PostgreSQL database.

## Overview

It uses Go with GORM for database operations and PostgreSQL for data storage.

## Features

- Account management (create, read, update, delete)
- Category management (create, read, update, delete)
- Transaction tracking (create, read, update, delete)
- Transaction aggregation (sum, count by account)

## Project Structure

- `compose/` - Docker Compose configuration for local development
- `db/` - Database related code
  - `entity/` - Data model definitions
  - `migrations/` - Database migration scripts
  - `repository/` - Data access layer implementations
- `pkg/` - Shared packages and utilities
- `vendor/` - Vendored dependencies

## Data Model

The application uses the following data model:

- **Account**: Represents financial accounts with ID, name, and type
- **Category**: Represents transaction categories with ID, name, and type
- **Transaction**: Represents financial transactions with amount, date, description, and relationships to accounts and categories

## Prerequisites

- Go 1.16 or higher
- Docker and Docker Compose
- PostgreSQL (or use the provided Docker setup)

## Getting Started

### Setup with Docker

1. Clone the repository
2. Start the PostgreSQL database:
   ```
   mage docker:up
   ```
3. Build and run the application:
   ```
   mage build
   ```

### Development Mode

To start the application in development mode:

```
mage dev
```

This will build the application and start the required Docker services.

## Testing

### Running Unit Tests

```
mage test
```

### Running Integration Tests

```
mage integrationtest
```

This will start the required Docker services, wait for PostgreSQL to initialize, and run the integration tests.

## Build Commands

This project uses [Mage](https://magefile.org/) for build automation. Available commands:

- `mage build` - Build the application using GoReleaser
- `mage compilebuild` - Build using standard Go build
- `mage test` - Run unit tests
- `mage integrationtest` - Run integration tests
- `mage lint` - Run linter
- `mage dev` - Build and start Docker services
- `mage docker:up` - Start Docker services
- `mage docker:down` - Stop Docker services
- `mage docker:logs` - Show Docker logs
- `mage docker:restart` - Restart Docker services
- `mage clean` - Clean build artifacts

## Database Configuration

The PostgreSQL database is configured with the following default settings:

- **Host**: localhost (or postgres in Docker network)
- **Port**: 5432
- **User**: jasoet
- **Password**: localhost
- **Database**: mcp_db

## License

[Add license information here]

## Contributing

[Add contribution guidelines here]