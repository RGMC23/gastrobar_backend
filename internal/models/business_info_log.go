package models

import "time"

// BusinessInfoLog representa la tabla business_info_logs
type BusinessInfoLog struct {
    ID            int       `json:"id"`
    BusinessInfoID int      `json:"business_info_id"`
    Operation     string    `json:"operation"`
    OldData       string    `json:"old_data"` // JSONB se maneja como string
    NewData       string    `json:"new_data"` // JSONB se maneja como string
    ChangedAt     time.Time `json:"changed_at"`
}