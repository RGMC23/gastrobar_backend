package services

import (
	"gastrobar-backend/internal/models"
	"gastrobar-backend/internal/repositories"
	"strings"

	"github.com/pkg/errors"
)

const MaxTables = 4

type TableService interface {
    CreateTable(table models.Table) (models.Table, error)
    GetTable(tableID int) (models.Table, error)
    ListTables() ([]models.Table, error)
    UpdateTable(table models.Table) (models.Table, error)
}

type tableService struct {
    tableRepo repositories.TableRepository
}

func NewTableService(tableRepo repositories.TableRepository) TableService {
    return &tableService{
        tableRepo: tableRepo,
    }
}

func (s *tableService) CreateTable(table models.Table) (models.Table, error) {
    // Validar que el nombre no esté vacío
    if table.TableName == "" {
        return models.Table{}, errors.New("table name cannot be empty")
    }

    // Validar que el nombre sea único
    _, err := s.tableRepo.FindByName(table.TableName)
    if err == nil {
        return models.Table{}, errors.New("table name already exists")
    }
    if !strings.Contains(err.Error(), "table not found") {
        return models.Table{}, errors.Wrap(err, "failed to check table name uniqueness")
    }

    // Validar el límite de mesas
    count, err := s.tableRepo.Count()
    if err != nil {
        return models.Table{}, errors.Wrap(err, "failed to count tables")
    }
    if count >= MaxTables {
        return models.Table{}, errors.New("maximum number of tables reached")
    }


    // Crear la mesa en el repositorio
    createdTable, err := s.tableRepo.Create(table)
    if err != nil {
        return models.Table{}, errors.Wrap(err, "failed to create table")
    }

    return createdTable, nil
}

func (s *tableService) GetTable(tableID int) (models.Table, error) {
    table, err := s.tableRepo.FindByID(tableID)
    if err != nil {
        return models.Table{}, errors.Wrap(err, "failed to get table")
    }
    return table, nil
}

func (s *tableService) ListTables() ([]models.Table, error) {
    tables, err := s.tableRepo.FindAll()
    if err != nil {
        return nil, errors.Wrap(err, "failed to list tables")
    }
    return tables, nil
}

func (s *tableService) UpdateTable(table models.Table) (models.Table, error) {
    // Validar que el nombre no esté vacío
    if table.TableName == "" {
        return models.Table{}, errors.New("table name cannot be empty")
    }

    // Obtener la mesa actual para verificar si el nombre cambió
    currentTable, err := s.tableRepo.FindByID(table.ID)
    if err != nil {
        return models.Table{}, errors.Wrap(err, "failed to find table")
    }

    // Validar unicidad del nombre si cambió
    if table.TableName != currentTable.TableName {
        _, err := s.tableRepo.FindByName(table.TableName)
        if err == nil {
            return models.Table{}, errors.New("table name already exists")
        }
        if !strings.Contains(err.Error(), "table not found") {
            return models.Table{}, errors.Wrap(err, "failed to check table name uniqueness")
        }
    }

    // Actualizar la mesa en el repositorio
    updatedTable, err := s.tableRepo.Update(table)
    if err != nil {
        return models.Table{}, errors.Wrap(err, "failed to update table")
    }
    return updatedTable, nil
}