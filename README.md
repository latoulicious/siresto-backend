# SiResto Backend

A robust restaurant management system backend built with Go, following Clean Architecture principles.

## Tech Stack

- **Language:** Go 1.24.2
- **Framework:** Fiber v2
- **Database:** PostgreSQL with GORM
- **Authentication:** JWT
- **File Storage:** Cloudflare R2 (S3-compatible)
- **Validation:** go-playground/validator
- **Logging:** Logrus
- **Development Tools:** Air (Live Reload)

## Project Structure

```
siresto/
├── cmd/                    # Application entry points
│   ├── hash/              # Password hashing utility
│   ├── keygen/            # Key generation utility
│   ├── migrate/           # Database migration tool
│   └── server/            # Main application server
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── domain/           # Business domain models
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── repository/       # Data access layer
│   ├── routes/           # Route definitions
│   ├── service/          # Business logic layer
│   ├── utils/            # Utility functions
│   └── validator/        # Input validation
├── migrations/           # Database migrations
├── pkg/                  # Public libraries
│   ├── core/            # Core functionality
│   ├── crypto/          # Cryptography utilities
│   ├── db/              # Database utilities
│   ├── dto/             # Data Transfer Objects
│   ├── jwt/             # JWT handling
│   └── logger/          # Logging utilities
└── test/                # Test files
```

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL
- Cloudflare R2 account (or compatible S3 storage)

## Setup and Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/latoulicious/siresto.git
   cd siresto
   ```

2. Copy the environment file and configure your variables:
   ```bash
   cp .env_example .env
   ```

3. Configure the following environment variables in `.env`:
   - `APP_ENV`: Application environment (development/production)
   - `PORT`: Server port (default: 3000)
   - `ALLOWED_ORIGINS`: CORS allowed origins
   - `DATABASE_URL`: PostgreSQL connection string
   - `R2_*`: Cloudflare R2 credentials and configuration
   - `JWT_SECRET_KEY`: Secret key for JWT token generation

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Run database migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

## Running the Application

### Development Mode (with Live Reload)
```bash
air
```

### Production Mode
```bash
go run cmd/server/main.go
```