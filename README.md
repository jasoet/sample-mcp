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
- **Transaction**: Represents financial transactions with amount, date, description, and relationships to accounts and
  categories

## Prerequisites

- Go 1.24 or higher (as specified in go.mod)
- [Mage](https://magefile.org/) build tool
- Docker and Docker Compose
- PostgreSQL (or use the provided Docker setup)

### Installing Mage

To install Mage, run:

```
go install github.com/magefile/mage@vlatest
```

Make sure the Go bin directory is in your PATH.

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

## Configuration

### Configuration File

The application can be configured using a `config.yml` file. The configuration file can be placed in the same directory as the binary or specified using the `MCP_SERVER_CONFIG` environment variable.

A sample configuration file is provided in `config.yml.sample`. You can copy this file to `config.yml` and modify it according to your needs.

```bash
# Copy the sample configuration file
cp config.yml.sample config.yml

# Edit the configuration file
nano config.yml
```

### Environment Variable

You can specify the path to the configuration file using the `MCP_SERVER_CONFIG` environment variable:

```bash
export MCP_SERVER_CONFIG=/path/to/your/config.yml
```

### Database Configuration

The PostgreSQL database is configured with the following default settings:

- **Type**: POSTGRES (also supports MYSQL and MSSQL)
- **Host**: localhost (or postgres in Docker network)
- **Port**: 5432
- **User**: jasoet
- **Password**: localhost
- **Database**: mcp_db
- **Timeout**: 3 seconds
- **Max Idle Connections**: 5
- **Max Open Connections**: 10

## License

This project is licensed under the MIT License - see below for details:

```
MIT License

Copyright (c) 2025 Sample MCP Server

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
