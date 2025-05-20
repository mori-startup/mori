package middleware

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
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

// List of allowed domains
var allowedDomains = []string{
	"localhost",
	"llm-service.yourdomain.com",
	// Add other trusted domains here
}

// isAllowedURL checks if the URL is allowed based on domain and IP restrictions
func isAllowedURL(urlStr string) (bool, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false, err
	}

	// Check protocol - only allow HTTP and HTTPS
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false, fmt.Errorf("only HTTP and HTTPS protocols are allowed")
	}

	// Check if domain is in the allowed list
	hostname := parsedURL.Hostname()
	for _, domain := range allowedDomains {
		if hostname == domain || strings.HasSuffix(hostname, "."+domain) {
			return true, nil
		}
	}

	// Check if the hostname resolves to a private IP
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return false, err
	}

	for _, ip := range ips {
		// Check if IP is private
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return false, fmt.Errorf("requests to private IP addresses are not allowed")
		}
	}

	return false, fmt.Errorf("domain not in allowed list: %s", hostname)
}

func SendRequestWithToken(url string, token string, jsonData []byte) error {
	// Validate URL before making the request
	allowed, err := isAllowedURL(url)
	if err != nil || !allowed {
		if err != nil {
			return fmt.Errorf("URL validation failed: %v", err)
		}
		return fmt.Errorf("URL is not allowed: %s", url)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Add the Authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Create a custom client with disabled redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Prevent redirects
		},
		Timeout: 30 * time.Second, // Set a reasonable timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Don't return the response body to the caller
	// Only return status information
	fmt.Println("Response status:", resp.Status)

	return nil
}
