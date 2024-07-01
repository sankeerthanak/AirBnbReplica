package property

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/sankeerthanak/airbnbreplica/services/auth"
	typesModel "github.com/sankeerthanak/airbnbreplica/types"
	"github.com/sankeerthanak/airbnbreplica/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	store     typesModel.PropertyStore
	userStore typesModel.UserStore
}

func NewHandler(store typesModel.PropertyStore, userStore typesModel.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/Property", auth.WithJWTAuth(h.showAllProperties, h.userStore)).Methods("GET")
	router.HandleFunc("/Property", auth.WithJWTAuth(h.createProperty, h.userStore)).Methods("POST")
	router.HandleFunc("/Property/{userId}", auth.WithJWTAuth(h.getUserProperties, h.userStore)).Methods("GET")
	router.HandleFunc("/Property/{propertyId}", auth.WithJWTAuth(h.deleteProperty, h.userStore)).Methods("DELETE")
	router.HandleFunc("/Property/{propertyId}", auth.WithJWTAuth(h.updateProperty, h.userStore)).Methods("POST")
}

func (h *Handler) createProperty(w http.ResponseWriter, r *http.Request) {

	//userId := auth.GetUserIdFromContext(r.Context())

	var property typesModel.Property

	if err := utils.ParseJson(r, &property); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(property); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	propertyId := primitive.NewObjectID()
	err := h.store.UploadToS3(property.Image, propertyId.Hex())
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	property.PropertyId = propertyId
	err = h.store.CreateProperty(property)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	utils.WriteJson(w, http.StatusOK, property)
}

func (h *Handler) getUserProperties(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	properties := h.store.GetPropertiesByUserId(params["userId"])

	if properties == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no property found"))
	}

	utils.WriteJson(w, http.StatusOK, properties)
}

func (h *Handler) updateProperty(w http.ResponseWriter, r *http.Request) {

	var property typesModel.Property
	params := mux.Vars(r)

	if err := utils.ParseJson(r, &property); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(property); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.UpdateProperty(params["propertyId"], property)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	utils.WriteJson(w, http.StatusOK, property)
}

func (h *Handler) showAllProperties(w http.ResponseWriter, r *http.Request) {

	properties := h.store.GetAllProperties()

	if properties == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no property found"))
	}

	utils.WriteJson(w, http.StatusOK, properties)

}

func (h *Handler) deleteProperty(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	error := h.store.DeleteProperty(params["propertyId"])

	if error != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no property found"))
	}
	utils.WriteJson(w, http.StatusOK, "successfully deleted the property")

}
