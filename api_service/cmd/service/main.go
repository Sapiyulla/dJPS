package main

import (
	userService "api-service/internal/application/user"
	userRepository "api-service/internal/infra/repository/user"
	"api-service/internal/infra/security/token"
	userHandler "api-service/internal/interfaces/http/user"
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type AppMode string

const (
	Dev  AppMode = "dev"
	Prod AppMode = "prod"
)

var (
	ENV_TOKEN_SECRET string
	ENV_MODE         AppMode
)

func init() {
	ENV_MODE = AppMode(os.Getenv("MODE"))
	if ENV_MODE == AppMode("") {
		ENV_MODE = Dev
	}

	ENV_TOKEN_SECRET = os.Getenv("TOKEN_SECRET")
	if ENV_TOKEN_SECRET == "" {
		slog.Warn("env load error")
		os.Exit(1)
	}
}

func main() {
	tokenService := token.NewTokenService(ENV_TOKEN_SECRET)

	// depends
	userRepo := userRepository.NewUserLocalStorage()
	userService := userService.NewUserService(userRepo, tokenService)

	userHandler := userHandler.NewUserHandler(userService, 5*time.Second)

	// http serve
	r := gin.Default()
	api := r.Group("/api")
	{
		userApi := api.Group("/user")
		{
			userApi.POST("/registration", userHandler.RegistrationHandler)
			userApi.POST("/signin", userHandler.SignInHandler)
		}
	}

	addr := ":"
	if ENV_MODE == Dev {
		addr += "8001"
	} else if ENV_MODE == Prod {
		addr += "443"
	}

	httpServer := http.Server{
		Addr:    addr,
		Handler: r,
	}

	wg := sync.WaitGroup{}

	wg.Go(func() {
		switch ENV_MODE {
		case Dev:
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("server error", "error", err)
			}
		case Prod:
			if err := httpServer.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
				slog.Error("server error", "error", err)
			}
		}
	})
	slog.Info("http server started", "mode", ENV_MODE, "addr", addr)

	wg.Wait()
	_ = httpServer.Shutdown(context.Background())
	slog.Info("http server stopped")
}
