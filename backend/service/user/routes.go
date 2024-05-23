package user

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/utsavll0/ecom/config"
	"github.com/utsavll0/ecom/service/auth"
	"github.com/utsavll0/ecom/types"
	"github.com/utsavll0/ecom/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		_ = utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		var e validator.ValidationErrors
		errors.As(err, &e)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", e))
		return
	}
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		_ = utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	if err != nil {
		log.Fatal(err)
	}

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get json payload
	// if user doesnt exists create new user
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		_ = utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		var e validator.ValidationErrors
		errors.As(err, &e)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", e))
		return
	}
	log.Println("Checking if email exists...")
	_, err := h.store.GetUserByEmail(payload.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	log.Println("User doesnt Exist")

	hasedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hasedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
