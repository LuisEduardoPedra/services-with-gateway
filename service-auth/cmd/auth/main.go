// cmd/auth/main.go
package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"

	"auth-service/internal/api/handlers"
	"auth-service/internal/api/responses"
	"auth-service/internal/core/auth"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// --- Helper Functions ---

func initFirestoreClient(ctx context.Context) *firestore.Client {
	projectID := "analise-sped-db"
	databaseID := "analise-sped-db"
	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		log.Fatalf("Erro ao inicializar cliente Firestore: %v\n", err)
	}
	log.Printf("Conectado com sucesso ao Firestore")
	return client
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		if os.IsNotExist(err) {
			log.Print("Arquivo .env não encontrado, prosseguindo com variáveis de ambiente")
		} else {
			log.Printf("Erro ao carregar .env: %v", err)
		}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
	log.Print("Variáveis de ambiente carregadas de .env")
}

// --- Main Service Runner ---
func main() {
	loadEnv()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("FATAL: Variável de ambiente JWT_SECRET não está configurada.")
	}

	responses.InitLogger()
	ctx := context.Background()
	firestoreClient := initFirestoreClient(ctx)
	defer firestoreClient.Close()

	authService := auth.NewService(firestoreClient, []byte(jwtSecret))
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/login", authHandler.Login)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "service": "auth-service"})
	})

	// Run on different port
	const port = "8081"
	log.Printf("🚀 Auth Service (Go) iniciado e escutando na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Falha ao iniciar o servidor de autenticação: ", err)
	}
}
