package models

import "time"

// Business representa la tabla business
type Business struct {
    ID              int       `json:"id"`
    BusinessName    string    `json:"business_name"`
    Address         string    `json:"address"`
    PhoneNumber     string    `json:"phone_number"`
    Email           string    `json:"email"`
    CorporateReason string    `json:"corporate_reason"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}