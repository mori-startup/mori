package handlers

import (
	"net/http"

	"mori/pkg/utils"
)

func (handler *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// 1) Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// 2) If it’s an OPTIONS preflight request, respond with allowed methods/headers
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// ... Now handle the GET request
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.RespondWithError(w, "Missing token", http.StatusBadRequest)
		return
	}

	err := handler.repos.UserRepo.VerifyEmail(token)
	if err != nil {
		utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondWithSuccess(w, "Email verified", http.StatusOK)
}
