package models

import "time"

type OrderDetail struct {
    ID         int       `json:"id"`
    OrderID    int       `json:"order_id"`
    MenuItemID int       `json:"menu_item_id,omitempty"`
    MenuItem   MenuItem  `json:"menu_item"`              
    Quantity   int       `json:"quantity"`
    CreatedAt  time.Time `json:"created_at"`
}