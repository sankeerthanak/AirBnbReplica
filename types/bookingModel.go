package typesModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookingStore interface {
	InsertBooking(Booking) error
	GetAllBookings() []primitive.M
	DeleteBookingById(string) (int64, error)
	GetUserBookingsbyId(string) []primitive.M
	UpdateBookingById(Booking) error
	SendEmail(Booking) error
}

type Booking struct {
	BookingId    primitive.ObjectID `json:"bookingId" bson:"_id"`
	PropertyId   string             `json:"propertyId"`
	NoOfGuests   int                `json:"noofguests"`
	CheckInDate  string             `json:"checkInDate"`
	CheckOutDate string             `json:"checkOutDate"`
	UserId       string             `json:"userId"`
	UserName     string             `json:"userName"`
	Email        string             `json:"email"`
	Message      string             `json:"message"`
	Amount       string             `json:"amount"`
	Reservation  bool               `json:"reservation"`
}
