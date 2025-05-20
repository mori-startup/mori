package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mori/pkg/models"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var accessSecret string
var refreshSecret string

// AllowedDomains is a list of domains that are allowed to be accessed
var AllowedDomains = []string{
	"127.0.0.1:3000", // Local development
	"localhost:3000", // Local development
	// Add your production domains here
}

// AllowedSchemes is a list of allowed URL schemes
var AllowedSchemes = []string{
	"http",
	"https",
}

// isPrivateIP checks if an IP address is in a private range
func isPrivateIP(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 || // 10.0.0.0/8
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
			(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
	}
	return false
}

// validateURL checks if the URL is allowed and not pointing to a private IP
func validateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Validate scheme
	schemeAllowed := false
	for _, scheme := range AllowedSchemes {
		if parsedURL.Scheme == scheme {
			schemeAllowed = true
			break
		}
	}
	if !schemeAllowed {
		return fmt.Errorf("scheme not allowed: %s", parsedURL.Scheme)
	}

	// Check if the host is in the allowed domains list
	hostAllowed := false
	for _, domain := range AllowedDomains {
		if parsedURL.Host == domain {
			hostAllowed = true
			break
		}
	}
	if !hostAllowed {
		return fmt.Errorf("domain not allowed: %s", parsedURL.Host)
	}

	// Resolve the hostname to check for private IPs
	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return fmt.Errorf("failed to resolve hostname: %v", err)
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("private IP addresses are not allowed")
		}
	}

	// Validate port if present
	if parsedURL.Port() != "" {
		port := parsedURL.Port()
		if port != "80" && port != "443" && port != "3000" {
			return fmt.Errorf("port not allowed: %s", port)
		}
	}

	return nil
}

type CustomClaims struct {
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
	// History        []struct {
	// 	UserRequest string `json:"user_request"`
	// 	LLMResponse string `json:"llm_response"`
	// } `json:"history"`
	// Stop           bool   `json:"stop"`

	jwt.RegisteredClaims
}

// type Conversation struct {
// 	UserID          string   `json:"user_id"`
// 	ConversationID  string   `json:"conversation_id"`
// 	Session         string   `json:"session"`
// 	UserRequest     string   `json:"user_request"`
// 	LLMResponse     string   `json:"llm_response"`
// 	NewConversation bool     `json:"new_conversation"`
// 	CreatedAt       string   `json:"created_at"`
// 	UpdateAt        string   `json:"update_at"`
// 	History         []string `json:"history"`
// }

type ServerPython struct {
	Status         string `json:"status"`
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Response       string `json:"response"`
	Timestamp      string `json:"timestamp"`
}

type ConversationResponse struct {
	Type string              `json:"type"`
	Data models.Conversation `json:"data"`
}

func (handler *Handler) LLMHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMHandler"+err.Error(),
				http.StatusInternalServerError)
			return
		}
		fmt.Println("conversation", conversation)

		userRequest := conversation.Convo[len(conversation.Convo)-1]["user_request"]

		accessToken, err := GenerateJWT(conversation.UserID, conversation.ConversationID, userRequest)
		if err != nil {
			http.Error(w, "Error generating JWT",
				http.StatusInternalServerError)
			fmt.Println("Error generating JWT: ", err)
			return

		}

		refreshToken, err := GenerateRefreshJWT(conversation.UserID, conversation.ConversationID, userRequest)
		if err != nil {
			http.Error(w, "Error generating refresh JWT",
				http.StatusInternalServerError)
			fmt.Println("Error generating refresh JWT: ", err)
			return

		}

		//put token in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "accessToken",
			Value:    accessToken,
			Path:     "http://127.0.0.1:3000/llm-protected",
			HttpOnly: true,
			Secure:   true, // Activez HTTPS en production
			SameSite: http.SameSiteStrictMode,
			MaxAge:   7 * 60, // 7 minutes
		})

		// println("JWT Token: ", accessToken)

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			Path:     "http://127.0.0.1:3000/llm-protected",
			HttpOnly: true,
			Secure:   true, // Activez HTTPS en production
			SameSite: http.SameSiteStrictMode,
			MaxAge:   7 * 24 * 60 * 60, // 7 jours
		})

		llmConversation := map[string]interface{}{
			"user_id":         conversation.UserID,
			"conversation_id": conversation.ConversationID,
			"message":         userRequest,
			// "history":         conversation.History,
		}

		data, err := json.Marshal(llmConversation)
		if err != nil {
			http.Error(w, "Error marshalling JSON",
				http.StatusInternalServerError)
			return

		}

		SendRequestWithToken("http://127.0.0.1:3000/llm-protected", accessToken, data, w)
		return
	}
}

