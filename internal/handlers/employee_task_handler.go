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

// EmployeeTaskHandler maneja las solicitudes relacionadas con las tareas de los empleados
type EmployeeTaskHandler struct {
    employeeTaskSvc services.EmployeeTaskService
}

// NewEmployeeTaskHandler crea una nueva instancia del manejador de tareas
func NewEmployeeTaskHandler(employeeTaskSvc services.EmployeeTaskService) *EmployeeTaskHandler {
    return &EmployeeTaskHandler{
        employeeTaskSvc: employeeTaskSvc,
    }
}

// CreateTaskHandler maneja la solicitud para crear una nueva tarea
func (h *EmployeeTaskHandler) CreateTaskHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var task models.EmployeeTask
        if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        createdTask, err := h.employeeTaskSvc.CreateTask(task)
        if err != nil {
            // Verificar si el error es debido a la validación del rol del empleado
            if strings.Contains(err.Error(), "tasks can only be assigned to employees with role 'empleado'") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Tasks can only be assigned to employees with role 'empleado'"})
                return
            }
            // Verificar si el error es porque el empleado no existe
            if strings.Contains(err.Error(), "employee does not exist") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Employee does not exist"})
                return
            }
            // Para otros errores inesperados, mantener el código 500
            log.Printf("Error creating task: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdTask)
    }
}

// GetTaskHandler maneja la solicitud para obtener una tarea por ID
func (h *EmployeeTaskHandler) GetTaskHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        taskIDStr := vars["id"]
        taskID, err := strconv.Atoi(taskIDStr)
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        task, err := h.employeeTaskSvc.GetTask(taskID)
        if err != nil {
            if strings.Contains(err.Error(), "task not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
                return
            }
            log.Printf("Error getting task: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(task)
    }
}

// ListTasksHandler maneja la solicitud para listar todas las tareas (admin y dueño)
func (h *EmployeeTaskHandler) ListTasksHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tasks, err := h.employeeTaskSvc.ListTasks()
        if err != nil {
            log.Printf("Error listing tasks: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tasks)
    }
}

// ListTasksByEmployeeHandler maneja la solicitud para listar las tareas de un empleado (para empleados)
func (h *EmployeeTaskHandler) ListTasksByEmployeeHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Obtener el employee_id del token JWT (del contexto, seteado por el middleware)
        employeeID, ok := r.Context().Value("employee_id").(int)
        if !ok {
            http.Error(w, "Invalid employee ID in token", http.StatusUnauthorized)
            return
        }

        tasks, err := h.employeeTaskSvc.ListTasksByEmployee(employeeID)
        if err != nil {
            log.Printf("Error listing tasks for employee: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tasks)
    }
}

// UpdateTaskHandler maneja la solicitud para actualizar una tarea (admin y dueño)
func (h *EmployeeTaskHandler) UpdateTaskHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        taskIDStr := vars["id"]
        taskID, err := strconv.Atoi(taskIDStr)
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        var task models.EmployeeTask
        if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        task.ID = taskID

        updatedTask, err := h.employeeTaskSvc.UpdateTask(task)
        if err != nil {
            // Verificar si el error es debido a la validación del rol del empleado
            if strings.Contains(err.Error(), "tasks can only be assigned to employees with role 'empleado'") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Tasks can only be assigned to employees with role 'empleado'"})
                return
            }
            // Verificar si el error es porque el empleado no existe
            if strings.Contains(err.Error(), "employee does not exist") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Employee does not exist"})
                return
            }
            // Verificar si el error es porque la tarea no existe
            if strings.Contains(err.Error(), "task not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
                return
            }
            // Para otros errores inesperados, mantener el código 500
            log.Printf("Error updating task: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedTask)
    }
}

// UpdateTaskStatusHandler maneja la solicitud para actualizar el estado de una tarea (admin y dueño)
func (h *EmployeeTaskHandler) UpdateTaskStatusHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        taskIDStr := vars["id"]
        taskID, err := strconv.Atoi(taskIDStr)
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        var request struct {
            Status models.EmployeeTaskStatus `json:"status"`
        }
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validar el estado
        if request.Status != models.TaskStatusPending && request.Status != models.TaskStatusInProgress && request.Status != models.TaskStatusCompleted {
            http.Error(w, "Invalid task status", http.StatusBadRequest)
            return
        }

        err = h.employeeTaskSvc.UpdateTaskStatus(taskID, request.Status)
        if err != nil {
            if strings.Contains(err.Error(), "task not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
                return
            }
            log.Printf("Error updating task status: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Task status updated successfully"))
    }
}

// DeleteTaskHandler maneja la solicitud para eliminar una tarea (admin y dueño)
func (h *EmployeeTaskHandler) DeleteTaskHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        taskIDStr := vars["id"]
        taskID, err := strconv.Atoi(taskIDStr)
        if err != nil {
            http.Error(w, "Invalid task ID", http.StatusBadRequest)
            return
        }

        err = h.employeeTaskSvc.DeleteTask(taskID)
        if err != nil {
            if strings.Contains(err.Error(), "task not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
                return
            }
            log.Printf("Error deleting task: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Task deleted successfully"))
    }
}