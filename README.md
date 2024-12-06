# Music Library

[![Go Version](https://img.shields.io/badge/Go-1.23.2-blue)](https://golang.org/)

Music Library is a service for managing a music library. The application is written in Go and follows modern design and code structuring practices.

This project is the result of a test assignment.

GitHub Repository: [Music Library on GitHub](https://github.com/FlutterDizaster/music-library)

---

## Features

- Adding information about music tracks.
- Updating information about existing tracks.
- Deleting tracks from the library.
- Retrieving song lyrics with support for pagination by verses.
- Retrieving library content with support for filtering across all fields and pagination.
- Working with data through a REST API (API documentation is available in Swagger).
- Integration with Prometheus for collecting and storing application metrics.

---

## Installation and Setup

1. Make sure you have the following installed:
   - [Go](https://golang.org/) version **1.23.2** or later.
   - [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).

2. Clone the repository:

   ```bash
   git clone https://github.com/FlutterDizaster/music-library
   cd music-library
   ```

3. Configure the application (see [Configuration](#configuration)).

4. Start the application containers:

   ```bash
   docker compose up -d
   ```

5. To stop the application:

   ```bash
   docker compose down
   ```

The application will be available at `http://<HTTP_ADDR>`.

---

## Configuration

The application loads configuration parameters from a `.env` file, which should be created in the project's root directory. Example `.env` file content:

```env
HTTP_ADDR=localhost:8080            # Server address and port
DETAILS_ADDR=localhost:9090         # Address of the details server
POSTGRES_DSN=postgres://user:pass@localhost:5432/dbname?sslmode=disable
```

Make sure all parameters are correctly set before starting the application.

---

## Project Structure

```plaintext
music-library/
├── cmd/               # Entry points (main files)
│   └── main.go        # Server-side application
├── internal/          # Application logic
│   ├── abstraction/   # Abstractions and interfaces
│   ├── apperrors/     # Custom errors
│   ├── application/   # Application bootstrap logic
│   ├── controller/    # REST API controllers
│   ├── metrics/       # Application metrics
│   ├── models/        # Data entities
│   ├── params/        # Request parameters
│   └── server/        # Server configuration
├── swagger/           # Swagger files for API documentation
├── Dockerfile         # Docker container for the application
├── docker-compose.yaml # Docker Compose configuration
├── Makefile           # Task automation
└── README.md          # Project documentation
```

---

## Usage

### Build and Run
- **Build the server**:

   ```bash
   make build
   ```

- **Start the containers**:

   ```bash
   docker compose up -d
   ```

- **Stop the containers**:

   ```bash
   docker compose down
   ```

---

## API

API documentation is available in the following files:
- [swagger.json](./swagger/swagger.json)
- [swagger.yaml](./swagger/swagger.yaml)

To view the documentation in a convenient interface, you can use [Swagger UI](https://swagger.io/tools/swagger-ui/) or any other compatible tool.

---

## Architecture Visualization

The project is described using the **C4** model:

### Level 2: Containers
![C4 Level 2](./docs/images/c2.png)

### Level 3: Components
![C4 Level 3](./docs/images/c3.png)

---

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for more details.