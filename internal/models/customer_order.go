package models

import (
    "time"

    "github.com/shopspring/decimal"
)

// CustomerOrder representa la tabla customer_orders
type CustomerOrder struct {
    ID               int             `json:"id"`
    TableID          int             `json:"table_id"`
    TotalAmount      decimal.Decimal `json:"total_amount"`
    Status           string          `json:"status"`
    CreatedAt        time.Time       `json:"created_at"`
    OrderDetails    []OrderDetail   `json:"order_details"`
}