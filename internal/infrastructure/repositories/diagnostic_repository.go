package repositories

import (
	"context"
	"time"

	"github.com/anshjamwal15/hsb_backend/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiagnosticRepositoryImpl struct {
	diagnosticsCollection *mongo.Collection
	bookingsCollection    *mongo.Collection
}

func NewDiagnosticRepository(db *mongo.Database) *DiagnosticRepositoryImpl {
	return &DiagnosticRepositoryImpl{
		diagnosticsCollection: db.Collection("diagnostics"),
		bookingsCollection:    db.Collection("diagnostic_bookings"),
	}
}

func (r *DiagnosticRepositoryImpl) FindAll(ctx context.Context) ([]*entities.Diagnostic, error) {
	cursor, err := r.diagnosticsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var diagnostics []*entities.Diagnostic
	if err = cursor.All(ctx, &diagnostics); err != nil {
		return nil, err
	}

	return diagnostics, nil
}

func (r *DiagnosticRepositoryImpl) CreateBooking(ctx context.Context, booking *entities.DiagnosticBooking) error {
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

func (r *DiagnosticRepositoryImpl) FindBookingsByUserID(ctx context.Context, userID string) ([]*entities.DiagnosticBooking, error) {
	filter := bson.M{"userId": userID}

	cursor, err := r.bookingsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*entities.DiagnosticBooking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *DiagnosticRepositoryImpl) UpdateBookingPayment(ctx context.Context, bookingID, paymentID, status string) error {
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