func SendRequestWithToken(urlStr string, token string, jsonData []byte, w http.ResponseWriter) {
	// Validate the URL before making the request
	if err := validateURL(urlStr); err != nil {
		http.Error(w, fmt.Sprintf("URL validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// Sanitize and parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid URL format: %v", err), http.StatusBadRequest)
		return
	}

	// Create a POST request with the sanitized URL
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the URL components directly
	req.URL = &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   parsedURL.Path,
	}

	// Add the Authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a client with timeout and disabled redirects
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Disable redirects
		},
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: -1, // Disable keep-alive
			}).DialContext,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Validate response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Unexpected status code: %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		// Basic validation - check if content is not empty and has reasonable length
		if len(line) == 0 /*|| len(line) > 2048*2048*/ { // 1MB max
			http.Error(w, "Invalid response content length", http.StatusBadGateway)
			return
		}

		// Check for binary content
		if !utf8.Valid(line) {
			http.Error(w, "Invalid response content encoding", http.StatusBadGateway)
			return
		}

		// Send each chunk to the frontend
		fmt.Fprintf(w, "%s", line)
		flusher.Flush() // Send data immediately to the client
	}

	fmt.Println("Response status stream:", resp.Status)
}

// isValidResponseContent validates the response content before sending it to the client
func isValidResponseContent(content []byte) bool {
	// Basic validation - check if content is not empty and has reasonable length
	if len(content) == 0 || len(content) > 1024*1024 { // 1MB max
		return false
	}

	// Check for binary content
	if !utf8.Valid(content) {
		return false
	}

	// Try to parse as JSON to validate structure
	var jsonData map[string]interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		// If it's not JSON, it might be a text response
		text := string(content)

		// Check for common SSRF attack patterns
		blockedPatterns := []string{
			"<?xml",
			"<!DOCTYPE",
			"<html",
			"<script",
			"<?php",
			"<?=",
			"<? ",
			"<?\n",
			"<?\r",
			"<?\t",
			"<? ",
			"<?\f",
			"<?\v",
		}

		for _, pattern := range blockedPatterns {
			if strings.Contains(strings.ToLower(text), strings.ToLower(pattern)) {
				return false
			}
		}

		return true
	}

	// For JSON responses, do basic validation
	// Only check if it's a valid JSON object
	return true
}

// Fonction pour générer un JWT
func GenerateJWT(username, conversationID, message string) (string, error) {
	// Définir les claims
	err := godotenv.Load()
	if err != nil {
		log.Printf("Erreur lors du chargement du fichier .env : %v", err)
	}
	accessSecret = os.Getenv("ACCESS_SECRET_KEY_LLM")
	claims := CustomClaims{
		UserID:         username,
		ConversationID: conversationID,
		Message:        message,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 7)), // Expire dans 1 heure
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mori",
			Subject:   username,
		},
	}

	// Créer le token avec les claims et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey := []byte(accessSecret)
	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return tokenString, nil
}

func GenerateRefreshJWT(username, conversationID, message string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Erreur lors du chargement du fichier .env : %v", err)
	}
	refreshSecret = os.Getenv("REFRESH_SECRET_KEY_LLM")

	// Définir les claims
	claims := CustomClaims{
		UserID:         username,
		ConversationID: conversationID,
		Message:        message,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Expire dans 1 heure
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mori",
			Subject:   username,
		},
	}

	// Créer le token avec les claims et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération du token : %v", err)
	}

	return tokenString, nil
}

// func SendRequestWithToken(url string, token string, jsonData []byte) {
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("Erreur lors de la création de la requête :", err)
// 		return
// 	}

// 	// Ajouter l'en-tête Authorization avec le token JWT
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'envoi de la requête :", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	fmt.Println("Statut de la réponse :", resp.Status)
// }

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Lire le refresh token envoyé par le client
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Refresh token manquant", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	// Valider le refresh token
	claims := &CustomClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})

	if err != nil {
		http.Error(w, "Refresh token invalide ou expiré", http.StatusUnauthorized)
		return
	}

	// Générer un nouvel access token
	newAccessToken, err := GenerateJWT(claims.UserID, claims.ConversationID, claims.Message)
	if err != nil {
		http.Error(w, "Erreur lors de la génération du nouvel access token", http.StatusInternalServerError)
		return
	}

	// Retourner le nouvel access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"accessToken": newAccessToken,
	})
}

func VerifyAndRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les tokens depuis les cookies
	accessTokenCookie, err := r.Cookie("accessToken")
	if err != nil {
		http.Error(w, "Access token manquant", http.StatusUnauthorized)
		return
	}
	refreshTokenCookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Refresh token manquant", http.StatusUnauthorized)
		return
	}

	accessToken := accessTokenCookie.Value
	refreshToken := refreshTokenCookie.Value

	// Vérifier l'expiration de l'access token
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	// Si le token est valide, vérifier s'il est proche de l'expiration
	if err == nil && token.Valid {
		timeRemaining := time.Until(claims.ExpiresAt.Time)
		if timeRemaining > 2*time.Minute {
			// Token encore valide, pas besoin de le renouveler
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":      "valid",
				"accessToken": accessToken,
			})
			return
		}
	}

	// Si l'access token est expiré ou proche de l'expiration, vérifier le refresh token
	refreshClaims := &CustomClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		http.Error(w, "Refresh token invalide ou expiré", http.StatusUnauthorized)
		return
	}

	// Générer un nouveau access token
	newAccessToken, err := GenerateJWT(refreshClaims.UserID, refreshClaims.ConversationID, refreshClaims.Message)
	if err != nil {
		http.Error(w, "Erreur lors de la génération du nouveau token", http.StatusInternalServerError)
		return
	}

	// Retourner le nouveau access token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":      "refreshed",
		"accessToken": newAccessToken,
	})
}

// LLMConvoSave saves the conversation to the database
func (handler *Handler) LLMConvoSave(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoSave "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		if conversation.NewConversation == true {
			// conversation.ConversationID = uuid.NewV4().String()
			fmt.Println("New conversation", conversation)
			err = handler.repos.LLMConvoRepo.SaveConvo(conversation)
			if err != nil {
				http.Error(w, "Error saving new conversation: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = handler.repos.LLMConvoRepo.SaveConvo(conversation)
			if err != nil {
				http.Error(w, "Error saving conversation: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		//send conversation.ConversationID to the client
		send_convoID := map[string]string{"conversation_id": conversation.ConversationID}

		// data, err := json.Marshal(send_convoID)
		// if err != nil {
		// 	http.Error(w, "Error marshalling JSON LLMConvoSave", http.StatusInternalServerError)
		// 	return
		// }
		fmt.Println("Conversation ID: ", conversation.ConversationID)
		fmt.Println("Send conversation ID: ", send_convoID)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(send_convoID)
		if err != nil {
			http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
}

// functon that get all the previous conversation
func (handler *Handler) LLMConvoGetAll(w http.ResponseWriter, convo models.Conversation) {

	var conversations []models.Conversation

	fmt.Println("Convo: ", convo)
	conversations, err := handler.repos.LLMConvoRepo.GetAllConvo(convo)
	if err != nil {
		http.Error(w, "Error getting all conversations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("ALL PREVIOUS CONVERSATIONS : ", conversations)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// LLMConvoGet gets the conversation from the database
func (handler *Handler) LLMConvoGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoGet "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		handler.LLMConvoGetAll(w, conversation)

		return
	}
}

func (handler *Handler) LLMConvoGetLast(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request

		type RequestConvo struct {
			UserID string `json:"user_id"`
		}

		var requestConvo RequestConvo

		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body LLMConvoGetLast request: ", string(body))

		err = json.Unmarshal(body, &requestConvo)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoGet "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		conversations, err := handler.repos.LLMConvoRepo.GetLastConvo()
		if err != nil {
			http.Error(w, "Error getting last conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Last conversation: ", conversations)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(conversations)
		if err != nil {
			http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
}

// LLMConvoDelete deletes the conversation from the database
func (handler *Handler) LLMConvoDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoDelete "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		err = handler.repos.LLMConvoRepo.DeleteConvo(conversation)
		if err != nil {
			http.Error(w, "Error deleting conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Conversation deleted")

		return
	}
}

// LLMConvoSelected gets the conversation from the database
func (handler *Handler) LLMConvoSelected(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoSelected "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		conversations, err := handler.repos.LLMConvoRepo.GetConvo(conversation)
		if err != nil {
			http.Error(w, "Error getting conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Conversation: ", conversations)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(conversations)
		if err != nil {
			http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
}
