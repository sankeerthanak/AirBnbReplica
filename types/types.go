package typesModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStore interface {
	CreateUser(User) error
	GetUserByEmail(string) (*User, error)
	GetUserById(string) (*User, error)
	InsertJwt(string, string) error
	ValidateSession(string, string) error
}

type User struct {
	UserId    primitive.ObjectID `json:"userId" bson:"_id"`
	UserName  string             `json:"username" validate:"required"`
	FirstName string             `json:"firstname" validate:"required"`
	LastName  string             `json:"lastname" validate:"required"`
	Email     string             `json:"email" validate:"required,email"`
	Password  string             `json:"password" validate:"required"`
	Role      string             `json:"rolename" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
