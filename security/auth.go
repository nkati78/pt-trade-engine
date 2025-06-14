package security

import (
	"github.com/golang-jwt/jwt"
)

type Auth struct {
	secret []byte
}

// GenerateToken generates a new jwt token
func (a Auth) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = userID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// NewAuth creates a new Auth instance
func NewAuth(secret []byte) Auth {
	return Auth{secret: secret}
}

// ValidateToken validates a jwt token
func (a Auth) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", err
	}

	userID := claims["sub"].(string)
	return userID, nil
}
