package services

import (
    "time"

    "gastrobar-backend/internal/models"
    "gastrobar-backend/internal/repositories"

    "github.com/pkg/errors"
)

// EmployeeTaskService define las operaciones relacionadas con las tareas de los empleados
type EmployeeTaskService interface {
    CreateTask(task models.EmployeeTask) (models.EmployeeTask, error)
    GetTask(taskID int) (models.EmployeeTask, error)
    ListTasks() ([]models.EmployeeTask, error)
    ListTasksByEmployee(employeeID int) ([]models.EmployeeTask, error)
    UpdateTask(task models.EmployeeTask) (models.EmployeeTask, error)
    UpdateTaskStatus(taskID int, status models.EmployeeTaskStatus) error
    DeleteTask(taskID int) error
}

type employeeTaskService struct {
    employeeTaskRepo repositories.EmployeeTaskRepository
    employeeRepo     repositories.EmployeeRepository
}

// NewEmployeeTaskService crea una nueva instancia del servicio de tareas
func NewEmployeeTaskService(employeeTaskRepo repositories.EmployeeTaskRepository, employeeRepo repositories.EmployeeRepository) EmployeeTaskService {
    return &employeeTaskService{
        employeeTaskRepo: employeeTaskRepo,
        employeeRepo:     employeeRepo,
    }
}

// CreateTask crea una nueva tarea para un empleado
func (s *employeeTaskService) CreateTask(task models.EmployeeTask) (models.EmployeeTask, error) {
    // Validar que el empleado exista
    employee, err := s.employeeRepo.FindByID(task.EmployeeID)
    if err != nil {
        return models.EmployeeTask{}, errors.Wrap(err, "employee does not exist")
    }

    // Validar que el empleado tenga el rol "empleado"
    if employee.Role != models.EmployeeRoleEmployee {
        return models.EmployeeTask{}, errors.New("tasks can only be assigned to employees with role 'empleado'")
    }

    // Establecer valores predeterminados
    task.Status = models.TaskStatusPending
    task.AssignedAt = time.Now()
    task.CompletedAt = nil // Inicialmente no completada

    // Crear la tarea en el repositorio
    createdTask, err := s.employeeTaskRepo.Create(task)
    if err != nil {
        return models.EmployeeTask{}, errors.Wrap(err, "failed to create task")
    }

    return createdTask, nil
}

// GetTask obtiene una tarea por ID
func (s *employeeTaskService) GetTask(taskID int) (models.EmployeeTask, error) {
    task, err := s.employeeTaskRepo.FindByID(taskID)
    if err != nil {
        return models.EmployeeTask{}, errors.Wrap(err, "failed to get task")
    }
    return task, nil
}

// ListTasks lista todas las tareas (para admin y dueño)
func (s *employeeTaskService) ListTasks() ([]models.EmployeeTask, error) {
    tasks, err := s.employeeTaskRepo.FindAll()
    if err != nil {
        return nil, errors.Wrap(err, "failed to list tasks")
    }
    return tasks, nil
}

// ListTasksByEmployee lista las tareas de un empleado específico (para empleados)
func (s *employeeTaskService) ListTasksByEmployee(employeeID int) ([]models.EmployeeTask, error) {
    tasks, err := s.employeeTaskRepo.FindByEmployeeID(employeeID)
    if err != nil {
        return nil, errors.Wrap(err, "failed to list tasks for employee")
    }
    return tasks, nil
}

// UpdateTask actualiza una tarea (para admin y dueño)
func (s *employeeTaskService) UpdateTask(task models.EmployeeTask) (models.EmployeeTask, error) {
    // Validar que el empleado exista
    employee, err := s.employeeRepo.FindByID(task.EmployeeID)
    if err != nil {
        return models.EmployeeTask{}, errors.Wrap(err, "employee does not exist")
    }

    // Validar que el empleado tenga el rol "empleado"
    if employee.Role != models.EmployeeRoleEmployee {
        return models.EmployeeTask{}, errors.New("tasks can only be assigned to employees with role 'empleado'")
    }

    updatedTask, err := s.employeeTaskRepo.Update(task)
    if err != nil {
        return models.EmployeeTask{}, errors.Wrap(err, "failed to update task")
    }
    return updatedTask, nil
}

// UpdateTaskStatus actualiza el estado de una tarea (para admin y dueño)
func (s *employeeTaskService) UpdateTaskStatus(taskID int, status models.EmployeeTaskStatus) error {
    // Obtener la tarea para validar su existencia
    task, err := s.employeeTaskRepo.FindByID(taskID)
    if err != nil {
        return errors.Wrap(err, "failed to find task")
    }

    // Determinar el valor de completed_at según el estado
    var completedAt *time.Time
    if status == models.TaskStatusCompleted {
        now := time.Now()
        completedAt = &now
    } else {
        completedAt = nil
    }

    // Actualizar el estado y completed_at
    err = s.employeeTaskRepo.UpdateStatus(taskID, status, completedAt)
    if err != nil {
        return errors.Wrap(err, "failed to update task status")
    }

    // Actualizar el campo en la tarea
    task.Status = status
    task.CompletedAt = completedAt

    return nil
}

// DeleteTask elimina una tarea (para admin y dueño)
func (s *employeeTaskService) DeleteTask(taskID int) error {
    err := s.employeeTaskRepo.Delete(taskID)
    if err != nil {
        return errors.Wrap(err, "failed to delete task")
    }
    return nil
}