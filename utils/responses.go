package utils

import (
	"github.com/delapaska/avito-rent/models"
	"github.com/google/uuid"
)

// @Description Error response structure
// @Param message body string true "Error message" example "Invalid input"
// @Param request_id body string true "Unique request identifier" example "12345"
// @Param code body int true "HTTP status code" example 400
type ErrorResponse struct {

	// @example "Invalid input"
	Message any `json:"message"`
	// @example "12345"
	RequestID any `json:"request_id"`
	// @example 400
	Code int `json:"code"`
}

// @Description Successful login response structure
// @Param token body string true "JWT token" example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
type LoginResponse struct {
	// @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string `json:"token"`
}

// @Description Successful registration response structure
// @Param user_id body string true "UUID of the new user" example "a4b4a122-11c1-4b52-bd95-4a5d3c4be616"
type RegisterResponse struct {
	// @example "a4b4a122-11c1-4b52-bd95-4a5d3c4be616"
	UserID uuid.UUID `json:"user_id"`
}

// @Description Successful dummy login response structure
// @Param token body string true "JWT token" example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
type DummyLoginResponse struct {
	// @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	Token string `json:"token"`
}

// @Description Response model for retrieving flats in a house
// @Name FlatsResponse
// @Example { "flats": [{"id": "4e58c7c8-4b4e-44a3-a1e0-df9b3485cb46", "address": "123 Elm Street", "price": 1200, "size": 75, "available": true}, {"id": "8a55c7f2-93d5-4c77-8c67-46aeb5e9d2d8", "address": "456 Oak Avenue", "price": 1500, "size": 90, "available": false}] }
type FlatsResponse struct {
	Flats []models.Flat `json:"flats"`
}

// @Description Response model for subscription confirmation
// @Name SubscriptionResponse
// @Example { "message": "Subscription successful", "request_id": "12345", "code": 201 }
type SubscriptionResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}
