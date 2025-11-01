// cmd/converter/main.go
package main

import (
	"log"

	"converter-service/internal/api/handlers"
	"converter-service/internal/api/responses"
	"converter-service/internal/core/converter"

	"github.com/gin-gonic/gin"
)

func main() {
	responses.InitLogger()

	converterService := converter.NewService()
	converterHandler := handlers.NewConverterHandler(converterService)

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/convert/francesinha", converterHandler.HandleSicrediConversion)
		apiV1.POST("/convert/receitas-acisa", converterHandler.HandleReceitasAcisaConversion)
		apiV1.POST("/convert/atolini-pagamentos", converterHandler.HandleAtoliniPagamentosConversion)
		apiV1.POST("/convert/atolini-recebimentos", converterHandler.HandleAtoliniRecebimentosConversion)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "converter-service"})
	})

	const port = "8083"
	log.Printf("🚀 Converter Service (Go) iniciado e escutando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Falha ao iniciar o servidor de conversão: ", err)
	}
}
