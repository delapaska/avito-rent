package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/delapaska/avito-rent/configs"
	"github.com/delapaska/avito-rent/models"
	"github.com/delapaska/avito-rent/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID, userType string) (string, error) {
	err := RoleGuard(userType)
	if err != nil {
		return "", err
	}

	claims := &models.Claims{
		UserID:   userID.String(),
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.Envs.JWTSecret))
}

func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId, _ := c.Get("RequestId")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("Retry-After", "30")
			utils.WriteJSON(c, http.StatusUnauthorized, gin.H{
				"message":    "Authorization header is required",
				"request_id": requestId,
				"code":       http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Envs.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.Header("Retry-After", "30")
			utils.WriteJSON(c, http.StatusUnauthorized, gin.H{
				"message":    "Invalid token",
				"request_id": requestId,
				"code":       http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if claims.UserType == role {
				userID, err := uuid.Parse(claims.UserID)
				if err != nil {
					utils.WriteJSON(c, http.StatusUnauthorized, gin.H{
						"message":    "Invalid user ID in token",
						"request_id": requestId,
						"code":       http.StatusUnauthorized,
					})
					c.Abort()
					return
				}
				c.Set("userID", userID)
				c.Set("userType", claims.UserType)
				c.Next()
				return
			}
		}

		utils.WriteJSON(c, http.StatusForbidden, gin.H{
			"message":    "Forbidden",
			"request_id": requestId,
			"code":       http.StatusForbidden,
		})
		c.Abort()
	}
}
