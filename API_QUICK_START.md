# Women's Health API - Quick Start Guide

## 🚀 Starting the Server

```bash
# Make sure your .env file is configured with:
# - MONGO_URI
# - JWT_SECRET
# - RAZORPAY_KEY
# - RAZORPAY_SECRET

# Run the server
go run cmd/main.go
```

The server will start on `http://localhost:8080` (or the port specified in your .env)

## 📚 API Documentation

Once the server is running, access the interactive API documentation:

- **Swagger UI**: http://localhost:8080/swagger-ui
- **Swagger YAML**: http://localhost:8080/swagger.yaml
- **Health Check**: http://localhost:8080/health

## 🔐 Authentication Flow

### 1. Register a New User
```bash
POST http://localhost:8080/user/register
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "phoneNumber": "+919876543210",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "689b36243fcf0bb96f7d2abf",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "phoneNumber": "+919876543210",
    "profileImage": ""
  }
}
```

### 2. Login
```bash
POST http://localhost:8080/user/login
Content-Type: application/json

{
  "email": "jane@example.com",
  "password": "SecurePass123!"
}
```

### 3. Use the Token
For all protected routes, include the token in the Authorization header:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 📋 Available Endpoints

### Public Endpoints (No Auth Required)
- `POST /user/register` - Register new user
- `POST /user/login` - Login
- `POST /user/forgot-password` - Request password reset
- `POST /user/verify-otp` - Verify OTP
- `POST /user/reset-password` - Reset password
- `GET /health` - Health check
- `GET /swagger-ui` - API documentation
- `GET /swagger.yaml` - OpenAPI specification

### Protected Endpoints (Auth Required)

#### User Profile
- `GET /user/profile` - Get user profile
- `PUT /user/profile` - Update profile
- `POST /user/change-password` - Change password

#### Doctors
- `GET /api/doctors` - List doctors (with pagination & search)
- `GET /api/doctors/:doctorId` - Get doctor details

#### Bookings
- `GET /api/time-slots?doctorId=xxx&date=2025-11-15` - Get available time slots
- `POST /api/bookings` - Create booking
- `POST /api/bookings/verify` - Verify payment
- `GET /api/bookings/my-with-doctors` - Get my bookings
- `GET /api/sessions/active` - Get active sessions

#### Journals
- `GET /api/journals` - List journals (with filters)
- `POST /api/journals` - Create journal
- `GET /api/journals/:journalId` - Get journal
- `PUT /api/journals/:journalId` - Update journal
- `DELETE /api/journals/:journalId` - Delete journal
- `GET /api/journals/user/:userId` - Get user's journals

## 🧪 Testing with Swagger UI

1. Open http://localhost:8080/swagger-ui
2. Click the "Authorize" button at the top
3. Enter your JWT token in the format: `Bearer <your-token>`
4. Click "Authorize" and then "Close"
5. Now you can test any endpoint directly from the UI

## 🔧 Example API Calls

### Get Doctors List
```bash
curl -X GET "http://localhost:8080/api/doctors?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Create a Journal Entry
```bash
curl -X POST "http://localhost:8080/api/journals" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "689b36243fcf0bb96f7d2abf",
    "title": "My Daily Reflection",
    "content": "Today was a good day...",
    "category": "Mental Health"
  }'
```

### Book a Doctor Appointment
```bash
curl -X POST "http://localhost:8080/api/bookings" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "doctorId": "60d5ec49f1b2c72b8c8e4f1a",
    "sessionType": "Video Call",
    "date": "2025-11-15",
    "timeSlot": "10:00 AM",
    "notes": "First consultation"
  }'
```

## 📊 Response Format

All API responses follow this format:

**Success Response:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Optional success message"
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## 🔍 Query Parameters

### Pagination
Most list endpoints support pagination:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)

### Search & Filters
- `search` - Search term
- `category` - Filter by category
- `date` - Filter by date (YYYY-MM-DD format)

Example:
```
GET /api/journals?page=1&limit=20&search=health&category=Mental Health
```

## 🛠️ Environment Variables

Required environment variables in `.env`:

```env
# Server
PORT=8080
ENVIRONMENT=development

# Database
MONGO_URI=mongodb://localhost:27017
DB_NAME=hsb_backend

# Security
JWT_SECRET=your-secret-key-here

# Payment
RAZORPAY_KEY=your-razorpay-key
RAZORPAY_SECRET=your-razorpay-secret

# Agora (for video calls)
AGORA_APP_ID=your-agora-app-id
AGORA_APP_CERT=your-agora-certificate

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## 📝 Notes

- All dates should be in `YYYY-MM-DD` format
- All timestamps are in ISO 8601 format
- File uploads use `multipart/form-data`
- Payment amounts are in INR (Indian Rupees)
- JWT tokens expire after 24 hours (configurable)

## 🐛 Troubleshooting

### "Unauthorized" Error
- Make sure you're including the Authorization header
- Check that your token hasn't expired
- Verify the token format: `Bearer <token>`

### "Doctor not found" Error
- Ensure the doctor exists in the database
- Check that you're using the correct doctor ID format (MongoDB ObjectID)

### Payment Verification Failed
- Verify Razorpay credentials are correct in .env
- Check that the payment signature is valid
- Ensure the order ID matches the booking

## 📞 Support

For issues or questions, refer to:
- Swagger documentation at `/swagger-ui`
- Implementation status in `ROUTES_IMPLEMENTATION_STATUS.md`
- Source code in the respective handler files
