package repositories

import (
    "database/sql"

    "gastrobar-backend/internal/models"

    "github.com/pkg/errors"
)

type EmployeeRepository interface {
    FindByEmail(email string) (models.Employee, error)
    FindByUsername(username string) (models.Employee, error)
    FindByID(employeeID int) (models.Employee, error)
    FindAll() ([]models.Employee, error)
    FindAllByRole(role models.EmployeeRole) ([]models.Employee, error)
    Create(employee models.Employee) (models.Employee, error)
    Update(employee models.Employee) (models.Employee, error)
    UpdatePassword(employeeID int, password string) error
}

type employeeRepository struct {
    db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
    return &employeeRepository{db: db}
}

func (r *employeeRepository) FindByEmail(email string) (models.Employee, error) {
    var employee models.Employee
    err := r.db.QueryRow(`
        SELECT id, employee_name, email, phone_number, role, username, password, created_at
        FROM employees
        WHERE email = $1`,
        email,
    ).Scan(&employee.ID, &employee.EmployeeName, &employee.Email, &employee.PhoneNumber, &employee.Role, &employee.Username, &employee.Password, &employee.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Employee{}, errors.Wrap(err, "employee not found")
        }
        return models.Employee{}, errors.Wrap(err, "failed to query employee by email")
    }
    return employee, nil
}

func (r *employeeRepository) FindByUsername(username string) (models.Employee, error) {
    var employee models.Employee
    err := r.db.QueryRow(`
        SELECT id, employee_name, email, phone_number, role, username, password, created_at
        FROM employees
        WHERE username = $1`,
        username,
    ).Scan(&employee.ID, &employee.EmployeeName, &employee.Email, &employee.PhoneNumber, &employee.Role, &employee.Username, &employee.Password, &employee.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Employee{}, errors.Wrap(err, "employee not found")
        }
        return models.Employee{}, errors.Wrap(err, "failed to query employee by username")
    }
    return employee, nil
}

func (r *employeeRepository) FindByID(employeeID int) (models.Employee, error) {
    var employee models.Employee
    err := r.db.QueryRow(`
        SELECT id, employee_name, email, phone_number, role, username, password, created_at
        FROM employees
        WHERE id = $1`,
        employeeID,
    ).Scan(&employee.ID, &employee.EmployeeName, &employee.Email, &employee.PhoneNumber, &employee.Role, &employee.Username, &employee.Password, &employee.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Employee{}, errors.Wrap(err, "employee not found")
        }
        return models.Employee{}, errors.Wrap(err, "failed to query employee by ID")
    }
    return employee, nil
}

func (r *employeeRepository) FindAll() ([]models.Employee, error) {
    rows, err := r.db.Query(`
        SELECT id, employee_name, email, phone_number, role, username, password, created_at
        FROM employees`)
    if err != nil {
        return nil, errors.Wrap(err, "failed to query employees")
    }
    defer rows.Close()

    var employees []models.Employee
    for rows.Next() {
        var employee models.Employee
        if err := rows.Scan(&employee.ID, &employee.EmployeeName, &employee.Email, &employee.PhoneNumber, &employee.Role, &employee.Username, &employee.Password, &employee.CreatedAt); err != nil {
            return nil, errors.Wrap(err, "failed to scan employee")
        }
        employees = append(employees, employee)
    }
    return employees, nil
}

func (r *employeeRepository) FindAllByRole(role models.EmployeeRole) ([]models.Employee, error) {
    rows, err := r.db.Query(`
        SELECT id, employee_name, email, phone_number, role, username, password, created_at
        FROM employees
        WHERE role = $1`,
        role,
    )
    if err != nil {
        return nil, errors.Wrap(err, "failed to query employees by role")
    }
    defer rows.Close()

    var employees []models.Employee
    for rows.Next() {
        var employee models.Employee
        if err := rows.Scan(&employee.ID, &employee.EmployeeName, &employee.Email, &employee.PhoneNumber, &employee.Role, &employee.Username, &employee.Password, &employee.CreatedAt); err != nil {
            return nil, errors.Wrap(err, "failed to scan employee")
        }
        employees = append(employees, employee)
    }
    return employees, nil
}


func (r *employeeRepository) Create(employee models.Employee) (models.Employee, error) {
    var createdEmployee models.Employee
    err := r.db.QueryRow(`
        INSERT INTO employees (employee_name, email, phone_number, role, username, password, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
        RETURNING id, employee_name, email, phone_number, role, username, password, created_at`,
        employee.EmployeeName, employee.Email, employee.PhoneNumber, employee.Role, employee.Username, employee.Password,
    ).Scan(&createdEmployee.ID, &createdEmployee.EmployeeName, &createdEmployee.Email, &createdEmployee.PhoneNumber, &createdEmployee.Role, &createdEmployee.Username, &createdEmployee.Password, &createdEmployee.CreatedAt)
    if err != nil {
        return models.Employee{}, errors.Wrap(err, "failed to create employee")
    }
    return createdEmployee, nil
}

func (r *employeeRepository) Update(employee models.Employee) (models.Employee, error) {
    var updatedEmployee models.Employee
    err := r.db.QueryRow(`
        UPDATE employees
        SET employee_name = $1, email = $2, phone_number = $3, role = $4, username = $5
        WHERE id = $6
        RETURNING id, employee_name, email, phone_number, role, username, password, created_at`,
        employee.EmployeeName, employee.Email, employee.PhoneNumber, employee.Role, employee.Username, employee.ID,
    ).Scan(&updatedEmployee.ID, &updatedEmployee.EmployeeName, &updatedEmployee.Email, &updatedEmployee.PhoneNumber, &updatedEmployee.Role, &updatedEmployee.Username, &updatedEmployee.Password, &updatedEmployee.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Employee{}, errors.Wrap(err, "employee not found")
        }
        return models.Employee{}, errors.Wrap(err, "failed to update employee")
    }
    return updatedEmployee, nil
}

func (r *employeeRepository) UpdatePassword(employeeID int, password string) error {
    result, err := r.db.Exec(`
        UPDATE employees
        SET password = $1
        WHERE id = $2`,
        password, employeeID,
    )
    if err != nil {
        return errors.Wrap(err, "failed to update employee password")
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.Wrap(err, "failed to check rows affected")
    }
    if rowsAffected == 0 {
        return errors.New("employee not found")
    }
    return nil
}