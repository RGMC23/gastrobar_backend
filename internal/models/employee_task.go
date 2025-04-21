package models

import "time"

// EmployeeTaskStatus define el tipo ENUM para los estados de las tareas
type EmployeeTaskStatus string

const (
    TaskStatusPending    EmployeeTaskStatus = "pendiente"
    TaskStatusInProgress EmployeeTaskStatus = "en_progreso"
    TaskStatusCompleted  EmployeeTaskStatus = "completada"
)

// EmployeeTask representa la tabla employee_tasks
type EmployeeTask struct {
    ID              int               `json:"id"`
    EmployeeID      int               `json:"employee_id"`
    TaskDescription string            `json:"task_description"`
    Status          EmployeeTaskStatus `json:"status"`
    AssignedAt      time.Time         `json:"assigned_at"`
    CompletedAt     *time.Time        `json:"completed_at"` // Puede ser NULL, usamos un puntero
    CreatedAt       time.Time         `json:"created_at"`
}