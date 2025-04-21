package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/pkg/errors"
    "gastrobar-backend/internal/models"
    "gastrobar-backend/internal/services"
)

// AuthHandler maneja las solicitudes relacionadas con la autenticación
type AuthHandler struct {
    authSvc services.AuthService
}

// NewAuthHandler crea una nueva instancia del manejador de autenticación
func NewAuthHandler(authSvc services.AuthService) *AuthHandler {
    return &AuthHandler{
        authSvc: authSvc,
    }
}

// LoginHandler maneja la solicitud de inicio de sesión
func (h *AuthHandler) LoginHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var loginReq models.LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        loginResp, err := h.authSvc.Login(loginReq)
        if err != nil {
            if errors.Is(err, errors.New("invalid email or password")) {
                http.Error(w, "Invalid email or password", http.StatusUnauthorized)
                return
            }
            log.Printf("Error during login: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(loginResp)
    }
}