package services

import (
    "crypto/rand"
    "encoding/base64"
    "gastrobar-backend/internal/models"
    "gastrobar-backend/internal/repositories"

    "github.com/pkg/errors"
    "golang.org/x/crypto/bcrypt"
)

// EmployeeService define las operaciones relacionadas con los empleados
type EmployeeService interface {
    CreateEmployee(employee models.Employee) (models.EmployeeCreateResponse, error)
    GetEmployee(employeeID int) (models.Employee, error)
    ListEmployees() ([]models.Employee, error)
    ListEmployeesByRole(role models.EmployeeRole) ([]models.Employee, error) // Nuevo
    UpdateEmployee(employee models.Employee) (models.Employee, error)
    UpdateEmployeePassword(employeeID int) (string, error)
}

type employeeService struct {
    employeeRepo repositories.EmployeeRepository
}

// NewEmployeeService crea una nueva instancia del servicio de empleados
func NewEmployeeService(employeeRepo repositories.EmployeeRepository) EmployeeService {
    return &employeeService{
        employeeRepo: employeeRepo,
    }
}

// generateRandomPassword genera una contraseña aleatoria segura
func generateRandomPassword(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", errors.Wrap(err, "failed to generate random password")
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// CreateEmployee crea un nuevo empleado con una contraseña aleatoria
func (s *employeeService) CreateEmployee(employee models.Employee) (models.EmployeeCreateResponse, error) {
    // Generar una contraseña aleatoria de 12 caracteres
    password, err := generateRandomPassword(12)
    if err != nil {
        return models.EmployeeCreateResponse{}, errors.Wrap(err, "failed to generate password")
    }

    // Encriptar la contraseña con bcrypt
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return models.EmployeeCreateResponse{}, errors.Wrap(err, "failed to hash password")
    }
    employee.Password = string(passwordHash)

    // Crear el empleado en el repositorio
    createdEmployee, err := s.employeeRepo.Create(employee)
    if err != nil {
        return models.EmployeeCreateResponse{}, errors.Wrap(err, "failed to create employee")
    }

    return models.EmployeeCreateResponse{
        Employee: createdEmployee,
        Password: password,
    }, nil
}

// GetEmployee obtiene un empleado por ID
func (s *employeeService) GetEmployee(employeeID int) (models.Employee, error) {
    employee, err := s.employeeRepo.FindByID(employeeID)
    if err != nil {
        return models.Employee{}, errors.Wrap(err, "failed to get employee")
    }
    return employee, nil
}

// ListEmployees lista todos los empleados
func (s *employeeService) ListEmployees() ([]models.Employee, error) {
    employees, err := s.employeeRepo.FindAll()
    if err != nil {
        return nil, errors.Wrap(err, "failed to list employees")
    }
    return employees, nil
}

// ListEmployeesByRole lista los empleados con un rol específico
func (s *employeeService) ListEmployeesByRole(role models.EmployeeRole) ([]models.Employee, error) {
    employees, err := s.employeeRepo.FindAllByRole(role)
    if err != nil {
        return nil, errors.Wrap(err, "failed to list employees by role")
    }
    return employees, nil
}

// UpdateEmployee actualiza los datos de un empleado (sin contraseña)
func (s *employeeService) UpdateEmployee(employee models.Employee) (models.Employee, error) {
    updatedEmployee, err := s.employeeRepo.Update(employee)
    if err != nil {
        return models.Employee{}, errors.Wrap(err, "failed to update employee")
    }
    return updatedEmployee, nil
}

// UpdateEmployeePassword cambia la contraseña de un empleado y devuelve la nueva contraseña generada
func (s *employeeService) UpdateEmployeePassword(employeeID int) (string, error) {
    // Generar una contraseña aleatoria de 12 caracteres
    newPassword, err := generateRandomPassword(12)
    if err != nil {
        return "", errors.Wrap(err, "failed to generate new password")
    }

    // Encriptar la nueva contraseña con bcrypt
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return "", errors.Wrap(err, "failed to hash new password")
    }

    // Actualizar la contraseña en el repositorio
    err = s.employeeRepo.UpdatePassword(employeeID, string(passwordHash))
    if err != nil {
        return "", errors.Wrap(err, "failed to update employee password")
    }

    return newPassword, nil
}