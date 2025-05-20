package handlers

import "mori/pkg/models"

// handler contains all repositories
type Handler struct {
	repos *models.Repositories
}

// initializing handler to return all repo connections
func InitHandlers(repos *models.Repositories) *Handler {
	return &Handler{repos: repos}
}
