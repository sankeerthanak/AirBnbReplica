package property

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sankeerthanak/airbnbreplica/config"
	typesModel "github.com/sankeerthanak/airbnbreplica/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "PropertyTable"

type Store struct {
	database *mongo.Database
}

func NewStore(database *mongo.Database) *Store {
	return &Store{database: database}
}

func (s *Store) CreateProperty(property typesModel.Property) error {

	collection := s.database.Collection(CollectionName)
	_, err := collection.InsertOne(context.Background(), property)

	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (s *Store) GetAllProperties() []primitive.M {
	collection := s.database.Collection(CollectionName)

	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var properties []primitive.M

	for cur.Next(context.Background()) {
		var property bson.M
		err := cur.Decode(&property)
		if err != nil {
			log.Fatal(err)
		}
		properties = append(properties, property)
	}

	defer cur.Close(context.Background())
	return properties
}

func (s *Store) GetPropertiesByUserId(userId string) []primitive.M {
	collection := s.database.Collection(CollectionName)
	filter := bson.M{"userId": userId}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	var properties []primitive.M

	for cur.Next(context.Background()) {
		var property bson.M
		err := cur.Decode(&property)
		if err != nil {
			log.Fatal(err)
		}
		properties = append(properties, property)
	}

	defer cur.Close(context.Background())
	return properties
}

func (s *Store) DeleteProperty(propertyId string) error {

	collection := s.database.Collection(CollectionName)

	pid, _ := primitive.ObjectIDFromHex(propertyId)
	filter := bson.M{"_id": pid}

	_, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	return err
}

// func (s *Store) GetProperty(userId string, propertyId string) typesModel.Property {

// }

func (s *Store) UpdateProperty(propertyId string, listing typesModel.Property) error {
	collection := s.database.Collection(CollectionName)
	id, _ := primitive.ObjectIDFromHex(propertyId)
	filter := bson.M{"propertyId": id}
	update := bson.M{"$set": bson.M{"title": listing.Title, "description": listing.Description, "StreetAddr": listing.StreetAddr, "City": listing.City, "Country": listing.Country, "ZipCode": listing.ZipCode, "Bedrooms": listing.Bedrooms, "Bathrooms": listing.Bathrooms,
		"accomodates": listing.Accomodates, "currency": listing.Currency, "price": listing.Price, "minStay": listing.MinStay,
		"maxStay": listing.MaxStay, "propertyType.privateBed": listing.PropertyType.PrivateBed, "PropertyType.whole": listing.PropertyType.Whole,
		"PropertyType.shared": listing.PropertyType.Shared, "amenities.ac": listing.Amenities.Ac, "amenities.heater": listing.Amenities.Heater, "amenities.tv": listing.Amenities.TV, "amenities.wifi": listing.Amenities.Wifi, "spaces.kitchen": listing.Spaces.Kitchen, "spaces.closets": listing.Spaces.Closets, "spaces.parking": listing.Spaces.Parking, "spaces.gym": listing.Spaces.Gym, "spaces.pool": listing.Spaces.Pool}}

	_, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (s *Store) UploadToS3(image string, key string) error {

	dec, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %v", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Envs.Region), // Replace with your AWS region
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)

	_, err = svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(config.Envs.S3BucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %v", err)
	}

	// Upload input parameters
	uploader := s3manager.NewUploader(sess)
	input := &s3manager.UploadInput{
		Bucket: aws.String(config.Envs.S3BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(dec),
	}

	// Perform an upload.
	_, err = uploader.Upload(input)
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %v", err)
	}

	fmt.Printf("Successfully uploaded image to %s/%s\n", config.Envs.S3BucketName, key)
	return nil
}

func (s *Store) RetrieveFromS3(fileHeader *multipart.FileHeader) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Replace with your AWS region
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)

	rawObject, err := svc.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String("toto"),
			Key:    aws.String("toto.txt"),
		})

	if err != nil {
		return "", fmt.Errorf("failed to read data %v", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rawObject.Body)

	if err != nil {
		return "", fmt.Errorf("failed to read data %v", err)
	}

	// Encode the contents to base64 string
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Str, nil
}
