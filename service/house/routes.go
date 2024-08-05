package house

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/delapaska/avito-rent/middleware"
	"github.com/delapaska/avito-rent/models"
	"github.com/delapaska/avito-rent/sender"
	"github.com/delapaska/avito-rent/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store models.HouseStore
}

func NewHandler(store models.HouseStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	moderationsOnly := router.Group("/")
	moderationsOnly.Use(middleware.AuthMiddleware("moderator"))
	{
		moderationsOnly.POST("/house/create", h.handleCreateHouse)
	}
	allUsers := router.Group("/")
	allUsers.Use(middleware.AuthMiddleware("moderator", "client"))
	{
		allUsers.POST("/house/:id/subscribe", h.handleSubscribeHouse)
		allUsers.GET("/house/:id", h.handleGetHouseFlats)
	}
}

// handleCreateHouse creates a new house
// @Summary Create House
// @Tags House
// @Description Create a new house. Requires moderator access.
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.HousePayload true "House details"
// @Success 201 {object} models.House "House created"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /house/create [post]
func (h *Handler) handleCreateHouse(c *gin.Context) {
	var payload models.HousePayload
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
	if payload.Developer != "" && len(strings.Trim(payload.Developer, " ")) == 0 {
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    "address cannot be empty",
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}

	if payload.Address != "" && len(strings.Trim(payload.Address, " ")) == 0 {
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    "address cannot be empty",
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}
	if payload.Year < 0 {
		utils.WriteJSON(c, http.StatusBadRequest, gin.H{
			"message":    "year must be more or equal than 0",
			"request_id": requestId,
			"code":       http.StatusBadRequest,
		})
		return
	}
	house, err := h.store.CreateHouse(models.House{
		Address:   payload.Address,
		Year:      payload.Year,
		Developer: payload.Developer,
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

	utils.WriteJSON(c, http.StatusCreated, house)
}

// @Summary Get House Flats
// @Description Retrieve flats for a specific house. Requires authorization for both moderator and client.
// @Tags House
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "House ID"
// @Success 200 {object} utils.FlatsResponse "Flats retrieved"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /house/{id} [get]
func (h *Handler) handleGetHouseFlats(c *gin.Context) {
	requestId, _ := c.Get("RequestId")
	houseID := c.Param("id")

	userType := c.GetString("userType")

	flats, err := h.store.GetHouseFlats(houseID, userType)
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    err.Error(),
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})

		return
	}

	utils.WriteJSON(c, http.StatusOK,
		gin.H{"flats": flats})
}

// @Summary Subscribe to House
// @Description Subscribe to updates for a specific house. Requires authorization for both moderator and client.
// @Tags House
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "House ID"
// @Param request body models.SubscribePayload true "Subscription details"
// @Success 201 {object} utils.SubscriptionResponse "Subscription successful"
// @Failure 400 {object} utils.ErrorResponse "Bad request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /house/{id}/subscribe [post]
func (h *Handler) handleSubscribeHouse(c *gin.Context) {
	requestId, _ := c.Get("RequestId")
	houseID := c.Param("id")

	var payload models.SubscribePayload

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

	err := h.store.AddSubscription(houseID, payload.Email)
	if err != nil {
		c.Header("Retry-After", "30")
		utils.WriteJSON(c, http.StatusInternalServerError, gin.H{
			"message":    "Failed to save subscription",
			"request_id": requestId,
			"code":       http.StatusInternalServerError,
		})
		return
	}

	utils.WriteJSON(c, http.StatusCreated, gin.H{
		"message":    "Subscription successful",
		"request_id": requestId,
		"code":       http.StatusCreated,
	})

	go h.notifyUser(houseID, payload.Email)
}

func (h *Handler) notifyUser(houseID string, email string) {
	ctx := context.Background()
	sender := sender.New()

	message := fmt.Sprintf("New flats are available in house %s. Check them out now!", houseID)

	err := sender.SendEmail(ctx, email, message)
	if err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", email, err)
	}
}
