package user

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	typesModel "github.com/sankeerthanak/airbnbreplica/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CollectionName = "UsersTable"

type Store struct {
	database  *mongo.Database
	rDatabase *redis.Client
}

func NewStore(db *mongo.Database, rDatabase *redis.Client) *Store {
	return &Store{database: db, rDatabase: rDatabase}
}

func (s *Store) GetUserByEmail(email string) (*typesModel.User, error) {

	filter := bson.M{"email": email}
	//fmt.Println(filter)
	collection := s.database.Collection(CollectionName)

	res := collection.FindOne(context.Background(), filter)
	//fmt.Println(collection.Name())
	user := new(typesModel.User)
	fmt.Println(user.Email)
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserById(id string) (*typesModel.User, error) {
	userId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": userId}
	collection := s.database.Collection(CollectionName)

	res := collection.FindOne(context.Background(), filter)
	user := new(typesModel.User)
	fmt.Println(user.Email)
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) CreateUser(user typesModel.User) error {
	collection := s.database.Collection(CollectionName)
	//fmt.Print(collection)
	_, err := collection.InsertOne(context.Background(), user)

	return err
}

func (s *Store) InsertJwt(token string, userId string) error {
	err := s.rDatabase.Set(context.Background(), userId, token, 0).Err()
	return err
}

func (s *Store) ValidateSession(userId string, token string) bool {
	//token = strings.TrimPrefix(token, "Bearer ")
	//redisToken, err := s.rDatabase.Get(context.Background(), userId).Result()

	// if err != nil || token != redisToken {
	// 	return false
	// }
	return true
}

func (s *Store) ValidateRole(users typesModel.User, role string) bool {
	for _, value := range users.Role {
		if role == value {
			return true
		}
	}
	return false
}
