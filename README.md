# Go E-Commerce API

A clean, production-ready REST API for e-commerce built with Go. This project implements user authentication, shopping cart management, and follows best practices for Go web development.

## Features

- **JWT Authentication** - Secure user registration and login with JWT tokens
- **Shopping Cart** - Full CRUD operations for cart items
- **PostgreSQL** - Robust database with migrations
- **Clean Architecture** - Well-organized codebase with separation of concerns
- **Type-Safe Queries** - Using sqlc for compile-time SQL query validation
- **Docker Support** - Easy database setup with Docker Compose
- **Structured Logging** - Request logging and error handling middleware

## Tech Stack

- **Go 1.24+** - Programming language
- **PostgreSQL 16** - Database
- **pgx/v5** - PostgreSQL driver
- **JWT** - Authentication tokens
- **sqlc** - Type-safe SQL code generation
- **Docker Compose** - Database containerization

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose (for database)
- PostgreSQL 16 (if not using Docker)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/angelchiav/go-ecommerce.git
cd go-ecommerce
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory:

```env
DB_URL=postgres://ecommerce:ecommerce@localhost:5432/ecommerce?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ADDR=:8080
```

**Important:** Change `JWT_SECRET` to a strong, random string in production!

### 3. Start the Database

Using Docker Compose (recommended):

```bash
docker-compose up -d
```

This will start a PostgreSQL 16 container with:
- Database: `ecommerce`
- User: `ecommerce`
- Password: `ecommerce`
- Port: `5432`

### 4. Run Database Migrations

The migrations are located in `db/migrations/`. You'll need to run them using your preferred migration tool (e.g., `migrate`, `golang-migrate`, or manually via `psql`).

Example with `golang-migrate`:

```bash
migrate -path db/migrations -database "$DB_URL" up
```

### 5. Generate SQL Code (if needed)

If you've modified SQL queries, regenerate the sqlc code:

```bash
sqlc generate
```

### 6. Run the Application

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Public Endpoints

#### Health Check
```
GET /health
```

#### Register
```
POST /v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Login
```
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}

Response:
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer"
}
```

### Protected Endpoints

All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

#### Get Current User
```
GET /v1/me
```

#### Get Cart
```
GET /v1/cart
```

#### Add Item to Cart
```
POST /v1/cart/items
Content-Type: application/json

{
  "product_id": 1,
  "qty": 2
}
```

#### Update Cart Item Quantity
```
PATCH /v1/cart/items/{id}
Content-Type: application/json

{
  "qty": 5
}
```

#### Remove Item from Cart
```
DELETE /v1/cart/items/{id}
```

## Project Structure

```
go-ecommerce/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── db/
│   ├── migrations/          # Database migration files
│   └── queries/             # SQL queries for sqlc
├── internal/
│   ├── app/                 # Application initialization
│   ├── config/              # Configuration management
│   ├── db/                  # Database connection
│   ├── handlers/            # HTTP request handlers
│   ├── httpx/               # HTTP utilities and middleware
│   ├── service/             # Business logic layer
│   └── sqlc/                # Generated SQL code
├── docker-compose.yml       # Database container setup
├── go.mod                   # Go module dependencies
└── sqlc.yaml               # sqlc configuration
```

## Development

### Running Tests

```bash
go test ./...
```

### Code Generation

After modifying SQL queries in `db/queries/`, regenerate the code:

```bash
sqlc generate
```

### Database Migrations

Migrations are versioned in `db/migrations/`. Always create both `up` and `down` migration files.

## Configuration

The application uses environment variables for configuration:

- `DB_URL` - PostgreSQL connection string (required)
- `JWT_SECRET` - Secret key for JWT token signing (required)
- `ADDR` - Server address (default: `:8080`)

The config package automatically loads a `.env` file from the project root if present.

## Security Considerations

- Passwords are hashed using `golang.org/x/crypto/bcrypt`
- JWT tokens are used for stateless authentication
- SQL injection protection via parameterized queries (sqlc)
- Input validation on all endpoints

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
