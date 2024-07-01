package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sankeerthanak/airbnbreplica/services/booking"
	"github.com/sankeerthanak/airbnbreplica/services/property"
	"github.com/sankeerthanak/airbnbreplica/services/user"
	"github.com/sankeerthanak/airbnbreplica/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type ApiServer struct {
	address   string
	dataBase  *mongo.Database
	rDatabase *redis.Client
}

func NewApiServer(address string, dataBase *mongo.Database, rDatabase *redis.Client) *ApiServer {
	return &ApiServer{
		address:   address,
		dataBase:  dataBase,
		rDatabase: rDatabase,
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()
	//subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.dataBase, s.rDatabase)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	bookingStore := booking.NewStore(s.dataBase)
	bookinHandler := booking.NewHandler(bookingStore, userStore)
	bookinHandler.RegisterRoutes(router)

	propertyStore := property.NewStore(s.dataBase)
	propertyHandler := property.NewHandler(propertyStore, userStore)
	propertyHandler.RegisterRoutes(router)

	log.Println("Listening on", s.address)

	return http.ListenAndServe(s.address, utils.EnableCORS(router))
	//return http.ListenAndServe(s.address, router)

}
