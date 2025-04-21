package services

import (
	"time"

	"gastrobar-backend/config"
	"gastrobar-backend/internal/models"
	"gastrobar-backend/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthService define las operaciones relacionadas con la autenticación
type AuthService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateJWT(employee models.Employee) (string, error)
	Login(loginReq models.LoginRequest) (models.LoginResponse, error)
}

type authService struct {
	employeeRepo repositories.EmployeeRepository
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(employeeRepo repositories.EmployeeRepository) AuthService {
	return &authService{
		employeeRepo: employeeRepo,
	}
}

// HashPassword genera un hash de la contraseña usando bcrypt
func (s *authService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}
	return string(bytes), nil
}

// CheckPasswordHash verifica si la contraseña coincide con el hash
func (s *authService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT genera un token JWT para un empleado
func (s *authService) GenerateJWT(employee models.Employee) (string, error) {
	// Obtener jwtSecret y jwtExpire desde la configuración
	jwtSecret := config.GetJWTSecret()
	jwtExpire := config.GetJWTExpire()

	// Convertir jwtExpire a una duración (por ejemplo, "100h" -> 100 horas)
	expireDuration, err := time.ParseDuration(jwtExpire)
	if err != nil {
		return "", errors.Wrap(err, "invalid jwtExpire format: must be a valid duration (e.g., '24h')")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"employee_id": employee.ID,
		"role":        string(employee.Role),
		"exp":         time.Now().Add(expireDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.Wrap(err, "failed to generate JWT")
	}
	return tokenString, nil
}

// Login autentica a un empleado y genera un token JWT
func (s *authService) Login(loginReq models.LoginRequest) (models.LoginResponse, error) {
	// Buscar al empleado por username usando el repositorio
	employee, err := s.employeeRepo.FindByUsername(loginReq.Username)
	if err != nil {
		return models.LoginResponse{}, errors.Wrap(err, "failed to authenticate employee")
	}

	// Verificar la contraseña
	if !s.CheckPasswordHash(loginReq.Password, employee.Password) {
		return models.LoginResponse{}, errors.New("invalid email or password")
	}

	// Generar el token JWT
	token, err := s.GenerateJWT(employee)
	if err != nil {
		return models.LoginResponse{}, errors.Wrap(err, "failed to generate JWT")
	}

	return models.LoginResponse{Token: token}, nil
}
