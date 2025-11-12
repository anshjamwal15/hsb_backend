package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClinicRepositoryImpl struct {
	clinicsCollection  *mongo.Collection
	bookingsCollection *mongo.Collection
}

func NewClinicRepository(db *mongo.Database) *ClinicRepositoryImpl {
	return &ClinicRepositoryImpl{
		clinicsCollection:  db.Collection("clinics"),
		bookingsCollection: db.Collection("clinic_bookings"),
	}
}

func (r *ClinicRepositoryImpl) FindAll(ctx context.Context) ([]*entities.Clinic, error) {
	cursor, err := r.clinicsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var clinics []*entities.Clinic
	if err = cursor.All(ctx, &clinics); err != nil {
		return nil, err
	}

	return clinics, nil
}

func (r *ClinicRepositoryImpl) CreateBooking(ctx context.Context, booking *entities.ClinicBooking) error {
	now := time.Now().Unix()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	result, err := r.bookingsCollection.InsertOne(ctx, booking)
	if err != nil {
		return err
	}

	booking.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *ClinicRepositoryImpl) FindBookingsByUserID(ctx context.Context, userID string) ([]*entities.ClinicBooking, error) {
	filter := bson.M{"userId": userID}

	cursor, err := r.bookingsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*entities.ClinicBooking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *ClinicRepositoryImpl) UpdateBookingPayment(ctx context.Context, bookingID, paymentID, status string) error {
	objectID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"paymentId":     paymentID,
			"paymentStatus": status,
			"status":        "confirmed",
			"updatedAt":     time.Now().Unix(),
		},
	}

	_, err = r.bookingsCollection.UpdateOne(ctx, filter, update)
	return err
}
