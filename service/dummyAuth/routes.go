package dummyauth

import (
	"net/http"

	"github.com/delapaska/avito-rent/middleware"
	"github.com/delapaska/avito-rent/models"
	"github.com/delapaska/avito-rent/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store models.DummyStore
}

func NewHandler(store models.DummyStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.Use(middleware.GenerateRequestId())
	router.GET("/dummyLogin", h.handleDummyLogin)
}

// handleDummyLogin обрабатывает запросы на логин
// @Summary Dummy login
// @Tags Authentication
// @Description Получение JWT токена для dummy пользователя
// @Accept json
// @Produce json
// @Param userType query string true "Type of the user" Enums(client, moderator) example(client)
// @Success 200 {object} utils.LoginResponse "Successful login"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dummyLogin [get]
func (h *Handler) handleDummyLogin(c *gin.Context) {
	userType := c.Query("userType")
	requestId, _ := c.Get("RequestId")

	if userType == "" {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    "userType is required",
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	validUserTypes := map[string]bool{"client": true, "moderator": true}
	if _, valid := validUserTypes[userType]; !valid {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    "Invalid userType",
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	userID := uuid.New()
	tokenString, err := middleware.GenerateJWT(userID, userType)
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    "Could not generate token",
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, gin.H{"token": tokenString})
}
