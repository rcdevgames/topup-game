# Waw Store - Topup Game Online Backend API

Backend API for Waw Store topup game online platform built with Go, Gin, GORM, PostgreSQL, and Redis.

## 🚀 Features

- **Authentication & Authorization**: JWT-based authentication with role-based access control
- **Database**: PostgreSQL with GORM ORM
- **Caching**: Redis for performance optimization (optional)
- **File Upload**: AWS S3 integration with local storage fallback
- **Payment Integration**: Ready for multiple payment gateways (Midtrans, Xendit, etc.)
- **WhatsApp Integration**: Notification system
- **Dynamic Forms**: Configurable game account validation forms
- **Voucher System**: Flexible discount voucher management
- **Analytics**: Dashboard with sales and transaction analytics
- **Documentation**: Swagger/OpenAPI documentation
- **Testing**: Comprehensive test coverage
- **Docker**: Containerized deployment

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Redis (optional)
- Docker & Docker Compose (for containerized deployment)

## 🛠️ Installation

### 1. Clone the repository

```bash
git clone <repository-url>
cd topup-game
```

### 2. Copy environment file

```bash
cp .env.example .env
```

### 3. Configure environment variables

Edit `.env` file with your database credentials and other configurations:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=topup_user
DB_PASSWORD=topup_password
DB_NAME=topup_db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Redis Configuration (Optional)
REDIS_URL=redis://localhost:6379

# AWS S3 Configuration (Optional)
AWS_S3_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Run the application

```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

## 🐳 Docker Deployment

### 1. Using Docker Compose (Recommended)

```bash
# Start all services (backend, postgres, redis)
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

### 2. Build Docker image only

```bash
# Build image
docker build -t topup-backend .

# Run container
docker run -p 8080:8080 --env-file .env topup-backend
```

## 📚 API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## 🗄️ Database

### Automatic Migration

The application automatically runs database migrations on startup. The following tables will be created:

- `users` - Customer users
- `admin_users` - Admin users
- `categories` - Product categories
- `products` - Game topup products
- `game_accounts` - User's game accounts
- `vouchers` - Discount vouchers
- `voucher_applications` - Voucher category/product associations
- `voucher_usage` - Voucher usage history
- `transactions` - Topup transactions
- `transaction_logs` - Transaction status history

### Sample Data

The application includes a seeder that creates sample data for development:

- Categories: Mobile Legends, Free Fire, Genshin Impact, PUBG Mobile
- Products: Various diamond/currency packages
- Default admin user: `username: admin, password: password123`

## 🔧 Key Endpoints

### Authentication

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh-token` - Refresh access token
- `GET /api/auth/profile` - Get user profile

### Products & Categories

- `GET /api/categories` - Get all categories
- `GET /api/products` - Get all products
- `GET /api/products/:id` - Get product details
- `GET /api/products/search` - Search products

### Transactions

- `POST /api/transactions` - Create transaction
- `GET /api/transactions` - Get user transactions
- `GET /api/transactions/:id` - Get transaction details

### Admin Panel

- `POST /api/admin/auth/login` - Admin login
- `GET /api/admin/analytics/dashboard` - Dashboard analytics
- `GET /api/admin/categories` - Manage categories
- `GET /api/admin/products` - Manage products
- `GET /api/admin/transactions` - Manage transactions

## 🏗️ Architecture

```
topup-backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection & migrations
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Database models
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic
│   └── utils/           # Utility functions
├── docs/                # Swagger documentation
├── cdn/                 # Local file storage
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Docker image definition
└── README.md
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./internal/services
```

## 🔐 Security Features

- JWT authentication with access & refresh tokens
- Password hashing with bcrypt
- Input validation and sanitization
- Rate limiting middleware
- CORS configuration
- Security headers middleware
- SQL injection prevention (GORM)

## 📱 Mobile API Support

The API is designed to be mobile-friendly with:

- Consistent response format
- Proper HTTP status codes
- Pagination support
- Error handling with meaningful messages
- Optimized queries for performance

## 🚀 Production Deployment

### Environment Variables

Make sure to set the following environment variables in production:

```env
NODE_ENV=production
JWT_SECRET=your-production-jwt-secret
DB_SSL_MODE=require
REDIS_URL=your-redis-url
AWS_S3_BUCKET=your-production-bucket
```

### Security Checklist

- [ ] Change default JWT secret
- [ ] Enable SSL/TLS for database connections
- [ ] Configure CORS for production domains
- [ ] Set up proper logging and monitoring
- [ ] Enable rate limiting
- [ ] Configure firewall rules
- [ ] Set up backup strategy

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For support and questions:

- Email: support@wawstore.com
- Documentation: [API Docs](http://localhost:8080/swagger/index.html)
- Issues: [GitHub Issues](https://github.com/your-repo/issues)

---

**Built with ❤️ for Waw Store**