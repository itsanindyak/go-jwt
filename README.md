# Go JWT Authentication System

A secure and scalable JWT-based authentication system built with Go, featuring user management, role-based access control, and comprehensive middleware support.

## 🚀 Features

- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **User Management**: Complete CRUD operations for user accounts
- **Role-Based Access Control (RBAC)**: Fine-grained permissions system
- **Middleware Support**: Authentication and authorization middleware
- **OTP Verification**: One-time password support for enhanced security
- **MongoDB Integration**: Robust database operations with MongoDB
- **Docker Support**: Containerized deployment with Docker Compose
- **Environment Configuration**: Flexible environment-based configuration
- **OAuth2 Integration**: Google OAuth2 support with OIDC
- **Logging**: Structured logging throughout the application

## 📁 Project Structure

```
go-jwt/
├── .github/workflows/     # GitHub Actions CI/CD workflows
├── cmd/                   # Application entry points
├── config/                # Configuration management
├── controllers/           # HTTP request handlers
├── database/              # Database connection and operations
├── middleware/            # Authentication and authorization middleware
├── models/                # Data models and structures
├── pkg/                   # Shared packages and utilities
├── routes/                # API route definitions
├── .dockerignore          # Docker ignore file
├── .env.sample           # Environment variables template
├── .gitignore            # Git ignore file
├── docker-compose.yaml   # Docker Compose configuration
├── dockerfile            # Docker build instructions
├── go.mod                # Go module dependencies
└── go.sum                # Go module checksums
```

## 🛠️ Technology Stack

- **Backend**: Go (Golang)
- **Web Framework**: Gin Web Framework
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **OAuth2**: Google OIDC
- **Containerization**: Docker & Docker Compose
- **CI/CD**: GitHub Actions

## ⚙️ Setup Instructions

### Prerequisites

- Go 1.19 or higher
- MongoDB
- Docker (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/itsanindyak/go-jwt.git
   cd go-jwt
   ```

2. **Set up environment variables**
   ```bash
   cp .env.sample .env
   # Edit .env file with your configuration
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the application**
   ```bash
   go run cmd/main.go
   ```

### Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

2. **Run in detached mode**
   ```bash
   docker-compose up -d
   ```

## 🔧 Configuration

Key environment variables (see `.env.sample` for complete list):

- `PORT`: Server port (default: 8080)
- `MONGODB_URI`: MongoDB connection string
- `JWT_SECRET`: Secret key for JWT tokens
- `JWT_REFRESH_SECRET`: Secret key for refresh tokens
- `GOOGLE_OIDC_ISSUER`: Google OIDC issuer URL

## 📖 API Usage

### Authentication Endpoints

```bash
# Register a new user
POST /auth/register

# User login
POST /auth/login

# Refresh token
POST /auth/refresh

# Verify OTP
POST /auth/verify-otp
```

### User Management Endpoints

```bash
# Get user by ID (requires authentication)
GET /users/:user_id

# Get all users (admin only)
GET /users/

# Update user
PATCH /users/:user_id

# Delete user
DELETE /users/:user_id
```

## 🔐 Permissions System

The application implements a comprehensive permissions system:

- `PermissionReadSelf`: Read own user data
- `PermissionReadUser`: Read any user data
- `PermissionUpdateUser`: Update user data
- `PermissionDeleteUser`: Delete user accounts

## 🧪 Testing

Run tests with:
```bash
go test ./...
```

## 🚀 Deployment

The project includes GitHub Actions workflow for automated testing and deployment. See `.github/workflows/go.yml` for CI/CD configuration.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Go-OIDC](https://github.com/coreos/go-oidc)

## 📧 Contact

For questions or support, please open an issue on GitHub.

---

**Made with ❤️ by [itsanindyak](https://github.com/itsanindyak)**
