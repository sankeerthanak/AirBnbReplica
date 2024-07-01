package typesModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PropertyStore interface {
	CreateProperty(Property) error
	GetAllProperties() []primitive.M
	GetPropertiesByUserId(string) []primitive.M
	DeleteProperty(string) error
	// GetProperty(string) Property
	UpdateProperty(string, Property) error
	UploadToS3(string, string) error
}

type Property struct {
	PropertyId  primitive.ObjectID `json:"propertyId" bson:"propertyId"`
	Username    string             `json:"userName" bson:"userName"`
	UserId      string             `json:"userId" bson:"userId"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	StreetAddr  string             `json:"streetAddr" bson:"streetAddr"`
	City        string             `json:"city" bson:"city"`
	Country     string             `json:"country" bson:"country"`
	ZipCode     string             `json:"zipCode" bson:"zipCode"`
	Bedrooms    string             `json:"bedRooms" bson:"bedRooms"`
	Bathrooms   string             `json:"bathRooms" bson:"bathRooms"`
	Accomodates string             `json:"accomodates" bson:"accomodates"`
	Currency    string             `json:"currency" bson:"currency"`
	Price       string             `json:"price" bson:"price"`
	MinStay     string             `json:"minStay" bson:"minStay"`
	MaxStay     string             `json:"maxStay" bson:"maxStay"`
	// StartDate    string             `json:"startDate" bson:"startDate"`
	// EndDate      string             `json:"endDate" bson:"endDate"`
	PropertyType PropertyType `json:"propertyType" bson:"propertyType"`
	Amenities    Amenities    `json:"amenities" bson:"amenities"`
	Spaces       Spaces       `json:"spaces" bson:"spaces"`
	Image        string       `json:"image" bson:"image"`
}

type PropertyType struct {
	PrivateBed bool `json:"privateBed" bson:"privateBed"`
	Whole      bool `json:"whole" bson:"whole"`
	Shared     bool `json:"shared" bson:"shared"`
}

type Amenities struct {
	Ac     bool `json:"ac" bson:"ac"`
	Heater bool `json:"heater" bson:"heater"`
	TV     bool `json:"tv" bson:"tv"`
	Wifi   bool `json:"wifi" bson:"wifi"`
}

type Spaces struct {
	Kitchen bool `json:"kitchen" bson:"kitchen"`
	Closets bool `json:"closets" bson:"closets"`
	Parking bool `json:"parking" bson:"parking"`
	Gym     bool `json:"gym" bson:"gym"`
	Pool    bool `json:"pool" bson:"pool"`
}
