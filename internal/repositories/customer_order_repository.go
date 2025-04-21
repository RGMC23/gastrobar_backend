package repositories

import (
    "database/sql"
    "gastrobar-backend/internal/models"

    "github.com/pkg/errors"
)

type CustomerOrderRepository interface {
    Create(order models.CustomerOrder) (models.CustomerOrder, error)
    FindByID(id int) (models.CustomerOrder, error)
    FindPendingByTableID(tableID int) (models.CustomerOrder, error)
    FindCompletedByTableID(tableID int) (models.CustomerOrder, error)
    UpdateStatus(id int, status string) (models.CustomerOrder, error)
}

type customerOrderRepository struct {
    db *sql.DB
}

func NewCustomerOrderRepository(db *sql.DB) CustomerOrderRepository {
    return &customerOrderRepository{db: db}
}

func (r *customerOrderRepository) Create(order models.CustomerOrder) (models.CustomerOrder, error) {
    var createdOrder models.CustomerOrder
    err := r.db.QueryRow(`
        INSERT INTO customer_orders (table_id, total_amount, status, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, table_id, total_amount, status, created_at`,
        order.TableID, order.TotalAmount, order.Status, order.CreatedAt,
    ).Scan(
        &createdOrder.ID,
        &createdOrder.TableID,
        &createdOrder.TotalAmount,
        &createdOrder.Status,
        &createdOrder.CreatedAt,
    )
    if err != nil {
        return models.CustomerOrder{}, errors.Wrap(err, "failed to create customer order")
    }

    return createdOrder, nil
}

func (r *customerOrderRepository) FindByID(id int) (models.CustomerOrder, error) {
    var order models.CustomerOrder
    err := r.db.QueryRow(`
        SELECT id, table_id, total_amount, status, created_at
        FROM customer_orders
        WHERE id = $1`,
        id,
    ).Scan(
        &order.ID,
        &order.TableID,
        &order.TotalAmount,
        &order.Status,
        &order.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.CustomerOrder{}, errors.New("customer order not found")
        }
        return models.CustomerOrder{}, errors.Wrap(err, "failed to find customer order")
    }
    return order, nil
}

func (r *customerOrderRepository) FindPendingByTableID(tableID int) (models.CustomerOrder, error) {
    var order models.CustomerOrder
    err := r.db.QueryRow(`
        SELECT id, table_id, total_amount, status, created_at
        FROM customer_orders
        WHERE table_id = $1 AND status = 'pending'`,
        tableID,
    ).Scan(
        &order.ID,
        &order.TableID,
        &order.TotalAmount,
        &order.Status,
        &order.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.CustomerOrder{}, errors.New("no pending customer order found for this table")
        }
        return models.CustomerOrder{}, errors.Wrap(err, "failed to find pending customer order")
    }
    return order, nil
}

// internal/repositories/customer_order_repository.go
func (r *customerOrderRepository) FindCompletedByTableID(tableID int) (models.CustomerOrder, error) {
    var order models.CustomerOrder
    err := r.db.QueryRow(`
        SELECT id, table_id, total_amount, status, created_at
        FROM customer_orders
        WHERE table_id = $1 AND status = 'completed'
        ORDER BY created_at DESC
        LIMIT 1`,
        tableID,
    ).Scan(
        &order.ID,
        &order.TableID,
        &order.TotalAmount,
        &order.Status,
        &order.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.CustomerOrder{}, errors.New("no completed customer order found for this table")
        }
        return models.CustomerOrder{}, errors.Wrap(err, "failed to query completed customer order by table id")
    }
    return order, nil
}

func (r *customerOrderRepository) UpdateStatus(id int, status string) (models.CustomerOrder, error) {
    var updatedOrder models.CustomerOrder
    err := r.db.QueryRow(`
        UPDATE customer_orders
        SET status = $2
        WHERE id = $1
        RETURNING id, table_id, total_amount, status, created_at`,
        id, status,
    ).Scan(
        &updatedOrder.ID,
        &updatedOrder.TableID,
        &updatedOrder.TotalAmount,
        &updatedOrder.Status,
        &updatedOrder.CreatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.CustomerOrder{}, errors.New("customer order not found")
        }
        return models.CustomerOrder{}, errors.Wrap(err, "failed to update customer order status")
    }
    return updatedOrder, nil
}