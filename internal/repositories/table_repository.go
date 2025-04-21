package repositories

import (
    "database/sql"
    "gastrobar-backend/internal/models"

    "github.com/pkg/errors"
)

type TableRepository interface {
    FindByID(tableID int) (models.Table, error)
    FindByName(tableName string) (models.Table, error)
    FindAll() ([]models.Table, error)
    Count() (int, error)
    Create(table models.Table) (models.Table, error)
    Update(table models.Table) (models.Table, error)
}

type tableRepository struct {
    db *sql.DB
}

func NewTableRepository(db *sql.DB) TableRepository {
    return &tableRepository{db: db}
}

func (r *tableRepository) FindByID(tableID int) (models.Table, error) {
    var table models.Table
    err := r.db.QueryRow(`
        SELECT id, table_name, created_at
        FROM tables
        WHERE id = $1`,
        tableID,
    ).Scan(&table.ID, &table.TableName, &table.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Table{}, errors.Wrap(err, "table not found")
        }
        return models.Table{}, errors.Wrap(err, "failed to query table by ID")
    }
    return table, nil
}

func (r *tableRepository) FindByName(tableName string) (models.Table, error) {
    var table models.Table
    err := r.db.QueryRow(`
        SELECT id, table_name, created_at
        FROM tables
        WHERE table_name = $1`,
        tableName,
    ).Scan(&table.ID, &table.TableName, &table.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Table{}, errors.Wrap(err, "table not found")
        }
        return models.Table{}, errors.Wrap(err, "failed to query table by name")
    }
    return table, nil
}

func (r *tableRepository) FindAll() ([]models.Table, error) {
    rows, err := r.db.Query(`
        SELECT id, table_name, created_at
        FROM tables`)
    if err != nil {
        return nil, errors.Wrap(err, "failed to query tables")
    }
    defer rows.Close()

    var tables []models.Table
    for rows.Next() {
        var table models.Table
        if err := rows.Scan(&table.ID, &table.TableName, &table.CreatedAt); err != nil {
            return nil, errors.Wrap(err, "failed to scan table")
        }
        tables = append(tables, table)
    }
    return tables, nil
}

func (r *tableRepository) Count() (int, error) {
    var count int
    err := r.db.QueryRow(`
        SELECT COUNT(*)
        FROM tables`).
        Scan(&count)
    if err != nil {
        return 0, errors.Wrap(err, "failed to count tables")
    }
    return count, nil
}

func (r *tableRepository) Create(table models.Table) (models.Table, error) {
    var createdTable models.Table
    err := r.db.QueryRow(`
        INSERT INTO tables (table_name, created_at)
        VALUES ($1, CURRENT_TIMESTAMP)
        RETURNING id, table_name, created_at`,
        table.TableName,
    ).Scan(&createdTable.ID, &createdTable.TableName, &createdTable.CreatedAt)
    if err != nil {
        return models.Table{}, errors.Wrap(err, "failed to create table")
    }
    return createdTable, nil
}

func (r *tableRepository) Update(table models.Table) (models.Table, error) {
    var updatedTable models.Table
    err := r.db.QueryRow(`
        UPDATE tables
        SET table_name = $1
        WHERE id = $2
        RETURNING id, table_name, created_at`,
        table.TableName, table.ID,
    ).Scan(&updatedTable.ID, &updatedTable.TableName, &updatedTable.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.Table{}, errors.Wrap(err, "table not found")
        }
        return models.Table{}, errors.Wrap(err, "failed to update table")
    }
    return updatedTable, nil
}