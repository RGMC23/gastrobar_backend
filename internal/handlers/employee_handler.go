package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gastrobar-backend/internal/models"
	"gastrobar-backend/internal/services"

	"github.com/gorilla/mux"
)

// EmployeeHandler maneja las solicitudes relacionadas con los empleados
type EmployeeHandler struct {
    employeeSvc services.EmployeeService
}

// NewEmployeeHandler crea una nueva instancia del manejador de empleados
func NewEmployeeHandler(employeeSvc services.EmployeeService) *EmployeeHandler {
    return &EmployeeHandler{
        employeeSvc: employeeSvc,
    }
}

// CreateEmployeeHandler maneja la solicitud para crear un nuevo empleado
func (h *EmployeeHandler) CreateEmployeeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var employee models.Employee
        if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        response, err := h.employeeSvc.CreateEmployee(employee)
        if err != nil {
            log.Printf("Error creating employee: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(response)
    }
}

// GetEmployeeHandler maneja la solicitud para obtener un empleado por ID
func (h *EmployeeHandler) GetEmployeeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        employeeIDStr := vars["id"]
        employeeID, err := strconv.Atoi(employeeIDStr)
        if err != nil {
            http.Error(w, "Invalid employee ID", http.StatusBadRequest)
            return
        }

        employee, err := h.employeeSvc.GetEmployee(employeeID)
        if err != nil {
            if strings.Contains(err.Error(), "employee not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Employee not found"})
                return
            }
            log.Printf("Error getting employee: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(employee)
    }
}

// ListEmployeesHandler maneja la solicitud para listar todos los empleados
func (h *EmployeeHandler) ListEmployeesHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        employees, err := h.employeeSvc.ListEmployees()
        if err != nil {
            log.Printf("Error listing employees: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(employees)
    }
}

// ListEmployeesByRoleHandler maneja la solicitud para listar empleados por rol (solo "empleado")
func (h *EmployeeHandler) ListEmployeesByRoleHandler(role models.EmployeeRole) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        employees, err := h.employeeSvc.ListEmployeesByRole(role)
        if err != nil {
            log.Printf("Error listing employees by role: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(employees)
    }
}

// UpdateEmployeeHandler maneja la solicitud para actualizar los datos de un empleado
func (h *EmployeeHandler) UpdateEmployeeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        employeeIDStr := vars["id"]
        employeeID, err := strconv.Atoi(employeeIDStr)
        if err != nil {
            http.Error(w, "Invalid employee ID", http.StatusBadRequest)
            return
        }

        var employee models.Employee
        if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        employee.ID = employeeID

        updatedEmployee, err := h.employeeSvc.UpdateEmployee(employee)
        if err != nil {
            if strings.Contains(err.Error(), "employee not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Employee not found"})
                return
            }
            log.Printf("Error updating employee: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedEmployee)
    }
}

// UpdateEmployeePasswordHandler maneja la solicitud para cambiar la contraseña de un empleado
func (h *EmployeeHandler) UpdateEmployeePasswordHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        employeeIDStr := vars["id"]
        employeeID, err := strconv.Atoi(employeeIDStr)
        if err != nil {
            http.Error(w, "Invalid employee ID", http.StatusBadRequest)
            return
        }

        // Generar una nueva contraseña aleatoria y actualizarla
        newPassword, err := h.employeeSvc.UpdateEmployeePassword(employeeID)
        if err != nil {
            if strings.Contains(err.Error(), "employee not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Employee not found"})
                return
            }
            log.Printf("Error updating employee password: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        // Devolver la nueva contraseña generada en la respuesta
        response := struct {
            Password string `json:"password"`
        }{Password: newPassword}

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(response)
    }
}