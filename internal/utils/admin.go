package utils

import (
	"context"
	"errors"
	"gin-app/config"
	"gin-app/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AUthenticate admin data
func GetLoggedInAdmin(c *gin.Context) (*models.User, int64, error) {
	accessToken, err := c.Cookie("admin_access")
	if err != nil || accessToken == "" {
		return nil, 0, errors.New("missing access token")
	}

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.App.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, 0, errors.New("invalid access token")
	}

	// Use jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, 0, errors.New("cannot parse claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, 0, errors.New("user_id not found in token claims")
	}
	adminID := int64(userIDFloat)

	var admin models.User
	err = config.DB.NewSelect().Model(&admin).
		// Column("id", "name", "email").  // if i need specific fields
		Where("id = ?", adminID).Scan(context.Background())
	if err != nil {
		return nil, adminID, err
	}

	return &admin, adminID, nil
}
