# Chirpy

A project from boot.dev's go backend course.  
Chirpy is a lightweight social media platform backend written in Go. It provides APIs for user management, authentication, and posting short messages called "chirps." The project uses Postgres for backend.

## Features

- User registration and authentication
- Password hashing using Argon2id
- JWT-based access tokens and refresh tokens
- CRUD operations for "chirps"
- Admin endpoints for metrics and database reset
- Input validation and basic content moderation
- API key-based webhook handling for user upgrades

## Project Structure

```
.
├── assets/                 # Static assets (e.g., logo)
├── internal/               # Internal packages
│   ├── auth/               # Authentication logic
│   └── database/           # Database queries and models
├── sql/                    # SQL schema and queries
│   ├── schema/             # Database schema migrations
│   └── queries/            # SQL queries for sqlc
├── main.go                 # Entry point of the application
├── mux.go                  # HTTP router setup
├── handler_*.go            # HTTP handlers for various endpoints
├── funcs.go                # Utility functions
├── consts.go               # Global constants and configuration
├── index.html              # Welcome page
├── go.mod                  # Go module file
└── sqlc.yaml               # sqlc configuration
```

## Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database
- Environment variables:
  - `DB_URL`: PostgreSQL connection string
  - `SECRET`: for JWT signing
  - `POLKA_KEY`: API key for webhook authentication

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo/chirpy.git
   cd chirpy
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Set up the database:
   - Apply the migrations in `sql/schema/` using a tool like `goose`.

4. Configure environment variables:
   - Create a `.env` file with the required variables:
     ```
     DB_URL=your_database_url  #format:postgresql://USERNAME:PASSWORD@HOST:PORT/DATABASE_NAME
     SECRET=your_secret_key
     POLKA_KEY=your_polka_key
     ```

## Usage

1. Start the server:
   ```sh
   go run main.go
   ```

2. The server will run on `http://localhost:8080`.

3. Use the following endpoints:

### API Endpoints

#### Public Endpoints
- `POST /api/users`: Register a new user
- `POST /api/login`: Log in and receive tokens
- `POST /api/refresh`: Refresh access token

#### Protected Endpoints (Require Bearer Token)
- `POST /api/chirps`: Create a new chirp
- `GET /api/chirps`: Get all chirps(Supported queries: `?sort=asc(default)/desc` which sorts by created_at, or `?author_id=ID` which returns chirps for given ID)
- `GET /api/chirps/{chirpID}`: Get a specific chirp
- `DELETE /api/chirps/{chirpID}`: Delete a chirp
- `PUT /api/users`: Update user details

#### Admin Endpoints
- `GET /admin/metrics`: View metrics
- `POST /admin/reset`: Reset the database (dev mode only)

#### Webhooks
- `POST /api/polka/webhooks`: Handle user upgrade events

## Testing

Run the tests using:
```sh
go test ./...
```

## Acknowledgments

- [sqlc](https://github.com/kyleconroy/sqlc) for generating type-safe database queries
- [Argon2id](https://github.com/alexedwards/argon2id) for secure password hashing
- [JWT](https://github.com/golang-jwt/jwt) for token-based authentication
