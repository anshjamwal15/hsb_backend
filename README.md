# HSB Backend - Healthcare Services Backend

A comprehensive healthcare backend system built with Go, featuring doctor management, booking system, and payment integration.

## 🚀 Quick Start

### Prerequisites

**Option A: Docker (Recommended)**
- Docker 20.10 or higher
- Docker Compose 2.0 or higher

**Option B: Local Development**
- Go 1.21 or higher
- MongoDB 4.4 or higher
- Razorpay account (for payment integration)

### Installation

#### 🐳 Using Docker (Easiest Way)

```bash
# 1. Copy environment file
cp .env.docker .env

# 2. Start everything with one command
make docker-run

# Or using docker-compose directly
docker-compose up -d
```

That's it! Your server is running at http://localhost:8080/swagger

See [DOCKER.md](DOCKER.md) for complete Docker documentation.

#### 💻 Local Development Setup

1. **Clone the repository**
```bash
cd /Users/ansh/volvrit/production/hsb_backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=hsb_backend
JWT_SECRET=your-secret-key
RAZORPAY_KEY_ID=your_razorpay_key
RAZORPAY_KEY_SECRET=your_razorpay_secret
```

4. **Run MongoDB**
```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or use your local MongoDB installation
mongod
```

5. **Run the server**
```bash
go run cmd/server/main.go
```

The server will start and display:
```
═══════════════════════════════════════════════════════════════
🚀 HSB Backend Server Starting...
═══════════════════════════════════════════════════════════════

📡 Server URL:          http://localhost:8080
📚 Swagger UI:          http://localhost:8080/swagger
📄 Swagger YAML:        http://localhost:8080/swagger.yaml
❤️  Health Check:        http://localhost:8080/health
```

## 📚 API Documentation

### Access Swagger UI

Open your browser and navigate to:
```
http://localhost:8080/swagger
```

This will display an interactive API documentation where you can:
- View all available endpoints
- Test APIs directly from the browser
- See request/response schemas
- Try out authentication flows

### API Endpoints

#### Authentication
- `POST /user/register` - Register new user
- `POST /user/login` - User login
- `POST /user/forgot-password` - Request password reset
- `POST /user/verify-otp` - Verify OTP
- `POST /user/reset-password` - Reset password
- `POST /api/user/change-password` - Change password (authenticated)

#### User Profile
- `GET /api/user/profile` - Get user profile
- `PUT /api/user/profile` - Update user profile

#### Doctors
- `GET /api/doctors` - List all doctors (with pagination & search)
- `GET /api/doctors/:doctorId` - Get doctor details

#### Bookings
- `GET /api/time-slots` - Get available time slots
- `POST /api/bookings` - Create new booking
- `POST /api/bookings/verify` - Verify payment
- `GET /api/bookings/my-with-doctors` - Get user's bookings
- `GET /api/sessions/active` - Get active sessions

## 🧪 Testing APIs

### Using Swagger UI (Recommended)

1. Open http://localhost:8080/swagger
2. Click on any endpoint
3. Click "Try it out"
4. Fill in the parameters
5. Click "Execute"

### Using cURL

**Register a new user:**
```bash
curl -X POST http://localhost:8080/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phoneNumber": "+919876543210",
    "password": "SecurePass123!"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Get doctors (authenticated):**
```bash
curl -X GET "http://localhost:8080/api/doctors?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Using Postman

1. Import the `swagger.yaml` file into Postman
2. Set up environment variables for the base URL and token
3. Test endpoints directly

## 🏗️ Project Structure

```
hsb_backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── application/
│   │   └── services/              # Business logic
│   ├── domain/
│   │   ├── entities/              # Domain models
│   │   └── repositories/          # Repository interfaces
│   ├── infrastructure/
│   │   ├── mongodb/               # MongoDB implementation
│   │   ├── payment/               # Payment integration
│   │   └── repositories/          # Repository implementations
│   └── presentation/
│       └── http/
│           ├── handlers/          # HTTP handlers
│           ├── middleware/        # Middleware
│           └── router/            # Route configuration
├── pkg/
│   └── auth/                      # Authentication utilities
├── swagger.yaml                   # API specification
├── .env.example                   # Environment variables template
└── README.md                      # This file
```

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication.

1. Register or login to get a JWT token
2. Include the token in the Authorization header:
   ```
   Authorization: Bearer YOUR_JWT_TOKEN
   ```
3. All `/api/*` endpoints require authentication

## 💳 Payment Integration

The system integrates with Razorpay for payment processing:

1. **Create Booking** - Generates a Razorpay order
2. **User Pays** - Frontend handles Razorpay checkout
3. **Verify Payment** - Backend verifies payment signature
4. **Update Status** - Booking status updated to confirmed

## 🗄️ Database

MongoDB collections:
- `users` - User accounts
- `doctors` - Doctor profiles
- `bookings` - Appointment bookings
- `otps` - OTP for password reset

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `MONGODB_URI` | MongoDB connection string | mongodb://localhost:27017 |
| `MONGODB_DATABASE` | Database name | hsb_backend |
| `JWT_SECRET` | Secret key for JWT | - |
| `RAZORPAY_KEY_ID` | Razorpay key ID | - |
| `RAZORPAY_KEY_SECRET` | Razorpay key secret | - |

## 🛠️ Development

### Run in development mode
```bash
go run cmd/server/main.go
```

### Build for production
```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

### Run tests
```bash
go test ./...
```

## 📊 Features

- ✅ User authentication & authorization
- ✅ Doctor management
- ✅ Booking system with time slot management
- ✅ Payment integration (Razorpay)
- ✅ JWT-based security
- ✅ MongoDB persistence
- ✅ RESTful API design
- ✅ Swagger documentation
- ✅ CORS support

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

---

**Built with ❤️ using Go and MongoDB**
