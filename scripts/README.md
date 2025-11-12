# Test Data Script

This directory contains scripts for populating the database with test data.

## Prerequisites

- Go 1.16 or higher
- MongoDB running locally on default port 27017
- Go modules initialized (run `go mod init` if not already done)

## Setup

1. Install required Go packages:
   ```bash
   go get go.mongodb.org/mongo-driver/mongo
   go get golang.org/x/crypto/bcrypt
   ```

## Running the Script

1. Make sure your MongoDB server is running locally
2. Navigate to the project root directory
3. Run the script:
   ```bash
   go run scripts/insert_test_data.go
   ```

## Test Data Details

The script will insert the following test data:

### Users
- **User 1** (Doctor)
  - Email: user1@example.com
  - Password: password123
  - Role: Doctor (Gynecologist)

- **User 2** (Patient)
  - Email: user2@example.com
  - Password: password123
  - Role: Regular User

### Doctors
- 1 doctor profile linked to User 1
- Specialties: Gynecology, Obstetrics
- Consultation Fee: ₹1000

### Bookings
- 1 confirmed booking for User 2 with the doctor
- Session Type: Video Call
- Status: confirmed
- Payment Status: paid

### Journals
- 2 sample journal entries (1 public, 1 private)
- Different moods and content for testing

## Notes
- All passwords are hashed using bcrypt
- The script will drop existing collections before inserting test data
- Timestamps are set to the current time when the script is run
