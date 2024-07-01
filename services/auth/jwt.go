package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sankeerthanak/airbnbreplica/config"
	typesModel "github.com/sankeerthanak/airbnbreplica/types"
	"github.com/sankeerthanak/airbnbreplica/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userId string) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":       userId,
		"expirationAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateBearer(secret []byte, userId primitive.ObjectID) (string, error) {
	token, err := CreateJWT(secret, userId.String())
	if err != nil {
		return "", err
	}
	return "Bearer " + token, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store typesModel.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)
		token, err := validateToken(tokenString)
		if err != nil {
			log.Println(err.Error())
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["userId"].(string)

		// u, err := store.GetUserById(userId)
		// if err != nil {
		// 	log.Printf("failed to get user by id: %v", err)
		// 	permissionDenied(w)
		// 	return
		// }

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userId)
		r = r.WithContext(ctx)

		validate := store.ValidateSession(userId, tokenString)
		if !validate {
			log.Printf("failed session %v", err)
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}
	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	t = strings.TrimPrefix(t, "Bearer ")
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) string {
	userId := ctx.Value(UserKey).(string)

	return userId
}
