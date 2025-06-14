package server

import (
	"net/http"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	"contract_ease/internal/config"
)

func (s *Server) RegisterRoutes(tp trace.TracerProvider) http.Handler {

	r := gin.Default()
	r.Use(otelgin.Middleware("contract_ease", otelgin.WithTracerProvider(tp)))

	cfg := config.LoadConfig()
	frontendURL := cfg.App.FrontendURL

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/health", s.healthHandler)
		api.GET("/docs", s.apiReferenceHandler)
	}

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health(c))
}

func (s *Server) apiReferenceHandler(c *gin.Context) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./docs/build/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "ContractEase API",
		},
		Theme:    "deepSpace",
		Layout:   "modern",
		DarkMode: true,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API documentation"})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, htmlContent)
}
