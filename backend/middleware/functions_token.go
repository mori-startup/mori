package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var accessSecret = os.Getenv("ACCESS_SECRET_KEY_LLM")
var refreshSecret = os.Getenv("REFRESH_SECRET_KEY_LLM")

type CustomClaims struct {
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	jwt.RegisteredClaims
}

// Fonction pour générer un JWT
func GenerateJWT(userID, conversationID string) (string, error) {
	// Définir les claims
	claims := CustomClaims{
		UserID:         userID,
		ConversationID: conversationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 7)), // Expire dans 1 heure
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mori",
			Subject:   userID,
		},
	}

	// Créer le token avec les claims et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString(accessSecret)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération du token : %v", err)
	}

	return tokenString, nil
}

func GenerateRefreshJWT(userID, conversationID string) (string, error) {
	// Définir les claims
	claims := CustomClaims{
		UserID:         userID,
		ConversationID: conversationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Expire dans 1 heure
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mori",
			Subject:   userID,
		},
	}

	// Créer le token avec les claims et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString(refreshSecret)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération du token : %v", err)
	}

	return tokenString, nil
}

func SendRequestWithToken(url string, token string, jsonData []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête :", err)
		return
	}

	// Ajouter l'en-tête Authorization avec le token JWT
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête :", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Statut de la réponse :", resp.Status)
}
