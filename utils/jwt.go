package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

// JWTClaims defines the custom claims for the JWT token
type JWTClaims struct {
	UserID  uint   `json:"userId"`
	IsAdmin bool   `json:"isAdmin"`
	Email   string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for the user
func GenerateJWT(userID uint, email string, isAdmin bool) (string, error) {
	claims := JWTClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		Email:   email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("jwt.expirationTime")) * time.Hour).Unix(),
			Issuer:    viper.GetString("jwt.issuer"),
			Audience:  viper.GetString("jwt.audience"),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt.secretKey")))
}

// ParseJWT parses and validates the JWT token
func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secretKey")), nil
	})
	// Specific error handling for different token issues
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid signature")
		} else if v, ok := err.(*jwt.ValidationError); ok && v.Errors == jwt.ValidationErrorExpired {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
