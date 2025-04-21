package repositories

import (
    "database/sql"

    "gastrobar-backend/internal/models"

    "github.com/pkg/errors"
)

type OrderDetailRepository interface {
    Create(orderDetail models.OrderDetail) (models.OrderDetail, error)
    FindByID(id int) (models.OrderDetail, error)
    FindByOrderID(orderID int) ([]models.OrderDetail, error)
    Update(orderDetail models.OrderDetail) (models.OrderDetail, error)
    Delete(id int) error
}

type orderDetailRepository struct {
    db *sql.DB
}

func NewOrderDetailRepository(db *sql.DB) OrderDetailRepository {
    if db == nil {
        panic("database connection cannot be nil")
    }
    return &orderDetailRepository{db: db}
}

func (r *orderDetailRepository) Create(orderDetail models.OrderDetail) (models.OrderDetail, error) {
    var createdDetail models.OrderDetail
    err := r.db.QueryRow(
        `
        INSERT INTO order_details (order_id, menu_item_id, quantity, created_at)
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
        RETURNING id, order_id, menu_item_id, quantity, created_at`,
        orderDetail.OrderID, orderDetail.MenuItemID, orderDetail.Quantity,
    ).Scan(
        &createdDetail.ID,
        &createdDetail.OrderID,
        &createdDetail.MenuItemID,
        &createdDetail.Quantity,
        &createdDetail.CreatedAt,
    )
    if err != nil {
        return models.OrderDetail{}, errors.Wrap(err, "failed to create order detail")
    }
    return createdDetail, nil
}

func (r *orderDetailRepository) FindByID(id int) (models.OrderDetail, error) {
    var orderDetail models.OrderDetail
    err := r.db.QueryRow(`
        SELECT id, order_id, menu_item_id, quantity, created_at
        FROM order_details
        WHERE id = $1`,
        id,
    ).Scan(
        &orderDetail.ID,
        &orderDetail.OrderID,
        &orderDetail.MenuItemID,
        &orderDetail.Quantity,
        &orderDetail.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.OrderDetail{}, errors.Wrap(err, "order detail not found")
        }
        return models.OrderDetail{}, errors.Wrap(err, "failed to query order detail by ID")
    }
    return orderDetail, nil
}

func (r *orderDetailRepository) FindByOrderID(orderID int) ([]models.OrderDetail, error) {
    rows, err := r.db.Query(`
        SELECT 
            od.id, 
            od.order_id, 
            od.quantity, 
            od.created_at,
            mi.id AS menu_item_id, 
            mi.item_name, 
            mi.category, 
            mi.price, 
            mi.stock, 
            mi.description, 
            mi.created_at AS menu_item_created_at
        FROM order_details od
        JOIN menu_items mi ON od.menu_item_id = mi.id
        WHERE od.order_id = $1`,
        orderID,
    )
    if err != nil {
        return nil, errors.Wrap(err, "failed to query order details by order ID")
    }
    defer rows.Close()

    var orderDetails []models.OrderDetail
    for rows.Next() {
        var od models.OrderDetail
        var menuItem models.MenuItem
        if err := rows.Scan(
            &od.ID,
            &od.OrderID,
            &od.Quantity,
            &od.CreatedAt,
            &menuItem.ID,
            &menuItem.ItemName,
            &menuItem.Category,
            &menuItem.Price,
            &menuItem.Stock,
            &menuItem.Description,
            &menuItem.CreatedAt,
        ); err != nil {
            return nil, errors.Wrap(err, "failed to scan order detail")
        }
        od.MenuItem = menuItem
        orderDetails = append(orderDetails, od)
    }
    return orderDetails, nil
}

func (r *orderDetailRepository) Update(orderDetail models.OrderDetail) (models.OrderDetail, error) {
    var updatedDetail models.OrderDetail
    err := r.db.QueryRow(
        `
        UPDATE order_details
        SET order_id = $1, menu_item_id = $2, quantity = $3
        WHERE id = $4
        RETURNING id, order_id, menu_item_id, quantity, created_at`,
        orderDetail.OrderID, orderDetail.MenuItemID, orderDetail.Quantity, orderDetail.ID,
    ).Scan(
        &updatedDetail.ID,
        &updatedDetail.OrderID,
        &updatedDetail.MenuItemID,
        &updatedDetail.Quantity,
        &updatedDetail.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.OrderDetail{}, errors.Wrap(err, "order detail not found")
        }
        return models.OrderDetail{}, errors.Wrap(err, "failed to update order detail")
    }
    return updatedDetail, nil
}

func (r *orderDetailRepository) Delete(id int) error {
    result, err := r.db.Exec(
        `DELETE FROM order_details WHERE id = $1`,
        id,
    )
    if err != nil {
        return errors.Wrap(err, "failed to delete order detail")
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return errors.Wrap(err, "failed to check rows affected")
    }
    if rowsAffected == 0 {
        return errors.New("order detail not found")
    }
    return nil
}