package auth

import (
	"fmt"
	"net/http"

	"github.com/delapaska/avito-rent/middleware"
	"github.com/delapaska/avito-rent/models"
	"github.com/delapaska/avito-rent/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store models.UserStore
}

func NewHandler(store models.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	router.POST("/login", h.handleLogin)
	router.POST("/register", h.handleRegister)

}

func (h *Handler) handleLogin(c *gin.Context) {
	var payload models.LoginUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	u, err := h.store.GetUserById(payload.ID)
	if err != nil {
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	if !middleware.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	token, err := middleware.GenerateJWT(u.User_id, u.UserType)
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c, http.StatusOK, gin.H{"token": token})
}
func (h *Handler) handleRegister(c *gin.Context) {
	var payload models.RegisterUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := middleware.RoleGuard(payload.UserType)
	if err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}
	user, err := h.store.GetUserByEmail(payload.Email)
	if err == nil && user.Email != "" {
		utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := middleware.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	newUser := models.User{
		User_id:  uuid.New(),
		Email:    payload.Email,
		Password: hashedPassword,
		UserType: payload.UserType,
	}

	err = h.store.CreateUser(newUser)
	if err != nil {
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(c, http.StatusCreated, gin.H{
		"user_id": newUser.User_id,
	})
}
