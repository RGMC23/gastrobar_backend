package repositories

import (
    "database/sql"

    "gastrobar-backend/internal/models"

    "github.com/pkg/errors"
    "github.com/shopspring/decimal"
)

type MenuItemRepository interface {
    FindByID(itemID int) (models.MenuItem, error)
    FindByName(itemName string) (models.MenuItem, error)
    FindAll() ([]models.MenuItem, error)
    Create(item models.MenuItem) (models.MenuItem, error)
    Update(item models.MenuItem) (models.MenuItem, error)
    Delete(itemID int) error
}

type menuItemRepository struct {
    db *sql.DB
}

func NewMenuItemRepository(db *sql.DB) MenuItemRepository {
    return &menuItemRepository{db: db}
}

func (r *menuItemRepository) FindByID(itemID int) (models.MenuItem, error) {
    var item models.MenuItem
    err := r.db.QueryRow(`
        SELECT id, item_name, category, price, stock, description, created_at
        FROM menu_items
        WHERE id = $1`,
        itemID,
    ).Scan(&item.ID, &item.ItemName, &item.Category, &item.Price, &item.Stock, &item.Description, &item.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.MenuItem{}, errors.Wrap(err, "item not found")
        }
        return models.MenuItem{}, errors.Wrap(err, "failed to query item by ID")
    }
    return item, nil
}

func (r *menuItemRepository) FindByName(itemName string) (models.MenuItem, error) {
    var item models.MenuItem
    err := r.db.QueryRow(`
        SELECT id, item_name, category, price, stock, description, created_at
        FROM menu_items
        WHERE item_name = $1`,
        itemName,
    ).Scan(&item.ID, &item.ItemName, &item.Category, &item.Price, &item.Stock, &item.Description, &item.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.MenuItem{}, errors.Wrap(err, "item not found")
        }
        return models.MenuItem{}, errors.Wrap(err, "failed to query item by name")
    }
    return item, nil
}

func (r *menuItemRepository) FindAll() ([]models.MenuItem, error) {
    rows, err := r.db.Query(`
        SELECT id, item_name, category, price, stock, description, created_at
        FROM menu_items`)
    if err != nil {
        return nil, errors.Wrap(err, "failed to query items")
    }
    defer rows.Close()

    var items []models.MenuItem
    for rows.Next() {
        var item models.MenuItem
        var price string // Temporal para manejar decimal.Decimal
        if err := rows.Scan(&item.ID, &item.ItemName, &item.Category, &price, &item.Stock, &item.Description, &item.CreatedAt); err != nil {
            return nil, errors.Wrap(err, "failed to scan item")
        }
        item.Price, _ = decimal.NewFromString(price) // Convertir a decimal.Decimal
        items = append(items, item)
    }
    return items, nil
}

func (r *menuItemRepository) Create(item models.MenuItem) (models.MenuItem, error) {
    var createdItem models.MenuItem
    err := r.db.QueryRow(`
        INSERT INTO menu_items (item_name, category, price, stock, description, created_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
        RETURNING id, item_name, category, price, stock, description, created_at`,
        item.ItemName, item.Category, item.Price.String(), item.Stock, item.Description,
    ).Scan(&createdItem.ID, &createdItem.ItemName, &createdItem.Category, &createdItem.Price, &createdItem.Stock, &createdItem.Description, &createdItem.CreatedAt)
    if err != nil {
        return models.MenuItem{}, errors.Wrap(err, "failed to create item")
    }
    return createdItem, nil
}

func (r *menuItemRepository) Update(item models.MenuItem) (models.MenuItem, error) {
    var updatedItem models.MenuItem
    err := r.db.QueryRow(`
        UPDATE menu_items
        SET item_name = $1, category = $2, price = $3, stock = $4, description = $5
        WHERE id = $6
        RETURNING id, item_name, category, price, stock, description, created_at`,
        item.ItemName, item.Category, item.Price.String(), item.Stock, item.Description, item.ID,
    ).Scan(&updatedItem.ID, &updatedItem.ItemName, &updatedItem.Category, &updatedItem.Price, &updatedItem.Stock, &updatedItem.Description, &updatedItem.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return models.MenuItem{}, errors.Wrap(err, "item not found")
        }
        return models.MenuItem{}, errors.Wrap(err, "failed to update item")
    }
    return updatedItem, nil
}

func (r *menuItemRepository) Delete(itemID int) error {
    result, err := r.db.Exec(`
        DELETE FROM menu_items
        WHERE id = $1`,
        itemID,
    )
    if err != nil {
        return errors.Wrap(err, "failed to delete item")
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.Wrap(err, "failed to check rows affected")
    }
    if rowsAffected == 0 {
        return errors.New("item not found")
    }
    return nil
}