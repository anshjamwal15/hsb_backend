# 🚀 Quick Start Guide

Get your HSB Backend up and running in 3 simple steps!

## Step 1: Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your settings (optional for local development)
nano .env
```

## Step 2: Start MongoDB

Choose one option:

**Option A: Using Docker (Recommended)**
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

**Option B: Local MongoDB**
```bash
mongod
```

## Step 3: Run the Server

```bash
# Make run script executable (first time only)
chmod +x run.sh

# Start the server
./run.sh
```

Or directly with Go:
```bash
go run cmd/server/main.go
```

## 🎉 You're Done!

When the server starts, you'll see:

```
═══════════════════════════════════════════════════════════════
🚀 HSB Backend Server Starting...
═══════════════════════════════════════════════════════════════

📡 Server URL:          http://localhost:8080
📚 Swagger UI:          http://localhost:8080/swagger
📄 Swagger YAML:        http://localhost:8080/swagger.yaml
❤️  Health Check:        http://localhost:8080/health

═══════════════════════════════════════════════════════════════
📋 Available Endpoints:
═══════════════════════════════════════════════════════════════

🔐 Authentication:
   • POST   /user/register
   • POST   /user/login
   • POST   /user/forgot-password
   • POST   /user/verify-otp
   • POST   /user/reset-password

👤 User Profile:
   • GET    /api/user/profile
   • PUT    /api/user/profile
   • POST   /api/user/change-password

👨‍⚕️ Doctors:
   • GET    /api/doctors
   • GET    /api/doctors/:doctorId

📅 Bookings:
   • GET    /api/time-slots
   • POST   /api/bookings
   • POST   /api/bookings/verify
   • GET    /api/bookings/my-with-doctors
   • GET    /api/sessions/active

═══════════════════════════════════════════════════════════════
✨ Server is running on port 8080
═══════════════════════════════════════════════════════════════

💡 Quick Start:
   1. Open Swagger UI: http://localhost:8080/swagger
   2. Test APIs directly from the browser
   3. Check health: curl http://localhost:8080/health

🎯 Ready to accept requests!
```

## 🧪 Test Your Setup

### 1. Check Health
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

### 2. Open Swagger UI

Open your browser and go to:
```
http://localhost:8080/swagger
```

You'll see an interactive API documentation where you can test all endpoints!

### 3. Register a User

Using Swagger UI:
1. Find the `POST /user/register` endpoint
2. Click "Try it out"
3. Fill in the request body:
```json
{
  "name": "Test User",
  "email": "test@example.com",
  "phoneNumber": "+919876543210",
  "password": "Test@123"
}
```
4. Click "Execute"

Or using cURL:
```bash
curl -X POST http://localhost:8080/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "phoneNumber": "+919876543210",
    "password": "Test@123"
  }'
```

### 4. Login

```bash
curl -X POST http://localhost:8080/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@123"
  }'
```

You'll get a JWT token in the response. Use this token for authenticated requests!

## 🔐 Using Authentication

For protected endpoints, add the JWT token to your requests:

**In Swagger UI:**
1. Click the "Authorize" button at the top
2. Enter: `Bearer YOUR_JWT_TOKEN`
3. Click "Authorize"

**In cURL:**
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🎯 What's Next?

- ✅ Explore all APIs in Swagger UI
- ✅ Test doctor listing and booking flows
- ✅ Set up Razorpay for payment testing
- ✅ Read the full [README.md](README.md) for more details

## 🆘 Troubleshooting

**MongoDB connection error?**
- Make sure MongoDB is running
- Check MONGODB_URI in .env file

**Port already in use?**
- Change PORT in .env file
- Or stop the process using port 8080

**Dependencies error?**
```bash
go mod tidy
go mod download
```

## 📚 Resources

- **Swagger UI**: http://localhost:8080/swagger
- **API Spec**: http://localhost:8080/swagger.yaml
- **Health Check**: http://localhost:8080/health
- **Full Documentation**: [README.md](README.md)
- **Implementation Status**: [API_IMPLEMENTATION_STATUS.md](API_IMPLEMENTATION_STATUS.md)

---

**Happy Coding! 🎉**
