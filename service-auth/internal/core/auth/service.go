// internal/core/auth/service.go
package auth

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type service struct {
	db        *firestore.Client
	jwtSecret []byte
}

func NewService(db *firestore.Client, jwtSecret []byte) Service {

	if len(jwtSecret) == 0 {
		if env := os.Getenv("JWT_SECRET"); env != "" {
			jwtSecret = []byte(env)
		}
	}

	return &service{db: db, jwtSecret: jwtSecret}
}

// User representa a estrutura de um usuário no Firestore.
type User struct {
	Username     string   `firestore:"username"`
	PasswordHash string   `firestore:"passwordHash"`
	Roles        []string `firestore:"roles"`
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	// 1. Encontrar o usuário no Firestore.
	query := s.db.Collection("users").Where("username", "==", username).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return "", errors.New("usuário ou senha inválidos")
	}
	if err != nil {
		log.Printf("Erro detalhado do Firestore: %v", err)
		return "", errors.New("erro ao consultar o banco de dados")
	}

	var user User
	if err := doc.DataTo(&user); err != nil {
		return "", errors.New("erro ao ler dados do usuário")
	}

	// 2. Comparar a senha fornecida com o hash armazenado.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("usuário ou senha inválidos")
	}

	// 3. Gerar o Token JWT com as permissões (roles).
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"roles":    user.Roles,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := claims.SignedString(s.jwtSecret)

	if err != nil {
		return "", errors.New("erro ao gerar token de acesso")
	}

	return tokenString, nil
}
