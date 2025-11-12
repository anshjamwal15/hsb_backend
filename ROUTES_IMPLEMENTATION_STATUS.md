# API Routes Implementation Status

## ✅ Fully Implemented Routes

### Authentication (Public Routes)
- `POST /user/register` - Register a new user
- `POST /user/login` - User login
- `POST /user/forgot-password` - Request password reset
- `POST /user/verify-otp` - Verify OTP
- `POST /user/reset-password` - Reset password with OTP

### User Profile (Protected Routes)
- `POST /user/change-password` - Change password for authenticated user
- `GET /user/profile` - Get user profile
- `PUT /user/profile` - Update user profile

### Doctors (Protected Routes)
- `GET /api/doctors` - Get list of doctors with pagination and search
- `GET /api/doctors/:doctorId` - Get doctor details by ID

### Sessions & Bookings (Protected Routes)
- `GET /api/time-slots` - Get available time slots for a doctor
- `POST /api/bookings` - Create a booking with a doctor
- `POST /api/bookings/verify` - Verify booking payment
- `GET /api/bookings/my-with-doctors` - Get user's bookings with doctor details
- `GET /api/sessions/active` - Get active sessions

### Journals (Protected Routes)
- `GET /api/journals` - Get journals with pagination, search, and filters
- `POST /api/journals` - Create a new journal entry
- `GET /api/journals/:journalId` - Get journal by ID
- `PUT /api/journals/:journalId` - Update journal entry
- `DELETE /api/journals/:journalId` - Delete journal entry
- `GET /api/journals/user/:userId` - Get journals by user ID

## 🚧 Routes Marked for Implementation (TODO)

### Health Data
- `GET /api/health-data` - Get dashboard health data

### Sessions & Bookings
- `POST /api/agora/session-token` - Generate Agora session token for video/audio calls

### Clinics
- `GET /api/clinics` - Get list of clinics
- `POST /api/clinic-bookings` - Book clinic appointment
- `GET /api/clinic-bookings/my-bookings` - Get user's clinic bookings
- `POST /api/clinic-bookings/verify-payment` - Verify clinic payment

### Diagnostics
- `GET /api/public/diagnostics` - Get list of diagnostic tests
- `GET /api/diagnosticsUsers` - Get list of diagnostic labs
- `POST /api/diagnostics-bookings` - Book diagnostic test
- `POST /api/diagnostics-bookings/verify-payment` - Verify diagnostics payment

### Period Tracker
- `GET /api/period-cycle` - Get period cycle data
- `POST /api/period-cycle` - Add period cycle entry
- `DELETE /api/period-cycle` - Reset period tracker

### Pregnancy Tracker
- `GET /api/pregnancy-tracker` - Get pregnancy tracking data
- `POST /api/pregnancy-tracker` - Add pregnancy tracking entry

### Sexual Wellness (FSFI)
- `GET /api/fsfi/test` - Get FSFI test questions
- `POST /api/fsfi/submit` - Submit FSFI test
- `GET /api/fsfi/my-results` - Get FSFI test results

### Mental Health Tests
- `GET /api/tests` - Get list of mental health tests
- `GET /api/tests/:testName` - Get specific mental health test
- `POST /api/test-results` - Submit test results
- `GET /api/test-results` - Get test results history

### PCOS Assessment
- `GET /api/pcos-assessment/questions` - Get PCOS assessment questions
- `POST /api/pcos-assessment/submit` - Submit PCOS assessment
- `GET /api/pcos-assessment/history` - Get PCOS assessment history
- `GET /api/pcos-assessment` - Get latest PCOS assessment

### Symptoms Tracking
- `POST /api/symptoms-tracking` - Submit symptoms tracking
- `GET /api/symptoms-tracking` - Get symptoms tracking history

### Weight & Metabolic Wellness
- `GET /api/weight-metabolic-wellness` - Get weight and metabolic data
- `POST /api/weight-metabolic-wellness` - Add weight/metabolic entry

### Groups & Community
- `GET /api/groups` - Get all groups
- `GET /api/groups/joined` - Get joined groups
- `POST /api/groups/:groupId/join` - Join a group
- `POST /api/groups/:groupId/leave` - Leave a group
- `GET /api/groups/:groupId/posts` - Get group posts
- `POST /api/groups/:groupId/posts` - Create group post
- `GET /api/groups/posts` - Get all community posts
- `GET /api/groups/comments/my` - Get user's comments

### Blogs
- `GET /api/blogs` - Get blog articles
- `GET /api/blogs/:blogId` - Get blog by ID

### Chat
- `POST /api/chat/send` - Send chat message
- `GET /api/chat/history/:userId` - Get chat history
- `GET /api/chat/rooms` - Get chat rooms

### Media Gallery
- `GET /api/media-gallery` - Get media gallery
- `POST /api/media-gallery` - Upload media

### Surveys
- `GET /api/surveys` - Get available surveys
- `GET /api/surveys/:surveyId/questions` - Get survey questions

## 📝 Implementation Notes

### Current Architecture
- **Repositories**: User, Doctor, Booking, and Journal repositories are implemented
- **Services**: Auth, User, Doctor, Booking, and Journal services are implemented
- **Handlers**: Auth, User, Doctor, Booking, TimeSlot, and Journal handlers are implemented
- **Middleware**: JWT authentication middleware is implemented
- **Payment**: Razorpay payment client is integrated

### To Implement Next
1. Create handlers, services, and repositories for the remaining features
2. Implement Agora integration for video/audio calls
3. Add clinic and diagnostics booking functionality
4. Implement health tracking features (period, pregnancy, FSFI, mental health, PCOS, symptoms, weight)
5. Add community features (groups, posts, comments)
6. Implement blog and chat functionality
7. Add media gallery and survey features

### Swagger Documentation
The API documentation is available via Swagger UI:
- **Swagger UI**: `http://localhost:8080/swagger-ui` (or `/swagger`)
- **Swagger YAML**: `http://localhost:8080/swagger.yaml`
- **Health Check**: `http://localhost:8080/health`

### Testing
- All implemented routes can be tested using the provided endpoints
- Use Swagger UI for interactive API testing and documentation
- Ensure proper JWT token is included in the Authorization header for protected routes
- Use the format: `Authorization: Bearer <token>`
- You can add the JWT token in Swagger UI by clicking the "Authorize" button
