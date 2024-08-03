package auth

import (
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
	requestId, _ := c.Get("RequestId")
	if err := utils.ParseJSON(c, &payload); err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    err.Error(),
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		c.Header("Retry-After", "30")
		errors := err.(validator.ValidationErrors)
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    utils.FormatValidationError(errors),
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	u, err := h.store.GetUserById(payload.ID)
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusNotFound, gin.H{
			"message":    "User not found",
			"request_id": requestId,
			"code":       http.StatusNotFound,
		})
		return
	}

	if !middleware.ComparePasswords(u.Password, []byte(payload.Password)) {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusUnauthorized, gin.H{
			"message":    "Invalid credentials",
			"request_id": requestId,
			"code":       http.StatusUnauthorized,
		})
		return
	}

	token, err := middleware.GenerateJWT(u.User_id, u.UserType)
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    "Could not generate token",
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, gin.H{"token": token})
}

func (h *Handler) handleRegister(c *gin.Context) {
	requestId, _ := c.Get("RequestId")
	var payload models.RegisterUserPayload
	if err := utils.ParseJSON(c, &payload); err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    err.Error(),
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		c.Header("Retry-After", "30")
		errors := err.(validator.ValidationErrors)
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    utils.FormatValidationError(errors),
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	err := middleware.RoleGuard(payload.UserType)
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusForbidden, gin.H{
			"message":    "Invalid user type",
			"request_id": requestId,
			"code":       http.StatusForbidden,
		})
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err == nil && user.Email != "" {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusConflict, gin.H{
			"message":    "User already exists",
			"request_id": requestId,
			"code":       http.StatusConflict,
		})
		return
	}

	hashedPassword, err := middleware.HashPassword(payload.Password)
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    "Could not hash password",
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})
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
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    "Could not create user",
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})
		return
	}

	utils.WriteJSON(c, http.StatusCreated, gin.H{
		"user_id": newUser.User_id,
	})
}
