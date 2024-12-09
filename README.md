# Music Library

[![Go Version](https://img.shields.io/badge/Go-1.23.2-blue)](https://golang.org/)

The Music Library project is the result of a test assignment.

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

3. Create a `.env` file in the root directory

4. Configure the application (see [Configuration](#configuration)).

5. Start the application containers:

   ```bash
   docker compose up -d
   ```

6. To stop the application:

   ```bash
   docker compose down
   ```

The application will be available at `http://<HTTP_ADDR>`.

---

## Configuration

The application loads configuration parameters from the `.env` file, which must be created in the project root directory, or from environment variables. Available configuration parameters are listed below: 
| Variable Name                 | Description                                            | Required |
|-------------------------------|--------------------------------------------------------|----------|
| HTTP_ADDR                      | Service address (default: `":8080"`)                   | **true** |
| DATABASE_DSN                   | Database connection string (default: none)            | **true** |
| DB_RETRY_COUNT                 | Number of retries for DB (default: `3`)               | false    |
| DB_RETRY_BACKOFF               | Time between DB reconnect retries (default: `"1s"`)   | false    |
| MIGRATIONS_PATH                | Path to migrations folder (default: `"/migrations"`)  | false    |
| DETAILS_SERVER_ADDR            | Address of the details server (default: `"http://localhost:8081"`) | **true** |
| DETAILS_SERVER_RETRY_COUNT     | Number of retries for details server (default: `3`)   | false    |
| DETAILS_SERVER_RETRY_BACKOFF   | Time between retries for details server (default: `"1s"`) | false    |
| DETAILS_SERVER_MAX_RETRY_BACKOFF | Max time between retries for details server (default: `"10s"`) | false    |


### Example `.env` file content:

```env
HTTP_ADDR=:8080
DETAILS_SERVER_ADDR=http://localhost:8081
DATABASE_DSN=postgres://user:pass@localhost:5432/dbname?sslmode=disable
```

Make sure all parameters are correctly set before starting the application.

---

## Project Structure

```plaintext
.
├── cmd
│   └── main.go
├── docs
│   └── images
│       ├── c2.png
│       └── c3.png
├── internal
│   ├── apperrors
│   │   └── apperrors.go
│   ├── application
│   │   ├── application.go
│   │   ├── config
│   │   │   └── config.go
│   │   └── service
│   │       └── service.go
│   ├── domain
│   │   ├── interfaces
│   │   │   ├── details_repository.go
│   │   │   ├── music_repository.go
│   │   │   └── music_service.go
│   │   └── models
│   │       ├── filter.go
│   │       ├── library.go
│   │       ├── library_easyjson.go
│   │       ├── lyrics.go
│   │       ├── lyrics_easyjson.go
│   │       ├── pagination.go
│   │       ├── song.go
│   │       ├── song_detail.go
│   │       ├── song_detail_easyjson.go
│   │       ├── song_easyjson.go
│   │       ├── song_title.go
│   │       └── song_title_easyjson.go
│   ├── infrastructure
│   │   ├── http
│   │   │   └── detailsclient
│   │   │       └── detailsclient.go
│   │   ├── metrics
│   │   │   ├── metrics.go
│   │   │   └── metrics_registry.go
│   │   └── persistance
│   │       ├── migrator
│   │       │   └── migrator.go
│   │       └── postgres
│   │           ├── postgres_repository.go
│   │           └── queries.go
│   └── presentation
│       ├── handler
│       │   ├── handler.go
│       │   └── music.go
│       ├── middleware
│       │   ├── logger.go
│       │   ├── memory_writer.go
│       │   ├── metrics.go
│       │   └── middleware.go
│       └── server
│           └── server.go
├── migrations
│   ├── 001_create_songs_table.down.sql
│   └── 001_create_songs_table.up.sql
├── swagger
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── Dockerfile
├── docker-compose.yaml
├── LICENSE
├── Makefile
├── README.md
├── go.mod
└── go.sum
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

## C4 Architecture Visualization

### Level 2: Containers
![C4 Level 2](./docs/images/c2.png)

### Level 3: Components
![C4 Level 3](./docs/images/c3.png)

---

## License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for more details.
