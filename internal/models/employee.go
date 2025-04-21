package models

import "time"

// EmployeeRole define el tipo ENUM para los roles de los empleados
type EmployeeRole string

const (
    EmployeeRoleOwner    EmployeeRole = "dueño"
    EmployeeRoleAdmin    EmployeeRole = "administrador"
    EmployeeRoleEmployee EmployeeRole = "empleado"
)

// Employee representa la tabla employees
type Employee struct {
    ID           int          `json:"id"`
    EmployeeName string       `json:"employee_name"`
    Email        string       `json:"email"`
    PhoneNumber  string       `json:"phone_number"`
    Role         EmployeeRole `json:"role"`
    Username     string       `json:"username"`
    Password     string       `json:"-"` // No se incluye en JSON
    CreatedAt    time.Time    `json:"created_at"`
}

// EmployeeCreateResponse representa la respuesta al crear un empleado
type EmployeeCreateResponse struct {
    Employee Employee `json:"employee"`
    Password string   `json:"password"` // Contraseña generada
}