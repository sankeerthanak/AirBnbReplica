package booking

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	typesModel "github.com/sankeerthanak/airbnbreplica/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "BookingsTable"

type Store struct {
	database *mongo.Database
}

func NewStore(database *mongo.Database) *Store {
	return &Store{database: database}
}

func (s *Store) InsertBooking(booking typesModel.Booking) error {

	collection := s.database.Collection(CollectionName)
	inserted, err := collection.InsertOne(context.Background(), booking)

	fmt.Println("Reservation successful for user", inserted.InsertedID)
	return err
}

func (s *Store) GetAllBookings() []primitive.M {

	collection := s.database.Collection(CollectionName)

	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var bookings []primitive.M

	for cur.Next(context.Background()) {
		var booking bson.M
		err := cur.Decode(&booking)
		if err != nil {
			log.Fatal(err)
		}
		bookings = append(bookings, booking)
	}

	defer cur.Close(context.Background())
	return bookings
}

func (s *Store) DeleteBookingById(userId string, bookingId string) (int64, error) {

	collection := s.database.Collection(CollectionName)

	bid, _ := primitive.ObjectIDFromHex(bookingId)
	uid, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": bid, "userId": uid}

	result, err := collection.DeleteOne(context.Background(), filter)

	fmt.Printf("Successfully deleted booking with userId %v", userId)

	res := result.DeletedCount

	return res, err
}

func (s *Store) GetUserBookingsbyId(userId string) []primitive.M {

	collection := s.database.Collection(CollectionName)

	id, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"userId": id}
	cur, err := collection.Find(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	var bookings []primitive.M

	for cur.Next(context.Background()) {
		var booking bson.M
		err := cur.Decode(&booking)
		if err != nil {
			log.Fatal(err)
		}
		bookings = append(bookings, booking)
	}

	defer cur.Close(context.Background())
	return bookings

}

func (s *Store) UpdateBookingById(booking typesModel.Booking) error {

	collection := s.database.Collection(CollectionName)

	filter := bson.M{"_id": booking.BookingId, "userId": booking.UserId}
	update := bson.M{"$set": bson.M{
		"NoofGuests":   booking.NoOfGuests,
		"CheckInDate":  booking.CheckInDate,
		"CheckOutDate": booking.CheckOutDate,
		"Amount":       booking.Amount,
		"Message":      booking.Message,
		"Email":        booking.Email,
	}}

	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err

}

func (s *Store) SendEmail(booking typesModel.Booking) error {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create an SES client
	svc := ses.New(sess)

	// Specify the email parameters
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(booking.Email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String("booking is successful"),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Regarding your Airbnb booking"),
			},
		},
		Source: aws.String("sankeerthana1234@gmail.com"), // Replace with your sender email address
	}

	// Send the email
	_, err = svc.SendEmail(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
