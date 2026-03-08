package ports

import (
	"github.com/google/uuid"
)

type TokenService interface {
	Generate(userID uuid.UUID) (string, error)
	Validate(tokenString string) (userID uuid.UUID, err error)
}
