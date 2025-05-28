# LecSens Backend

Backend service for LecSens application built with Go and PostgreSQL.

## Prerequisites

Before you begin, ensure you have the following installed:
- Go (version 1.22 or later)
- Docker & Docker Compose
- PostgreSQL (if running without Docker)
- Git

## Getting Started

### Clone Repository
```bash
git clone <repository-url>
cd be-lecsens
```

### Environment Setup
Create .env file in root directory:
```bash
PORT=3160
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=lecsens
JWT_SECRET=your-secret-key
```

### Run with Docker

1. Build and start containers:
    ```bash
    docker-compose up --build
    ```
2. Stop containers:
    ```bash
    docker-compose down
    ```
3. Remove volumes along with containers:
    ```bash
    docker-compose down -v
    ```

### Run Locally (Without Docker)

1. Install dependencies:
    ```bash
    go mod download
    ```
2. Update database connection in `.env`:
    ```bash
    DB_HOST=localhost
    ```
3. Run the application:
    ```bash
    go run main.go 
    ```