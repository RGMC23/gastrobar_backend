package repositories

import (
    "database/sql"
    "time"

    "gastrobar-backend/internal/models"

    "github.com/pkg/errors"
)

type EmployeeTaskRepository interface {
    FindByID(taskID int) (models.EmployeeTask, error)
    FindAll() ([]models.EmployeeTask, error)
    FindByEmployeeID(employeeID int) ([]models.EmployeeTask, error)
    Create(task models.EmployeeTask) (models.EmployeeTask, error)
    Update(task models.EmployeeTask) (models.EmployeeTask, error)
    UpdateStatus(taskID int, status models.EmployeeTaskStatus, completedAt *time.Time) error
    Delete(taskID int) error
}

type employeeTaskRepository struct {
    db *sql.DB
}

func NewEmployeeTaskRepository(db *sql.DB) EmployeeTaskRepository {
    return &employeeTaskRepository{db: db}
}

func (r *employeeTaskRepository) FindByID(taskID int) (models.EmployeeTask, error) {
    var task models.EmployeeTask
    var completedAt sql.NullTime

    err := r.db.QueryRow(`
        SELECT id, employee_id, task_description, status, assigned_at, completed_at, created_at
        FROM employee_tasks
        WHERE id = $1`,
        taskID,
    ).Scan(&task.ID, &task.EmployeeID, &task.TaskDescription, &task.Status, &task.AssignedAt, &completedAt, &task.CreatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return models.EmployeeTask{}, errors.Wrap(err, "task not found")
        }
        return models.EmployeeTask{}, errors.Wrap(err, "failed to query task by ID")
    }

    // Convertir completedAt (sql.NullTime) a *time.Time
    if completedAt.Valid {
        task.CompletedAt = &completedAt.Time
    }

    return task, nil
}

func (r *employeeTaskRepository) FindAll() ([]models.EmployeeTask, error) {
    rows, err := r.db.Query(`
        SELECT id, employee_id, task_description, status, assigned_at, completed_at, created_at
        FROM employee_tasks`)
    if err != nil {
        return nil, errors.Wrap(err, "failed to query tasks")
    }
    defer rows.Close()

    var tasks []models.EmployeeTask
    for rows.Next() {
        var task models.EmployeeTask
        var completedAt sql.NullTime

        if err := rows.Scan(&task.ID, &task.EmployeeID, &task.TaskDescription, &task.Status, &task.AssignedAt, &completedAt, &task.CreatedAt); err != nil {
            return nil, errors.Wrap(err, "failed to scan task")
        }

        // Convertir completedAt (sql.NullTime) a *time.Time
        if completedAt.Valid {
            task.CompletedAt = &completedAt.Time
        }

        tasks = append(tasks, task)
    }
    return tasks, nil
}

func (r *employeeTaskRepository) FindByEmployeeID(employeeID int) ([]models.EmployeeTask, error) {
    rows, err := r.db.Query(`
        SELECT id, employee_id, task_description, status, assigned_at, completed_at, created_at
        FROM employee_tasks
        WHERE employee_id = $1`,
        employeeID,
    )
    if err != nil {
        return nil, errors.Wrap(err, "failed to query tasks by employee ID")
    }
    defer rows.Close()

    var tasks []models.EmployeeTask
    for rows.Next() {
        var task models.EmployeeTask
        var completedAt sql.NullTime

        if err := rows.Scan(&task.ID, &task.EmployeeID, &task.TaskDescription, &task.Status, &task.AssignedAt, &completedAt, &task.CreatedAt); err != nil {
            return nil, errors.Wrap(err, "failed to scan task")
        }

        // Convertir completedAt (sql.NullTime) a *time.Time
        if completedAt.Valid {
            task.CompletedAt = &completedAt.Time
        }

        tasks = append(tasks, task)
    }
    return tasks, nil
}

func (r *employeeTaskRepository) Create(task models.EmployeeTask) (models.EmployeeTask, error) {
    var createdTask models.EmployeeTask
    var completedAt sql.NullTime

    err := r.db.QueryRow(`
        INSERT INTO employee_tasks (employee_id, task_description, status, assigned_at, completed_at, created_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
        RETURNING id, employee_id, task_description, status, assigned_at, completed_at, created_at`,
        task.EmployeeID, task.TaskDescription, task.Status, task.AssignedAt, task.CompletedAt,
    ).Scan(&createdTask.ID, &createdTask.EmployeeID, &createdTask.TaskDescription, &createdTask.Status, &createdTask.AssignedAt, &completedAt, &createdTask.CreatedAt)

    if err != nil {
        return models.EmployeeTask{}, errors.Wrap(err, "failed to create task")
    }

    // Convertir completedAt (sql.NullTime) a *time.Time
    if completedAt.Valid {
        createdTask.CompletedAt = &completedAt.Time
    }

    return createdTask, nil
}

func (r *employeeTaskRepository) Update(task models.EmployeeTask) (models.EmployeeTask, error) {
    var updatedTask models.EmployeeTask
    var completedAt sql.NullTime

    err := r.db.QueryRow(`
        UPDATE employee_tasks
        SET employee_id = $1, task_description = $2, status = $3, assigned_at = $4, completed_at = $5
        WHERE id = $6
        RETURNING id, employee_id, task_description, status, assigned_at, completed_at, created_at`,
        task.EmployeeID, task.TaskDescription, task.Status, task.AssignedAt, task.CompletedAt, task.ID,
    ).Scan(&updatedTask.ID, &updatedTask.EmployeeID, &updatedTask.TaskDescription, &updatedTask.Status, &updatedTask.AssignedAt, &completedAt, &updatedTask.CreatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return models.EmployeeTask{}, errors.Wrap(err, "task not found")
        }
        return models.EmployeeTask{}, errors.Wrap(err, "failed to update task")
    }

    // Convertir completedAt (sql.NullTime) a *time.Time
    if completedAt.Valid {
        updatedTask.CompletedAt = &completedAt.Time
    }

    return updatedTask, nil
}

func (r *employeeTaskRepository) UpdateStatus(taskID int, status models.EmployeeTaskStatus, completedAt *time.Time) error {
    result, err := r.db.Exec(`
        UPDATE employee_tasks
        SET status = $1, completed_at = $2
        WHERE id = $3`,
        status, completedAt, taskID,
    )
    if err != nil {
        return errors.Wrap(err, "failed to update task status")
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.Wrap(err, "failed to check rows affected")
    }
    if rowsAffected == 0 {
        return errors.New("task not found")
    }
    return nil
}

func (r *employeeTaskRepository) Delete(taskID int) error {
    result, err := r.db.Exec(`
        DELETE FROM employee_tasks
        WHERE id = $1`,
        taskID,
    )
    if err != nil {
        return errors.Wrap(err, "failed to delete task")
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.Wrap(err, "failed to check rows affected")
    }
    if rowsAffected == 0 {
        return errors.New("task not found")
    }
    return nil
}