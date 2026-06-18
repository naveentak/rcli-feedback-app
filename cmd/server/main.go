package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rcli/feedback/internal/auth"
	"github.com/rcli/feedback/internal/config"
	"github.com/rcli/feedback/internal/feedback"
	ghclient "github.com/rcli/feedback/internal/github"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	gh := ghclient.NewClient(cfg.GitHubToken, cfg.GitHubOwner, cfg.GitHubRepo)
	svc := feedback.NewService(gh)
	handler := feedback.NewHandler(svc)

	r := gin.Default()
	r.Use(corsMiddleware(cfg.AllowedOrigins))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// User-facing page
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"apps": feedback.ValidApps,
		})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/feedback", handler.List)
		api.GET("/feedback/:number", handler.Get)

		if cfg.DevMode || cfg.PublicSubmit {
			api.POST("/feedback", handler.Submit)
		}

		if cfg.HMACSecret != "" {
			api.POST("/feedback/signed", handler.SignedSubmit(cfg.HMACSecret))
		}

		if len(cfg.APIKeys) > 0 {
			protected := api.Group("")
			protected.Use(auth.APIKeyMiddleware(cfg.APIKeys))
			if !cfg.DevMode && !cfg.PublicSubmit {
				protected.POST("/feedback", handler.Submit)
			}
			protected.POST("/feedback/:number/comments", handler.Comment)
			protected.PATCH("/feedback/:number", handler.UpdateStatus)
		} else if !cfg.DevMode {
			api.POST("/feedback", handler.Submit)
			api.POST("/feedback/:number/comments", handler.Comment)
			api.PATCH("/feedback/:number", handler.UpdateStatus)
		}
	}

	log.Printf("feedback service listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	defaults := []string{
		"https://rclip.refactory.co.za",
		"http://localhost:5173",
		"http://localhost:3000",
	}
	origins := append(defaults, allowedOrigins...)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		for _, o := range origins {
			if origin == o {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, X-App, Authorization")
				break
			}
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}