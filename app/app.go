package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ximofam/chess-online/config"
	"github.com/ximofam/chess-online/database"
	Auth "github.com/ximofam/chess-online/services/auth"
	"github.com/ximofam/chess-online/services/user"
	"github.com/ximofam/chess-online/services/ws"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB
	r  *gin.Engine
}

func New(db *gorm.DB) *App {
	r := gin.Default()

	return &App{
		db: db,
		r:  r,
	}
}

func (a *App) Run() {
	a.registerRoutes()

	sv := http.Server{
		Addr:    config.Envs.Server.Addr,
		Handler: a.r,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("Listening on port: %s", sv.Addr)

		if err := sv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Printf("Server error: %v", err)

	case sig := <-sigCh:
		log.Printf("Received signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := sv.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}

	log.Println("Server stopped")
}

func (a *App) registerRoutes() {
	a.r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	gormDatabase := &database.GormDatabase{DB: a.db}

	userHandler := user.NewHandler(gormDatabase)

	g := a.r.Group("/api/v1")

	a.r.GET("/ws", Auth.AuthenticateByPathValue(), ws.Server.ServeWS)

	auth := g.Group("/auth")
	{
		auth.GET("/me", Auth.Authenticate(), userHandler.GetMe)
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)
	}

	user := g.Group("/users", Auth.Authenticate())
	{
		user.GET("/:id", userHandler.GetUser)
	}
}
