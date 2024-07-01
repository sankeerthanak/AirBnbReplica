package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/sankeerthanak/airbnbreplica/config"
	"github.com/sankeerthanak/airbnbreplica/services/auth"
	typesModel "github.com/sankeerthanak/airbnbreplica/types"
	"github.com/sankeerthanak/airbnbreplica/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
)

type Handler struct {
	store typesModel.UserStore
}

func NewHandler(store typesModel.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payLoad typesModel.LoginUserPayload

	if err := utils.ParseJson(r, &payLoad); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(payLoad); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(payLoad.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user does not exist"))
		return
	}

	if !h.store.ValidateRole(*u, payLoad.Role) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user does not have access to given role"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(payLoad.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.UserId.Hex())

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error while creating token"))
		return
	}

	err = h.store.InsertJwt(token, u.UserId.Hex())

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error pushing token"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		//Secure:   true, // Set to true in production for HTTPS
		Expires: time.Now().Add(time.Hour),
	})

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token, "userId": u.UserId.Hex()})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payLoad typesModel.User
	//json.NewDecoder(r.Body).Decode(&payLoad)
	//fmt.Print("1")
	if err := utils.ParseJson(r, &payLoad); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	//fmt.Print("2")

	if err := utils.Validate.Struct(payLoad); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	_, err := h.store.GetUserByEmail(payLoad.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payLoad.Email))
		return
	}
	//fmt.Print(user)

	hashedPassword, err := auth.HashPassword(payLoad.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	Error := h.store.CreateUser(typesModel.User{
		UserId:    primitive.NewObjectID(),
		FirstName: payLoad.FirstName,
		LastName:  payLoad.LastName,
		Email:     payLoad.Email,
		Password:  hashedPassword,
		//UserName:  payLoad.UserName,
		Role: payLoad.Role,
	})
	//fmt.Print("4")

	if Error != nil {
		fmt.Println(("In eror block/n"))
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	//fmt.Print("5")

	utils.WriteJson(w, http.StatusOK, payLoad)
}
