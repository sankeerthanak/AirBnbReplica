package booking

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/sankeerthanak/airbnbreplica/services/auth"
	typesModel "github.com/sankeerthanak/airbnbreplica/types"
	"github.com/sankeerthanak/airbnbreplica/utils"
)

type Handler struct {
	store     typesModel.BookingStore
	userStore typesModel.UserStore
}

func NewHandler(store typesModel.BookingStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/Bookings", auth.WithJWTAuth(h.showAllBookings, h.userStore)).Methods("GET")
	router.HandleFunc("/Booking", auth.WithJWTAuth(h.createBooking, h.userStore)).Methods("POST")
	router.HandleFunc("/Booking/{userId}/{bookingId}", auth.WithJWTAuth(h.deleteBooking, h.userStore)).Methods("DELETE")
	router.HandleFunc("/Booking/{userId}", auth.WithJWTAuth(h.showUserBookings, h.userStore)).Methods("GET")
	router.HandleFunc("/Booking/{userId}/{bookingId}", auth.WithJWTAuth(h.updateUserBooking, h.userStore)).Methods("POST")
}

func (h *Handler) createBooking(w http.ResponseWriter, r *http.Request) {

	var booking typesModel.Booking

	if err := utils.ParseJson(r, &booking); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(booking); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.InsertBooking(booking)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	err = h.store.SendEmail(booking)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to send notification to user"))
	}

	utils.WriteJson(w, http.StatusOK, booking)
}

func (h *Handler) showAllBookings(w http.ResponseWriter, r *http.Request) {

	bookings := h.store.GetAllBookings()

	if bookings == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no booking found"))
	}

	utils.WriteJson(w, http.StatusOK, bookings)
}

func (h *Handler) deleteBooking(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	res, err := h.store.DeleteBookingById(params["bookingId"])

	if res == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no booking found to delete"))
	}

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	json.NewEncoder(w).Encode(params["bookingId"])
}

func (h *Handler) showUserBookings(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	bookings := h.store.GetUserBookingsbyId(params["userId"])

	if bookings == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no booking found"))
	}

	utils.WriteJson(w, http.StatusOK, bookings)
}

func (h *Handler) updateUserBooking(w http.ResponseWriter, r *http.Request) {

	var booking typesModel.Booking

	if err := utils.ParseJson(r, &booking); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(booking); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.store.UpdateBookingById(booking)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

}
