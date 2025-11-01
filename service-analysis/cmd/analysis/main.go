// cmd/analysis/main.go
package main

import (
	"log"

	"analysis-service/internal/api/handlers"
	"analysis-service/internal/api/responses"
	"analysis-service/internal/core/analysis"

	"github.com/gin-gonic/gin"
)

func main() {
	responses.InitLogger()

	analysisService := analysis.NewService()
	analysisHandler := handlers.NewAnalysisHandler(analysisService)

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		// Sem Middleware -- Gateway lida com isso
		apiV1.POST("/analyze/icms", analysisHandler.HandleAnalysisIcms)
		apiV1.POST("/analyze/ipi-st", analysisHandler.HandleAnalysisIpiSt)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "analysis-service"})
	})

	const port = "8082"
	log.Printf("🚀 Analysis Service (Go) iniciado e escutando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Falha ao iniciar o servidor de análise: ", err)
	}
}
