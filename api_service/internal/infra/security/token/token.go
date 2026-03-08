package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	secret []byte

	jwt.RegisteredClaims
}

func NewTokenService(secret string) *UserClaims {
	return &UserClaims{
		secret: []byte(secret),
	}
}

func (c *UserClaims) Generate(userID uuid.UUID) (string, error) {
	c.RegisteredClaims.Subject = userID.String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(c.secret)
}

func (c *UserClaims) Validate(tokenString string) (userID uuid.UUID, err error) {
	token, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		return c.secret, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		subject, ok := (*claims)["sub"].(string)
		if !ok {
			return uuid.Nil, jwt.ErrInvalidKeyType
		}
		return uuid.Parse(subject)
	}

	return uuid.Nil, jwt.ErrInvalidKey
}
