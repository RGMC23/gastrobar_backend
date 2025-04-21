package models

import (
    "time"

    "github.com/shopspring/decimal"
)

// MenuItem representa la tabla menu_items
type MenuItem struct {
    ID          int             `json:"id"`
    ItemName    string          `json:"item_name"`
    Category    string          `json:"category"`
    Price       decimal.Decimal `json:"price"` // Usamos decimal.Decimal para NUMERIC
    Stock       int             `json:"stock"`
    Description string          `json:"description"`
    CreatedAt   time.Time       `json:"created_at"`
}