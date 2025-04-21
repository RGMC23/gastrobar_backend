package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "strings"

    "gastrobar-backend/internal/models"
    "gastrobar-backend/internal/services"

    "github.com/gorilla/mux"
)

type OrderDetailHandler struct {
    orderDetailSvc services.OrderDetailService
}

func NewOrderDetailHandler(orderDetailSvc services.OrderDetailService) *OrderDetailHandler {
    return &OrderDetailHandler{
        orderDetailSvc: orderDetailSvc,
    }
}

type CreateOrderDetailRequest struct {
    OrderID    int `json:"order_id"` // Opcional
    TableID    int `json:"table_id"`
    MenuItemID int `json:"menu_item_id"`
    Quantity   int `json:"quantity"`
}

func (h *OrderDetailHandler) CreateOrderDetailHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var request CreateOrderDetailRequest
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validar campos requeridos
        if request.TableID <= 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "table_id is required"})
            return
        }
        if request.MenuItemID <= 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "menu_item_id is required"})
            return
        }
        if request.Quantity <= 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "quantity must be greater than 0"})
            return
        }

        // Crear el order_detail
        orderDetail := models.OrderDetail{
            OrderID:    request.OrderID,
            MenuItemID: request.MenuItemID,
            Quantity:   request.Quantity,
        }
        createdDetail, err := h.orderDetailSvc.CreateOrderDetail(orderDetail, request.TableID)
        if err != nil {
            if strings.Contains(err.Error(), "failed to find table") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "table not found"})
                return
            }
            if strings.Contains(err.Error(), "failed to find menu item") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "menu item not found"})
                return
            }
            if strings.Contains(err.Error(), "insufficient stock for item") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }
            if strings.Contains(err.Error(), "a pending customer order already exists") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }
            if strings.Contains(err.Error(), "customer order is already completed") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }
            if strings.Contains(err.Error(), "customer order does not belong to the specified table") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }
            log.Printf("Error creating order detail: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdDetail)
    }
}

func (h *OrderDetailHandler) GetOrderDetailByIDHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "invalid order detail ID"})
            return
        }

        orderDetail, err := h.orderDetailSvc.GetOrderDetailByID(id)
        if err != nil {
            if strings.Contains(err.Error(), "order detail not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "order detail not found"})
                return
            }
            log.Printf("Error getting order detail: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(orderDetail)
    }
}

func (h *OrderDetailHandler) GetOrderDetailsByOrderIDHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        orderIDStr := vars["order_id"]
        orderID, err := strconv.Atoi(orderIDStr)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "invalid order ID"})
            return
        }

        orderDetails, err := h.orderDetailSvc.GetOrderDetailsByOrderID(orderID)
        if err != nil {
            if strings.Contains(err.Error(), "failed to find customer order") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "customer order not found"})
                return
            }
            log.Printf("Error getting order details: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(orderDetails)
    }
}

func (h *OrderDetailHandler) UpdateOrderDetailHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "invalid order detail ID"})
            return
        }

        var orderDetail models.OrderDetail
        if err := json.NewDecoder(r.Body).Decode(&orderDetail); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Asegurar que el ID de la URL coincida con el del cuerpo
        orderDetail.ID = id

        // Validar campos requeridos
        if orderDetail.OrderID <= 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "order_id is required"})
            return
        }
        if orderDetail.MenuItemID <= 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "menu_item_id is required"})
            return
        }
        if orderDetail.Quantity <= 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "quantity must be greater than 0"})
            return
        }

        // Actualizar el order_detail
        updatedDetail, err := h.orderDetailSvc.UpdateOrderDetail(orderDetail)
        if err != nil {
            if strings.Contains(err.Error(), "order detail not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "order detail not found"})
                return
            }
            if strings.Contains(err.Error(), "failed to find customer order") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "customer order not found"})
                return
            }
            if strings.Contains(err.Error(), "failed to find menu item") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "menu item not found"})
                return
            }
            if strings.Contains(err.Error(), "insufficient stock for item") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }
            if strings.Contains(err.Error(), "customer order is already completed") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }
            log.Printf("Error updating order detail: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedDetail)
    }
}

func (h *OrderDetailHandler) DeleteOrderDetailHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        idStr := vars["id"]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "invalid order detail ID"})
            return
        }

        err = h.orderDetailSvc.DeleteOrderDetail(id)
        if err != nil {
            if strings.Contains(err.Error(), "order detail not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "order detail not found"})
                return
            }
            if strings.Contains(err.Error(), "failed to find customer order") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "customer order not found"})
                return
            }
            if strings.Contains(err.Error(), "customer order is already completed") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                return
            }
            log.Printf("Error deleting order detail: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}