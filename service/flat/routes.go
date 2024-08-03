package flat

import (
	"fmt"
	"net/http"

	"github.com/delapaska/avito-rent/middleware"
	"github.com/delapaska/avito-rent/models"
	"github.com/delapaska/avito-rent/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	store models.FlatStore
}

func NewHandler(store models.FlatStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	allUsers := router.Group("/")
	allUsers.Use(middleware.AuthMiddleware("moderator", "client"))
	{
		allUsers.POST("/flat/create", h.handleCreateFlat)

	}
	moderationsOnly := router.Group("/")
	moderationsOnly.Use(middleware.AuthMiddleware("moderator"))
	{
		moderationsOnly.POST("/flat/update", h.handleUpdateFlatStatus)
	}

}

func (h *Handler) handleCreateFlat(c *gin.Context) {
	requestId, _ := c.Get("RequestId")
	var payload models.FlatPayload
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

	flat, err := h.store.CreateFlat(models.Flat{
		House_id: payload.House_id,
		Price:    payload.Price,
		Rooms:    payload.Rooms,
	})
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    err.Error(),
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, flat)
}

func (h *Handler) handleUpdateFlatStatus(c *gin.Context) {
	var payload models.UpdateStatusPayload
	requestId, _ := c.Get("RequestId")
	userID, exists := c.Get("userID")
	if !exists {
		utils.WriteJSON(c, http.StatusUnauthorized, gin.H{
			"message":    "userID not found in context",
			"request_id": requestId,
			"code":       http.StatusUnauthorized,
		})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.WriteJSON(c, http.StatusUnauthorized, gin.H{
			"message":    "userID is not of type uuid.UUID",
			"request_id": requestId,
			"code":       http.StatusUnauthorized,
		})
		return
	}

	if err := utils.ParseJSON(c, &payload); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if !models.ValidStatuses[payload.Status] {
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    fmt.Sprintf("invalid status %s", payload.Status),
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	flat, err := h.store.UpdateFlatStatus(userIDUUID, payload)
	if err != nil {
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    err.Error(),
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, gin.H{
		"flat":       flat,
		"request_id": requestId,
		"code":       http.StatusOK,
	})
}
