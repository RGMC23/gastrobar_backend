package middleware

import (
	"context"
	"net/http"
	"strings"

	"gastrobar-backend/config"
	"gastrobar-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware valida el token JWT y verifica los roles permitidos
func AuthMiddleware(allowedRoles []models.EmployeeRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener el token del header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// Obtener jwtSecret desde la configuración
			jwtSecret := config.GetJWTSecret()

			// Parsear y validar el token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Obtener los claims del token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Verificar el rol del usuario
			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "Role not found in token", http.StatusUnauthorized)
				return
			}

			// Si se especificaron roles permitidos, verificar que el rol del usuario esté en la lista
			if len(allowedRoles) > 0 {
				roleAllowed := false
				for _, allowedRole := range allowedRoles {
					if role == string(allowedRole) {
						roleAllowed = true
						break
					}
				}
				if !roleAllowed {
					http.Error(w, "Insufficient permissions", http.StatusForbidden)
					return
				}
			}

			// Agregar el employee_id y el role al contexto para que las rutas puedan usarlos
			employeeID, _ := claims["employee_id"].(float64)
			ctx := context.WithValue(r.Context(), "employee_id", int(employeeID))
			ctx = context.WithValue(ctx, "role", role)

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
