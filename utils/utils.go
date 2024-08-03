package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(c *gin.Context, payload any) error {
	if c.Request == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(c.Request.Body).Decode(payload)
}

func WriteJSON(c *gin.Context, status int, v any) error {
	c.Header("Content-type", "application/json")
	c.Status(status)

	return json.NewEncoder(c.Writer).Encode(v)
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			message := fmt.Sprintf("field validation for '%s' failed on the '%s' tag", fieldError.Field(), fieldError.Tag())
			errors[field] = message
		}
	}
	return errors
}
