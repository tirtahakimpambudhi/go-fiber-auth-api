# Auth API Service

A robust authentication and authorization service built with Go Fiber, PostgreSQL, and Redis.

## 🚀 Tech Stack

- **Framework:** [Go Fiber](https://gofiber.io/) - Express-inspired web framework
- **Database:** PostgreSQL - Primary data storage
- **Cache:** Redis - Session and token management
- **Migration:** [golang-migrate](https://github.com/golang-migrate/migrate) - Database versioning
- **Environment:** Docker & Docker Compose
- **Deployment:** Railway

## 📋 Prerequisites

Before running this project, make sure you have the following installed:

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Redis 6 or higher
- Docker & Docker Compose (optional)
- Make (for running helper commands)

## 🛠️ Installation

1. Clone the repository
```bash
git clone https://github.com/tirtahakimpambudhi/go-fiber-auth-api
```

2. Copy environment example and configure your variables
```bash
cp .env.example .env
```

3. Install dependencies
```bash
go mod download
```

4. Run database migrations
```bash
make db_up
```

5. Start the server
```bash
make build && ./main
```

## 🔑 Environment Variables

Key environment variables required:

```env
# RedisConfig
CACHE_DB_NAME=
CACHE_DB_HOST=
CACHE_DB_PORT=
CACHE_DB_USER=
CACHE_DB_PASS=
CACHE_DB_MAX_CON=100
CACHE_DB_MIN_CON=10
CACHE_DB_MAX_TIME=10
CACHE_DB_MIN_TIME=2

# Casbin
MODEL_PATH=resource/model

# SSLConfig
FIBER_SSL_PATH=resource/ssl
FIBER_SSL_CERT=certificate.crt
FIBER_SSL_KEY=private.key

# FiberConfig
FIBER_HOST=localhost
FIBER_PORT=8081
FIBER_PREFORK=true
FIBER_STRICT_ROUTING=true
FIBER_CASE_SENSITIVE=true
FIBER_BODY_LIMIT=4
FIBER_READ_TIMEOUT=4
FIBER_WRITE_TIMEOUT=5
FIBER_REDUCE_MEMU=true
FIBER_JSON=json

# Bcrypt
HASH_SALT=10

# LoggerConfig
LOG_PATH=resource/logs
LOG_MAX_SIZE=10
LOG_MAX_BACKUP=5
LOG_MAX_SIZE_ROTATE=20
LOG_TIME_FORMAT=2006-01-02
LOG_COLOR_OUTPUT=false
LOG_QUOTE_STR=false
LOG_END_WITH_MESSAGE=false

# SqlConfig
DB_DRIVER=
DB_PROTOCOL=
DB_NAME=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_MAX_CON=100
DB_MIN_CON=10
DB_MAX_TIME=30
DB_MIN_TIME=5

# JWTToken
TOKEN_NAME=

# SecretKey
SECRET_KEY_ACCESS_TOKEN=
SECRET_KEY_REFRESH_TOKEN=
SECRET_KEY_FP_TOKEN=
SECRET_KEY_CSRF=
SECRET_TEST_CLIENT=

# Timeout
CACHE_TIMEOUT=8
DB_TIMEOUT=20
DOWN_STREAM_TIMEOUT=30

CORS_ALLOW_METHODS=
CORS_ALLOW_HEADERS=
CORS_ALLOW_ORIGINS=
CORS_ALLOW_CREDENTIALS=
CORS_EXPOSE_HEADERS=
```

## 📁 Project Structure

```
.
├── artifact
│   ├── api
│   └── coverage
├── build
├── client
├── db
│   ├── migrations
│   ├── postgresql
│   └── redis
├── docs
│   └── v1
├── internal
│   ├── configs
│   │   ├── bootstrap
│   │   ├── cache
│   │   ├── casbin
│   │   ├── fiber
│   │   │   └── resource
│   │   │       └── ssl
│   │   ├── hash
│   │   ├── logger
│   │   ├── orm
│   │   ├── pub-sub
│   │   ├── security
│   │   ├── sql
│   │   ├── timeout
│   │   └── token
│   ├── delivery
│   │   └── http
│   │       ├── middleware
│   │       └── route
│   ├── entity
│   ├── errors
│   ├── gateway
│   ├── model
│   │   ├── mapper
│   │   ├── request
│   │   └── response
│   ├── repository
│   ├── usecase
│   └── validation
├── openapi
│   ├── parameters
│   │   ├── path
│   │   └── query
│   ├── requests
│   │   ├── json
│   │   └── jsonapi
│   ├── resources
│   ├── responses
│   │   ├── json
│   │   └── jsonapi
│   ├── schemas
│   └── security
├── pkg
│   └── helper
│       ├── path
│       └── reflect
├── resource
│   ├── logs
│   ├── model
│   └── ssl
├── script
└── tests
    └── end2end
```

## 🔄 API Docs
https://go-fiber-auth-api-production.up.railway.app/api/v1/docs

## 🔒 Security Features

- JWT-based authentication
- Password hashing using bcrypt
- Rate limiting
- CORS protection
- XSS protection
- CSRF protection
- Request validation
- Session management with Redis

## 🚦 Development Commands

```bash
# Run development server
make run

# Run production server
make build


# Run tests
make tests

# Run linter
make lint

# Generate swagger docs
npm run build # with redocly

# Database commands
make db_up      # Run migrations up
make db_down    # Run migrations down
make db_version # Run migrations version to show version migration
```

## 🌐 Deployment

This service is configured for deployment on Railway. GitHub Actions workflow is included for CI/CD:

- Automated testing
- Database migrations
- Docker image building
- Railway deployment

## 📝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the APACHE License - see the [LICENSE](LICENSE.md) file for details.

## 🤝 Support

For support, email your-email@domain.com or open an issue in this repository.
