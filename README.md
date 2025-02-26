## A market REST API for entering NIC at IITU.

A simple, minimalistic REST API built in Go for managing market data. <br />
The project includes an integrated Swagger UI for interactive API documentation.

### Installation & Run

1. Clone the repository

```bash
git clone https://github.com/DaniilKalts/market-rest-api.git
cd market-rest-api
```

2. Install dependencies

```bash
go mod tidy
```

3. Create .env (environment variables) file and fill it with your values

```bash
PORT=8080
DATABASE_DSN="host=localhost user=youruser password=yourpassword dbname=market port=5432 sslmode=disable TimeZone=UTC"
```

4. Run the project

```bash
make run
```

### Swagger UI

Access the interactive API documentation at:

```bash
http://localhost:8080/swagger/index.html
```

![Swagger UI](https://github.com/user-attachments/assets/66b83921-cb30-4d40-a9ed-7fa83626dbc0)

### Project Structure

```bash
market-rest-api
├── app
│   └── app.go           // App initialization & routing
├── docs
│   ├── docs.go          // Generated Swagger metadata
│   ├── swagger.json     // Swagger spec in JSON
│   └── swagger.yaml     // Swagger spec in YAML
├── go.mod               // Module dependencies
├── go.sum               // Module checksums
├── handlers
│   ├── item_handler.go  // HTTP handlers for items
│   └── user_handler.go  // HTTP handlers for users
├── logger
│   └── logger.go        // Logging utilities
├── main.go              // Application entry point
├── Makefile             // Build/run/clean commands
├── models
│   ├── item.go          // Item domain model & Swagger annotations
│   └── user.go          // User domain model & Swagger annotations
├── README.md            // Project overview
├── repositories
│   ├── item_repository.go // Data access for items
│   └── user_repository.go // Data access for users
└── services
    ├── item_service.go  // Business logic for items
    └── user_service.go  // Business logic for users
```
