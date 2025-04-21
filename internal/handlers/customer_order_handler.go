package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "strings"

    "gastrobar-backend/internal/services"

    "github.com/gorilla/mux"
)

type CustomerOrderHandler struct {
    customerOrderSvc services.CustomerOrderService
}

func NewCustomerOrderHandler(customerOrderSvc services.CustomerOrderService) *CustomerOrderHandler {
    return &CustomerOrderHandler{
        customerOrderSvc: customerOrderSvc,
    }
}

func (h *CustomerOrderHandler) GetOrderWithDetailsHandler() http.HandlerFunc {
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

        order, err := h.customerOrderSvc.GetOrderWithDetails(orderID)
        if err != nil {
            if strings.Contains(err.Error(), "failed to find customer order") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "customer order not found"})
                return
            }
            log.Printf("Error getting order with details: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(order)
    }
}

func (h *CustomerOrderHandler) CompleteOrderHandler() http.HandlerFunc {
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

        updatedOrder, err := h.customerOrderSvc.CompleteOrder(orderID)
        if err != nil {
            if strings.Contains(err.Error(), "failed to find customer order") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "customer order not found"})
                return
            }
            if strings.Contains(err.Error(), "order is not in 'pending' state") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "order is not in 'pending' state"})
                return
            }
            log.Printf("Error completing order: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedOrder)
    }
}

func (h *CustomerOrderHandler) CompleteOrderByEmployeeHandler() http.HandlerFunc {
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

        updatedOrder, err := h.customerOrderSvc.CompleteOrderByEmployee(orderID)
        if err != nil {
            if strings.Contains(err.Error(), "failed to find customer order") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "customer order not found"})
                return
            }
            if strings.Contains(err.Error(), "order is not in 'pending' state") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "order is not in 'pending' state"})
                return
            }
            log.Printf("Error completing order by employee: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedOrder)
    }
}