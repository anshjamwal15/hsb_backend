package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"github.com/anshjamwal15/hsb_backend/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bookingRepositoryImpl struct {
	collection *mongo.Collection
}

func NewBookingRepository(db *mongo.Database) repositories.BookingRepository {
	return &bookingRepositoryImpl{
		collection: db.Collection("bookings"),
	}
}

func (r *bookingRepositoryImpl) Create(ctx context.Context, booking *entities.Booking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, booking)
	if err != nil {
		return err
	}
	booking.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *bookingRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid booking ID format: %w", err)
	}

	var booking entities.Booking
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepositoryImpl) FindByUserID(ctx context.Context, userID string, page, limit int) ([]*entities.Booking, int64, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := bson.M{"userId": userObjID}
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var bookings []*entities.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (r *bookingRepositoryImpl) FindByDoctorAndTimeSlot(ctx context.Context, doctorID string, date time.Time, timeSlot string) ([]*entities.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID format: %w", err)
	}

	timeRange := strings.Split(timeSlot, "-")
	if len(timeRange) != 2 {
		return nil, fmt.Errorf("invalid time slot format")
	}

	filter := bson.M{
		"doctorId": objID,
		"date": bson.M{
			"$gte": time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, time.UTC),
		},
		"timeSlot": timeSlot,
		"status": bson.M{"$nin": []string{"cancelled", "rejected"}},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding bookings: %w", err)
	}
	defer cursor.Close(ctx)

	var bookings []*entities.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("error decoding bookings: %w", err)
	}

	return bookings, nil
}

func (r *bookingRepositoryImpl) FindByDoctorID(ctx context.Context, doctorID string, startDate, endDate time.Time) ([]*entities.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(doctorID)
	if err != nil {
		return nil, fmt.Errorf("invalid doctor ID format: %w", err)
	}

	filter := bson.M{
		"doctorId": objID,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"status": bson.M{"$nin": []string{"cancelled", "rejected"}},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding bookings: %w", err)
	}
	defer cursor.Close(ctx)

	var bookings []*entities.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("error decoding bookings: %w", err)
	}

	return bookings, nil
}

func (r *bookingRepositoryImpl) UpdatePaymentStatus(ctx context.Context, bookingID, status, razorpayOrderID, razorpayPaymentID string) error {
	objID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return fmt.Errorf("invalid booking ID format: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"paymentStatus":     status,
			"razorpayOrderId":  razorpayOrderID,
			"razorpayPaymentId": razorpayPaymentID,
			"updatedAt":         time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return fmt.Errorf("error updating payment status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

func (r *bookingRepositoryImpl) Update(ctx context.Context, booking *entities.Booking) error {
	booking.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": booking.ID},
		bson.M{"$set": booking},
	)
	return err
}

func (r *bookingRepositoryImpl) FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Booking, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := bson.M{
		"userId": userObjID,
		"status": bson.M{"$in": []string{"pending", "confirmed"}},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*entities.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
