package user

import (
	"api-service/internal/application/ports"
	"api-service/internal/application/user"
	"api-service/internal/interfaces/http/user/dto"
	"api-service/internal/pkg/errs"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *user.UserService

	requestTimeout time.Duration
}

func NewUserHandler(userService *user.UserService, requestTimeout time.Duration) *UserHandler {
	return &UserHandler{
		userService:    userService,
		requestTimeout: requestTimeout,
	}
}

func (h *UserHandler) RegistrationHandler(c *gin.Context) {
	var registrationDto dto.RegistrationDto
	if err := c.ShouldBindJSON(&registrationDto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	token, err := h.userService.Registration(ctx,
		registrationDto.Username,
		registrationDto.Email,
		registrationDto.Password)
	if err != nil {
		switch err {
		case ports.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, errs.DefaultError{Error: err.Error()})
			return
		case context.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, errs.DefaultError{Error: "request timeout"})
			return
		case ports.ErrInternal:
			c.JSON(http.StatusInternalServerError, errs.DefaultError{Error: err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, errs.DefaultError{Error: "internal error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) SignInHandler(c *gin.Context) {
	var registrationDto dto.SignInDto
	if err := c.ShouldBindJSON(&registrationDto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
	defer cancel()

	token, err := h.userService.SignIn(ctx,
		registrationDto.Email,
		registrationDto.Password)
	if err != nil {
		switch err {
		case ports.ErrUserNotFound:
			c.JSON(http.StatusNotFound, errs.DefaultError{Error: err.Error()})
			return
		case user.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, errs.DefaultError{Error: err.Error()})
			return
		case context.DeadlineExceeded:
			c.JSON(http.StatusGatewayTimeout, errs.DefaultError{Error: "request timeout"})
			return
		case ports.ErrInternal:
			c.JSON(http.StatusInternalServerError, errs.DefaultError{Error: err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, errs.DefaultError{Error: "internal error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
