package services

import (
	"gastrobar-backend/internal/models"
	"gastrobar-backend/internal/repositories"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type MenuItemService interface {
    CreateMenuItem(item models.MenuItem) (models.MenuItem, error)
    GetMenuItem(itemID int) (models.MenuItem, error)
    ListMenuItems() ([]models.MenuItem, error)
    UpdateMenuItem(item models.MenuItem) (models.MenuItem, error)
    DeleteMenuItem(itemID int) error
}

type menuItemService struct {
    menuItemRepo repositories.MenuItemRepository
}

func NewMenuItemService(menuItemRepo repositories.MenuItemRepository) MenuItemService {
    return &menuItemService{
        menuItemRepo: menuItemRepo,
    }
}

func (s *menuItemService) CreateMenuItem(item models.MenuItem) (models.MenuItem, error) {
    // Validar que el nombre no esté vacío
    if item.ItemName == "" {
        return models.MenuItem{}, errors.New("item name cannot be empty")
    }

    // Validar que el nombre sea único
    _, err := s.menuItemRepo.FindByName(item.ItemName)
    if err == nil {
        return models.MenuItem{}, errors.New("item name already exists")
    }
    if !strings.Contains(err.Error(), "item not found") {
        return models.MenuItem{}, errors.Wrap(err, "failed to check item name uniqueness")
    }

    // Validar que la descripción no esté vacía
    if item.Description == "" {
        return models.MenuItem{}, errors.New("description cannot be empty")
    }

    // Validar que el precio sea mayor a 0
    if item.Price.LessThanOrEqual(decimal.Zero) {
        return models.MenuItem{}, errors.New("price must be greater than 0")
    }

    // Validar que el stock no sea negativo
    if item.Stock < 0 {
        return models.MenuItem{}, errors.New("stock cannot be negative")
    }

    // Crear el ítem en el repositorio
    createdItem, err := s.menuItemRepo.Create(item)
    if err != nil {
        return models.MenuItem{}, errors.Wrap(err, "failed to create item")
    }
    return createdItem, nil
}

func (s *menuItemService) GetMenuItem(itemID int) (models.MenuItem, error) {
    item, err := s.menuItemRepo.FindByID(itemID)
    if err != nil {
        return models.MenuItem{}, errors.Wrap(err, "failed to get item")
    }
    return item, nil
}

func (s *menuItemService) ListMenuItems() ([]models.MenuItem, error) {
    items, err := s.menuItemRepo.FindAll()
    if err != nil {
        return nil, errors.Wrap(err, "failed to list items")
    }
    return items, nil
}

func (s *menuItemService) UpdateMenuItem(item models.MenuItem) (models.MenuItem, error) {
    // Validar que el nombre no esté vacío
    if item.ItemName == "" {
        return models.MenuItem{}, errors.New("item name cannot be empty")
    }

    // Obtener el ítem actual para verificar si el nombre cambió
    currentItem, err := s.menuItemRepo.FindByID(item.ID)
    if err != nil {
        return models.MenuItem{}, errors.Wrap(err, "failed to find item")
    }

    // Validar unicidad del nombre si cambió
    if item.ItemName != currentItem.ItemName {
        _, err := s.menuItemRepo.FindByName(item.ItemName)
        if err == nil {
            return models.MenuItem{}, errors.New("item name already exists")
        }
        if !strings.Contains(err.Error(), "item not found") {
            return models.MenuItem{}, errors.Wrap(err, "failed to check item name uniqueness")
        }
    }

    // Validar que la descripción no esté vacía
    if item.Description == "" {
        return models.MenuItem{}, errors.New("description cannot be empty")
    }

    // Validar que el precio sea mayor a 0
    if item.Price.LessThanOrEqual(decimal.Zero) {
        return models.MenuItem{}, errors.New("price must be greater than 0")
    }

    // Validar que el stock no sea negativo
    if item.Stock < 0 {
        return models.MenuItem{}, errors.New("stock cannot be negative")
    }

    // Actualizar el ítem en el repositorio
    updatedItem, err := s.menuItemRepo.Update(item)
    if err != nil {
        return models.MenuItem{}, errors.Wrap(err, "failed to update item")
    }
    return updatedItem, nil
}

func (s *menuItemService) DeleteMenuItem(itemID int) error {
    err := s.menuItemRepo.Delete(itemID)
    if err != nil {
        return errors.Wrap(err, "failed to delete item")
    }
    return nil
}