package models

// LoginRequest y LoginResponse para el proceso de autenticación
type LoginRequest struct {
    Username    string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
}