// backend/utils/jwt_utils.go
package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/your-module/backend/config"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	TenantID uint   `json:"tenant_id"`
	RoleID   uint   `json:"role_id"`
	Username string `json:"username"`
	IssuedAt int64  `json:"issued_at"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, tenantID, roleID uint, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		RoleID:   roleID,
		Username: username,
		IssuedAt: time.Now().Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	cfg := config.GetConfig()
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
} 