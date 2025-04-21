package models

import "time"

// Table representa la tabla tables
type Table struct {
    ID          int       `json:"id"`
    TableName   string    `json:"table_name"`
    CreatedAt   time.Time `json:"created_at"`
}