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

type MenuItemHandler struct {
    menuItemSvc services.MenuItemService
}

func NewMenuItemHandler(menuItemSvc services.MenuItemService) *MenuItemHandler {
    return &MenuItemHandler{
        menuItemSvc: menuItemSvc,
    }
}

func (h *MenuItemHandler) CreateMenuItemHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var item models.MenuItem
        if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        createdItem, err := h.menuItemSvc.CreateMenuItem(item)
        if err != nil {
            if strings.Contains(err.Error(), "item name cannot be empty") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Item name cannot be empty"})
                return
            }
            if strings.Contains(err.Error(), "item name already exists") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Item name already exists"})
                return
            }
            if strings.Contains(err.Error(), "description cannot be empty") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Description cannot be empty"})
                return
            }
            if strings.Contains(err.Error(), "price must be greater than 0") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Price must be greater than 0"})
                return
            }
            if strings.Contains(err.Error(), "stock cannot be negative") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Stock cannot be negative"})
                return
            }
            log.Printf("Error creating item: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdItem)
    }
}

func (h *MenuItemHandler) GetMenuItemHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        itemIDStr := vars["id"]
        itemID, err := strconv.Atoi(itemIDStr)
        if err != nil {
            http.Error(w, "Invalid item ID", http.StatusBadRequest)
            return
        }

        item, err := h.menuItemSvc.GetMenuItem(itemID)
        if err != nil {
            if strings.Contains(err.Error(), "item not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
                return
            }
            log.Printf("Error getting item: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(item)
    }
}

func (h *MenuItemHandler) ListMenuItemsHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        items, err := h.menuItemSvc.ListMenuItems()
        if err != nil {
            log.Printf("Error listing items: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(items)
    }
}

func (h *MenuItemHandler) UpdateMenuItemHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        itemIDStr := vars["id"]
        itemID, err := strconv.Atoi(itemIDStr)
        if err != nil {
            http.Error(w, "Invalid item ID", http.StatusBadRequest)
            return
        }

        var item models.MenuItem
        if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        item.ID = itemID

        updatedItem, err := h.menuItemSvc.UpdateMenuItem(item)
        if err != nil {
            if strings.Contains(err.Error(), "item name cannot be empty") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Item name cannot be empty"})
                return
            }
            if strings.Contains(err.Error(), "item name already exists") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Item name already exists"})
                return
            }
            if strings.Contains(err.Error(), "description cannot be empty") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Description cannot be empty"})
                return
            }
            if strings.Contains(err.Error(), "price must be greater than 0") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Price must be greater than 0"})
                return
            }
            if strings.Contains(err.Error(), "stock cannot be negative") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"error": "Stock cannot be negative"})
                return
            }
            if strings.Contains(err.Error(), "item not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
                return
            }
            log.Printf("Error updating item: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(updatedItem)
    }
}

func (h *MenuItemHandler) DeleteMenuItemHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        itemIDStr := vars["id"]
        itemID, err := strconv.Atoi(itemIDStr)
        if err != nil {
            http.Error(w, "Invalid item ID", http.StatusBadRequest)
            return
        }

        err = h.menuItemSvc.DeleteMenuItem(itemID)
        if err != nil {
            if strings.Contains(err.Error(), "item not found") {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"error": "Item not found"})
                return
            }
            log.Printf("Error deleting item: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Item deleted successfully"))
    }
}